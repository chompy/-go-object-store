package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego/session"
	"github.com/pkg/errors"
)

var globalSessions *session.Manager
var store *Store

func listen(config *Config) error {
	// init store
	store = NewStore(config)
	sessions = make([]*UserSession, 0)

	u := NewUser()
	u.Username = "nathan"
	u.SetPassword("test1234")
	u.Groups = []string{"admin"}
	store.SetUser(u)

	// start sessions
	var err error
	globalSessions, err = session.NewManager(
		"memory",
		&session.ManagerConfig{
			CookieName:     "sessid",
			CookieLifeTime: 3600,
		},
	)
	if err != nil {
		return errors.WithStack(err)
	}
	go globalSessions.GC()
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

func errorResponse(w http.ResponseWriter, err error) {
	logWarnErr(err, "")
	sendResponse(w, errHTTPResponseCode(err), &APIResponse{
		Success: false,
		Message: err.Error(),
	})
}

func sendResponse(w http.ResponseWriter, status int, resp *APIResponse) {
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

func request(res APIResource, w http.ResponseWriter, r *http.Request) {
	// read request
	apiReq := APIRequest{}
	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}
	if len(rBody) > 0 {
		if err := json.Unmarshal(rBody, &apiReq); err != nil {
			errorResponse(w, err)
			return
		}
	}
	apiReq.sanitizeValues()
	// handle request
	switch res {
	case APILogin:
		{
			if apiReq.Username == "" || apiReq.Password == "" {
				errorResponse(w, ErrInvalidCreds)
				return
			}
			// check username/password
			user, err := store.GetUserByUsername(apiReq.Username)
			if err != nil {
				if errors.Is(err, ErrNotFound) {
					errorResponse(w, ErrInvalidCredientials)
					return
				}
				errorResponse(w, err)
				return
			}
			if !user.CheckPassword(apiReq.Password) {
				errorResponse(w, ErrInvalidCredientials)
				return
			}
			// prepare user session
			sess, key := user.NewSession(user, r.RemoteAddr)
			if sess == nil || key == "" {
				errorResponse(w, ErrUnknown)
				return
			}
			checkSessions()
			sessions = append(sessions, sess)
			// send response
			sendResponse(w, http.StatusOK, &APIResponse{
				Success: true,
				Key:     key,
			})
			return
		}
	case APILogout:
		{
			sendResponse(w, http.StatusOK, &APIResponse{
				Success: true,
			})
			return
		}
	case APIGet:
		{
			return
		}
	case APISet:
		{
			sess := getSessionFromKey(apiReq.SessionKey)
			if sess == nil {
				errorResponse(w, ErrPermission)
				return
			}
			user, err := store.GetUser(sess.UserUID)
			if err != nil {
				errorResponse(w, err)
				return
			}
			respObjs := make([]APIObject, 0)
			for _, o := range apiReq.Objects {
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
			sendResponse(w, http.StatusOK, &APIResponse{
				Success: true,
				Objects: respObjs,
			})
			return
		}
	}
	errorResponse(w, ErrInvalidResource)
}

func login(w http.ResponseWriter, r *http.Request) {
	request(APILogin, w, r)
}

func set(w http.ResponseWriter, r *http.Request) {
	request(APISet, w, r)
}

func get(w http.ResponseWriter, r *http.Request) {
	request(APIGet, w, r)
}

func delete(w http.ResponseWriter, r *http.Request) {
	request(APIDelete, w, r)
}

func index(w http.ResponseWriter, r *http.Request) {

	/*switch r.Method {
	case http.MethodPut, http.MethodPost:
		{
			request(APISet, w, r)
			return
		}
	case http.MethodDelete:
		{
			request(APIDelete, w, r)
			return
		}
	}*/
	errorResponse(w, ErrInvalidResource)
}
