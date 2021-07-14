package main

import (
	"fmt"
	"strings"

	"gitlab.com/contextualcode/go-object-store/types"
)

func sanitizeValues(req *types.APIRequest) {
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.Password = strings.TrimSpace(req.Password)
}

func logAPIRequest(req types.APIRequest, res types.APIResource) {
	userIdentity := anonymousUser
	if req.Username != "" {
		userIdentity = req.Username
	} else if req.SessionKey != "" {
		userIdentity = req.SessionKey
	}
	objString := ""
	if len(req.Objects) > 0 {
		objList := req.ObjectUIDs()
		newObjCt := len(req.Objects) - len(objList)
		objString = fmt.Sprintf(" %s", strings.Join(objList, ","))
		if newObjCt > 0 {
			objString += fmt.Sprintf(" + %d new", newObjCt)
		}
	} else if req.Query != "" {
		objString += " " + req.Query
	}
	logInfo(
		fmt.Sprintf("@%s - %s%s", userIdentity, res.Name(), objString),
	)
}
