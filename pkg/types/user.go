package types

type User struct {
	UID      string                 `bson:"uid"`
	Username string                 `bson:"username"`
	Password []byte                 `bson:"password"`
	Access   []string               `bson:"access"`
	Data     map[string]interface{} `bson:"data"`
}
