package types

// APIResource defines an action to perform via the API.
type APIResource int

const (
	// APILogin defines login action.
	APILogin APIResource = 1
	// APILogout defines logout action.
	APILogout APIResource = 2
	// APIGet defines get object action.
	APIGet APIResource = 3
	// APISet defines set object action.
	APISet APIResource = 4
	// APIDelete defines delete object action.
	APIDelete APIResource = 5
	// APIQuery defines query object action.
	APIQuery APIResource = 6
)

// Name returns string name for API resource.
func (r APIResource) Name() string {
	switch r {
	case APILogin:
		{
			return "LOGIN"
		}
	case APILogout:
		{
			return "LOGOUT"
		}
	case APIGet:
		{
			return "GET"
		}
	case APISet:
		{
			return "SET"
		}
	case APIDelete:
		{
			return "DELETE"
		}
	case APIQuery:
		{
			return "QUERY"
		}
	}
	return ""
}
