package main

import "time"

// IndexObject defines a indexed object.
type IndexObject struct {
	UID      string                 `json:"uid"`
	Author   string                 `json:"author"`
	Created  time.Time              `json:"created"`
	Modified time.Time              `json:"modified"`
	Data     map[string]interface{} `json:"data"`
}

// QueryMap returns map of queryable data.
func (i *IndexObject) QueryMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["_created"] = i.Created.UTC().Unix()
	out["_modified"] = i.Modified.UTC().Unix()
	out["_author"] = i.Author
	for k, v := range i.Data {
		out[k] = v
	}
	return out
}
