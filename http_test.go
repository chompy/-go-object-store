package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestHTTPGetNotFound(t *testing.T) {
	c := &Config{}
	c.UserGroups = map[string]UserGroup{
		"anonymous": UserGroup{
			Get:    false,
			Set:    false,
			Update: false,
			Delete: false,
		},
	}
	c.HTTP.Port = 31582
	go listen(c)
	time.Sleep(time.Second)
	resp, err := http.Get(
		fmt.Sprintf("http://localhost:%d/get?uid=1", c.HTTP.Port),
	)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("expected not found status")
	}
}

func TestHTTPGet(t *testing.T) {
	c := &Config{}
	c.UserGroups = map[string]UserGroup{
		"anonymous": UserGroup{
			Get:    true,
			Set:    false,
			Update: false,
			Delete: false,
		},
	}
	c.HTTP.Port = 31582
	go listen(c)
	time.Sleep(time.Second)

	// create object
	o := &Object{
		Data: map[string]interface{}{
			"test": "hello world",
		},
	}
	store.Set(o, nil)

	// http request
	resp, err := http.Get(
		fmt.Sprintf("http://localhost:%d/get?uid=%s", c.HTTP.Port, o.UID),
	)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("unexpected status")
	}

	// read response
	respRaw, _ := ioutil.ReadAll(resp.Body)
	apiResp := APIResponse{}
	json.Unmarshal(respRaw, &apiResp)
	if !apiResp.Success {
		t.Error("expected success")
	}

	// compare api obj to local obj
	apiObj := apiResp.Objects[0].Object()
	if apiObj.UID != o.UID || apiObj.Data["test"] != o.Data["test"] {
		t.Error("expected object in api response")
	}

}

func TestHTTPLogin(t *testing.T) {
	c := &Config{}
	c.UserGroups = map[string]UserGroup{
		"anonymous": UserGroup{
			Get:    "level = 1",
			Set:    false,
			Update: false,
			Delete: false,
		},
		"admin": UserGroup{
			Get:    true,
			Set:    true,
			Update: true,
			Delete: true,
		},
	}
	c.HTTP.Port = 31582
	go listen(c)
	time.Sleep(time.Second)

	// create user
	u := NewUser()
	u.Username = "testuser1"
	password := "test1234"
	u.SetPassword(password)
	u.Groups = []string{"admin"}
	store.SetUser(u)

	// create request
	req := APIRequest{
		Username: u.Username,
		Password: password,
	}
	reqJSON, _ := json.Marshal(req)
	reqReader := bytes.NewReader(reqJSON)

	// submit request
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/login", c.HTTP.Port), "application/json", reqReader)
	if err != nil {
		t.Error(err)
		return
	}

	// check response
	if resp.StatusCode != http.StatusOK {
		t.Error("expected ok status")
		return
	}

	// read api response, check key
	apiResp := APIResponse{}
	rawResp, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rawResp, &apiResp)
	if apiResp.Key == "" {
		t.Error("expected key in response")
	}

	// test set
	req = APIRequest{
		SessionKey: apiResp.Key,
		Objects: []APIObject{
			APIObject{"test": "hello world"},
		},
	}
	reqJSON, _ = json.Marshal(req)
	reqReader = bytes.NewReader(reqJSON)
	resp, err = http.Post(fmt.Sprintf("http://localhost:%d/set", c.HTTP.Port), "application/json", reqReader)
	if err != nil {
		t.Error(err)
		return
	}
	rawResp, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(rawResp, &apiResp)
	if apiResp.Key == "" {
		t.Error("expected key in response")
	}
	if !apiResp.Success {
		t.Error("expected success response")
	}
	if apiResp.Objects[0].Object().UID == "" {
		t.Error("expected response object to have uid")
	}
	if apiResp.Objects[0].Object().Data["test"] != "hello world" {
		t.Error("unexpected value in response object")
	}

}
