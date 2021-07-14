package types

import "time"

// IndexObject defines a indexed object.
type IndexObject struct {
	UID      string                 `json:"uid"`
	Author   string                 `json:"author"`
	Modifier string                 `json:"modifier"`
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
	out["_modifier"] = i.Modified
	for k, v := range i.Data {
		out[k] = v
	}
	return out
}

// API converts index object to API object.
func (i *IndexObject) API() APIObject {
	out := make(APIObject)
	out["_uid"] = i.UID
	out["_author"] = i.Author
	out["_modifier"] = i.Modifier
	out["_created"] = i.Created.Format(time.RFC3339)
	out["_modified"] = i.Modified.Format(time.RFC3339)
	for k, v := range i.Data {
		out[k] = v
	}
	return out
}
