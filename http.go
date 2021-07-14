package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/go-object-store/types"
)

var store *Store

func listen(config *Config) error {
	// init store
	store = NewStore(config)
	sessions = make([]*UserSession, 0)
	// init anonymous user
	u, _ := store.GetUserByUsername(anonymousUser)
	if u == nil {
		u := &types.User{
			Username: anonymousUser,
			Groups:   []string{anonymousUser},
		}
		if err := store.SetUser(u); err != nil {
			return errors.WithStack(err)
		}
	}
	// endpoints
	http.HandleFunc("/login", login)
	http.HandleFunc("/set", set)
	http.HandleFunc("/get", get)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/", index)
	// serve http
	logInfo(fmt.Sprintf("Start HTTP server on port %d.", config.HTTP.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.HTTP.Port), nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func parsePostBody(r *http.Request) (types.APIRequest, error) {
	apiReq := types.APIRequest{
		IP: r.RemoteAddr,
	}
	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return apiReq, errors.WithStack(err)
	}
	if len(rBody) > 0 {
		if err := json.Unmarshal(rBody, &apiReq); err != nil {
			return apiReq, errors.WithStack(err)
		}
	}
	sanitizeValues(&apiReq)
	return apiReq, nil
}

func getUserFromSessionKey(key string) (*types.User, error) {
	if key == "" {
		user, err := store.GetUserByUsername(anonymousUser)
		return user, errors.WithStack(err)
	}
	sess := getSessionFromKey(key)
	if sess == nil {
		return nil, errors.WithStack(ErrPermission)
	}
	user, err := store.GetUser(sess.UserUID)
	return user, errors.WithStack(err)
}

func errorResponse(w http.ResponseWriter, err error) {
	logWarnErr(err, "")
	sendResponse(w, errHTTPResponseCode(err), &types.APIResponse{
		Success: false,
		Message: err.Error(),
	})
}

func sendResponse(w http.ResponseWriter, status int, resp *types.APIResponse) {
	w.WriteHeader(status)
	if resp == nil {
		io.WriteString(w, `{"success":false,"message":"An unknown error occurred."`)
		logWarnErr(ErrEmptyReponse, "")
		return
	}
	data, err := json.Marshal(resp)
	if err != nil {
		io.WriteString(w, `{"success":false,"message":"An unknown error occurred."`)
		logWarnErr(err, "failed to encode response")
		return
	}
	if _, err := w.Write(data); err != nil {
		logWarnErr(err, "failed to send response")
	}
}

func request(res types.APIResource, req types.APIRequest, w http.ResponseWriter) {
	// log request
	logAPIRequest(req, res)
	// handle request
	switch res {
	case types.APILogin:
		{
			if req.Username == "" || req.Password == "" {
				errorResponse(w, ErrInvalidCreds)
				return
			}
			// check username/password
			user, err := store.GetUserByUsername(req.Username)
			if err != nil {
				if errors.Is(err, ErrNotFound) {
					errorResponse(w, ErrInvalidCredientials)
					return
				}
				errorResponse(w, err)
				return
			}
			if !checkPassword(req.Password, user.PasswordHash) {
				errorResponse(w, ErrInvalidCredientials)
				return
			}
			// prepare user session
			sess, key := newSession(user, req.IP)
			if sess == nil || key == "" {
				errorResponse(w, ErrUnknown)
				return
			}
			checkSessions()
			sessions = append(sessions, sess)
			// send response
			sendResponse(w, http.StatusOK, &types.APIResponse{
				Success: true,
				Key:     key,
			})
			return
		}
	case types.APILogout:
		{
			sendResponse(w, http.StatusOK, &types.APIResponse{
				Success: true,
			})
			return
		}
	case types.APIGet:
		{
			if len(req.Objects) == 0 {
				errorResponse(w, ErrObjectNotSpecified)
				return
			}
			user, err := getUserFromSessionKey(req.SessionKey)
			if err != nil {
				errorResponse(w, err)
				return
			}
			respObjs := make([]types.APIObject, 0)
			for _, o := range req.Objects {
				// ensure object isn't already in response
				hasObj := false
				for _, ro := range respObjs {
					if ro.UID() == o.UID() {
						hasObj = true
						break
					}
				}
				if hasObj {
					continue
				}
				// fetch
				respObj, err := store.Get(o.Object().UID, user)
				if err != nil {
					errorResponse(w, err)
					return
				}
				respObjs = append(respObjs, respObj.API())
			}
			sendResponse(w, http.StatusOK, &types.APIResponse{
				Success: true,
				Objects: respObjs,
			})
			return
		}
	case types.APISet:
		{
			user, err := getUserFromSessionKey(req.SessionKey)
			if err != nil {
				errorResponse(w, err)
				return
			}
			respObjs := make([]types.APIObject, 0)
			for _, o := range req.Objects {
				if o == nil {
					continue
				}
				fullObj := o.Object()
				if err := store.Set(fullObj, user); err != nil {
					errorResponse(w, err)
					return
				}
				respObjs = append(respObjs, fullObj.API())
			}
			sendResponse(w, http.StatusOK, &types.APIResponse{
				Success: true,
				Objects: respObjs,
			})
			return
		}
	case types.APIDelete:
		{
			user, err := getUserFromSessionKey(req.SessionKey)
			if err != nil {
				errorResponse(w, err)
				return
			}
			for _, o := range req.Objects {
				if o == nil {
					continue
				}
				if err := store.Delete(o.Object(), user); err != nil {
					errorResponse(w, err)
					return
				}
			}
			sendResponse(w, http.StatusOK, &types.APIResponse{
				Success: true,
			})
		}
	case types.APIQuery:
		{
			user, err := getUserFromSessionKey(req.SessionKey)
			if err != nil {
				errorResponse(w, err)
				return
			}
			if req.Query == "" {
				errorResponse(w, ErrInvalidArg)
				return
			}
			objs, err := store.Query(req.Query, user)
			if err != nil {
				errorResponse(w, err)
				return
			}
			respObjs := make([]types.APIObject, 0)
			for _, o := range objs {
				respObjs = append(respObjs, o.API())
			}
			sendResponse(w, http.StatusOK, &types.APIResponse{
				Success: true,
				Objects: respObjs,
			})
			return
		}
	}
	errorResponse(w, ErrInvalidResource)
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		{
			req, err := parsePostBody(r)
			if err != nil {
				errorResponse(w, err)
				return
			}
			request(types.APILogin, req, w)
			return
		}
	}
	errorResponse(w, ErrAPIInvalidMethod)
}

func set(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		{
			req, err := parsePostBody(r)
			if err != nil {
				errorResponse(w, err)
				return
			}
			request(types.APISet, req, w)
			return
		}
	}
	errorResponse(w, ErrAPIInvalidMethod)
}

func get(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			uids := strings.Split(r.URL.Query().Get("uid"), ",")
			req := types.APIRequest{
				SessionKey: r.URL.Query().Get("key"),
				Objects:    make([]types.APIObject, 0),
			}
			for _, uid := range uids {
				if uid != "" {
					req.Objects = append(req.Objects, types.APIObject{"_uid": uid})
				}
			}
			request(types.APIGet, req, w)
			return
		}
	case http.MethodPost:
		{
			req, err := parsePostBody(r)
			if err != nil {
				errorResponse(w, err)
				return
			}
			request(types.APIGet, req, w)
			return
		}
	}
	errorResponse(w, ErrAPIInvalidMethod)
}

func delete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodDelete:
		{
			req, err := parsePostBody(r)
			if err != nil {
				errorResponse(w, err)
				return
			}
			request(types.APIDelete, req, w)
			return
		}
	}
	errorResponse(w, ErrAPIInvalidMethod)
}

func index(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, ErrInvalidResource)
}
