package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"gitlab.com/contextualcode/go-object-store/store"
	"gitlab.com/contextualcode/go-object-store/types"
)

const testHTTPPort = 31582

func initTestServer() {
	if client != nil {
		return
	}
	c := &store.Config{}
	c.UserGroups = map[string]store.UserGroup{
		"anonymous": store.UserGroup{
			Get:    true,
			Set:    false,
			Update: false,
			Delete: false,
		},
		"admin": store.UserGroup{
			Get:    true,
			Set:    true,
			Update: true,
			Delete: true,
		},
	}
	c.HTTP.Port = testHTTPPort
	go Listen(c)
	time.Sleep(time.Second)
}

func TestHTTPGetNotFound(t *testing.T) {
	initTestServer()
	resp, err := http.Get(
		fmt.Sprintf("http://localhost:%d/get?uid=1", testHTTPPort),
	)
	if err != nil {
		t.Error(err)
		return
	}
	raw, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(raw))
	if resp.StatusCode != http.StatusNotFound {
		t.Error("expected not found status")
		return
	}
}

func TestHTTPGet(t *testing.T) {

	initTestServer()

	// create object
	o := &types.Object{
		Data: map[string]interface{}{
			"test": "hello world",
		},
	}
	client.Set(o, nil)

	// http request
	resp, err := http.Get(
		fmt.Sprintf("http://localhost:%d/get?uid=%s", testHTTPPort, o.UID),
	)
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("unexpected status")
		return
	}

	// read response
	respRaw, _ := ioutil.ReadAll(resp.Body)
	apiResp := types.APIResponse{}
	json.Unmarshal(respRaw, &apiResp)
	if !apiResp.Success {
		t.Error("expected success")
		return
	}

	// compare api obj to local obj
	apiObj := apiResp.Objects[0].Object()
	if apiObj.UID != o.UID || apiObj.Data["test"] != o.Data["test"] {
		t.Error("expected object in api response")
		return
	}

}

func TestHTTPLogin(t *testing.T) {
	initTestServer()

	// create user
	u := &types.User{
		Username: "testuser1",
	}
	u.Username = "testuser1"
	password := "test1234"
	store.SetPassword(password, u)
	u.Groups = []string{"admin"}
	client.SetUser(u)

	// create request
	req := types.APIRequest{
		Username: u.Username,
		Password: password,
	}
	reqJSON, _ := json.Marshal(req)
	reqReader := bytes.NewReader(reqJSON)

	// submit request
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/login", testHTTPPort), "application/json", reqReader)
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
	apiResp := types.APIResponse{}
	rawResp, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(rawResp, &apiResp)
	if apiResp.Key == "" {
		t.Error("expected key in response")
	}

	// test set
	req = types.APIRequest{
		SessionKey: apiResp.Key,
		Objects: []types.APIObject{
			types.APIObject{"test": "hello world"},
		},
	}
	reqJSON, _ = json.Marshal(req)
	reqReader = bytes.NewReader(reqJSON)
	resp, err = http.Post(fmt.Sprintf("http://localhost:%d/set", testHTTPPort), "application/json", reqReader)
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
