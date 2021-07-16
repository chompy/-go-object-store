package http

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/go-object-store/store"
)

var (
	ErrAPIInvalidMethod = errors.New("method not supported")
	ErrInvalidSession   = errors.New("invalid session key")
	ErrInvalidResource  = errors.New("invalid resource")
	ErrEmptyReponse     = errors.New("empty response")
)

func errHTTPResponseCode(err error) int {
	switch errors.Cause(err) {
	case store.ErrNotFound:
		{
			return http.StatusNotFound
		}
	case store.ErrInvalidArg, store.ErrObjectNotSpecified:
		{
			return http.StatusBadRequest
		}
	case store.ErrInvalidCreds, ErrInvalidSession:
		{
			return http.StatusUnauthorized
		}
	case store.ErrPermission:
		{
			return http.StatusUnauthorized
		}
	case ErrAPIInvalidMethod:
		{
			return http.StatusMethodNotAllowed
		}
	}
	return http.StatusInternalServerError
}
