package types

// APIResource defines an action to perform via the API.
type APIResource int

const (
	// APILogin defines login action.
	APILogin APIResource = 1
	// APIGet defines get object action.
	APIGet APIResource = 2
	// APISet defines set object action.
	APISet APIResource = 3
	// APIDelete defines delete object action.
	APIDelete APIResource = 4
	// APIQuery defines query object action.
	APIQuery APIResource = 5
)

// Name returns string name for API resource.
func (r APIResource) Name() string {
	switch r {
	case APILogin:
		{
			return "LOGIN"
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
