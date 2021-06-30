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

func listen(port uint16) error {
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

	// serve http
	logInfo(fmt.Sprintf("Start HTTP server on port %d.", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func errorResponse(w http.ResponseWriter, err error) {
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

func request(res APIResource, w http.ResponseWriter, r *http.Request) error {
	// read request
	apiReq := APIRequest{}
	rBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(rBody) > 0 {
		if err := json.Unmarshal(rBody, &apiReq); err != nil {
			return errors.WithStack(err)
		}
	}
	apiReq.sanitizeValues()
	// handle request
	switch res {
	case APILogin:
		{
			if apiReq.Username == "" || apiReq.Password == "" {
				errorResponse(w, ErrInvalidCreds)
				return errors.WithStack(ErrInvalidCreds)
			}
			sendResponse(w, http.StatusOK, &APIResponse{
				Success: true,
			})
			return nil
		}
	case APILogout:
		{
			sendResponse(w, http.StatusOK, &APIResponse{
				Success: true,
			})
			return nil
		}
	}
	errorResponse(w, ErrInvalidResource)
	return errors.WithStack(ErrInvalidResource)
}

func login(w http.ResponseWriter, r *http.Request) {
	request(APILogin, w, r)
}
