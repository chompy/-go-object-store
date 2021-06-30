package main

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const indexValueMaxSize = 128

// Object defines a storable object.
type Object struct {
	UID      string                 `json:"uid"`
	Author   string                 `json:"author"`
	Created  time.Time              `json:"created"`
	Modified time.Time              `json:"modified"`
	Data     map[string]interface{} `json:"data"`
}

// NewObject returns a new object.
func NewObject(u *User) *Object {
	uid, err := gonanoid.New()
	if err != nil {
		return nil
	}
	author := ""
	if u != nil {
		author = u.UID
	}
	return &Object{
		UID:      uid,
		Author:   author,
		Created:  time.Now(),
		Modified: time.Now(),
		Data:     make(map[string]interface{}),
	}
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
		Created:  o.Created,
		Modified: o.Modified,
		Data:     indexData,
	}
}
