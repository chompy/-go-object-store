package types

import (
	"encoding/json"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const indexMaxAttrSize = 1024

// Object defines a storable object.
type Object struct {
	UID      string                 `json:"uid"`
	Created  time.Time              `json:"created"`
	Modified time.Time              `json:"modified"`
	Type     string                 `json:"type"`
	Data     map[string]interface{} `json:"data"`
	Access   struct {
		Read  []string `json:"read"`
		Write []string `json:"write"`
	} `json:"access"`
	Indexed bool `json:"-"`
}

// NewObject returns a new object.
func NewObject() *Object {
	uid, err := gonanoid.New()
	if err != nil {
		return nil
	}
	return &Object{
		UID:      uid,
		Created:  time.Now(),
		Modified: time.Now(),
		Data:     make(map[string]interface{}),
		Indexed:  false,
	}
}

// Serialize serializes object data.
func (o *Object) Serialize() ([]byte, error) {
	return json.Marshal(o)
}

// Unserialize unserializes object data.
func (o *Object) Unserialize(data []byte) error {
	return json.Unmarshal(data, o)
}

// Index returns version of object with large data sets removed. Used to index for queries.
func (o *Object) Index() *Object {
	out := &Object{
		UID:      o.UID,
		Created:  o.Created,
		Modified: o.Modified,
		Type:     o.Type,
		Data:     make(map[string]interface{}),
		Indexed:  true,
	}
	for k, v := range o.Data {
		switch v := v.(type) {
		case string:
			{
				if len(v) > indexMaxAttrSize {
					out.Data[k] = ""
					break
				}
				out.Data[k] = v
				break
			}
		case []byte:
			{
				if len(v) > indexMaxAttrSize {
					out.Data[k] = []byte{}
					break
				}
				out.Data[k] = v
				break
			}
		default:
			{
				out.Data[k] = v
				break
			}
		}
	}
	return out
}
