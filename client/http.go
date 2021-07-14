package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.com/contextualcode/go-object-store/types"
)

var APIURLPrefix = "http://locahost"

func request(res types.APIResource, req types.APIRequest) (types.APIResponse, error) {
	// determine endpoint
	endpoint := ""
	resp := types.APIResponse{}
	switch res {
	case types.APILogin:
		{
			endpoint = APIURLPrefix + "/login"
			break
		}
	case types.APIGet:
		{
			endpoint = APIURLPrefix + "/get"
			break
		}
	case types.APISet:
		{
			endpoint = APIURLPrefix + "/set"
			break
		}
	case types.APIDelete:
		{
			endpoint = APIURLPrefix + "/delete"
			break
		}
	case types.APIQuery:
		{
			endpoint = APIURLPrefix + "/query"
			break
		}
	}
	// encode request to json
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}
	httpResp, err := http.Post(endpoint, "application/json", bytes.NewReader(reqJSON))
	if err != nil {
		return resp, err
	}
	defer httpResp.Body.Close()
	httpRespJSON, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(httpRespJSON, &resp); err != nil {
		return resp, err
	}
	return resp, err
}
