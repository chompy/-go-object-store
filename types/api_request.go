package types

// APIRequest defines an API request.
type APIRequest struct {
	IP         string      `json:"-"`
	SessionKey string      `json:"key,omitempty"`
	Username   string      `json:"username,omitempty"`
	Password   string      `json:"password,omitempty"`
	Objects    []APIObject `json:"objects,omitempty"`
	Query      string      `json:"query,omitempty"`
}

// ObjectUIDs return list of object uids in api request.
func (a APIRequest) ObjectUIDs() []string {
	out := make([]string, 0)
	for _, o := range a.Objects {
		uid := o.Object().UID
		if uid != "" {
			out = append(out, uid)
		}
	}
	return out
}
