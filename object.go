package main

import (
	"time"
)

const indexValueMaxSize = 128

// Object defines a storable object.
type Object struct {
	UID      string                 `json:"uid"`
	Author   string                 `json:"author"`
	Modifier string                 `json:"modifier"`
	Created  time.Time              `json:"created"`
	Modified time.Time              `json:"modified"`
	Data     map[string]interface{} `json:"data"`
}

// Index returns version of object with large data sets removed. Used to index for queries.
func (o *Object) Index() *IndexObject {
	indexData := make(map[string]interface{})
	for k, v := range o.Data {
		if len(k) > indexValueMaxSize {
			continue
		}
		switch v := v.(type) {
		case string:
			{
				indexData[k] = v
				if len(v) > indexValueMaxSize {
					indexData[k] = v[0:indexValueMaxSize]
				}
				break
			}
		case bool:
			{
				indexData[k] = v
				break
			}
		case int:
			{
				indexData[k] = float64(v)
				break
			}
		case float32:
			{
				indexData[k] = float64(v)
				break
			}
		case float64:
			{
				indexData[k] = v
				break
			}
		}
	}
	return &IndexObject{
		UID:      o.UID,
		Author:   o.Author,
		Created:  o.Created,
		Modified: o.Modified,
		Data:     indexData,
	}
}

// API converts object to API object.
func (o *Object) API() APIObject {
	out := make(APIObject)
	out["_uid"] = o.UID
	out["_author"] = o.Author
	out["_modifier"] = o.Modifier
	out["_created"] = o.Created.Format(time.RFC3339)
	out["_modifiered"] = o.Modified.Format(time.RFC3339)
	for k, v := range o.Data {
		out[k] = v
	}
	return out
}
