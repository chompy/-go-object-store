package main

import "time"

// IndexObject defines a indexed object.
type IndexObject struct {
	UID      string                 `json:"uid"`
	Created  time.Time              `json:"created"`
	Modified time.Time              `json:"modified"`
	Data     map[string]interface{} `json:"data"`
}

// Match checks if given value matches tag.
func (i IndexObject) Match(v string) bool {

	return false
}
