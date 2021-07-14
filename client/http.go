package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
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

func requestObj(res types.APIResource, objs []string, key string) ([]*types.Object, error) {
	apiObjs := make([]types.APIObject, 0)
	for _, obj := range objs {
		apiObjs = append(apiObjs, types.APIObject{"_uid": obj})
	}
	req := types.APIRequest{
		SessionKey: key,
		Objects:    apiObjs,
	}
	resp, err := request(res, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !resp.Success {
		return nil, errors.WithStack(errors.WithMessage(ErrResponse, resp.Message))
	}
	returnObjs := make([]*types.Object, 0)
	for _, obj := range resp.Objects {
		returnObjs = append(returnObjs, obj.Object())
	}
	return returnObjs, nil
}
