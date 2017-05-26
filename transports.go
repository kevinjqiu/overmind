package overmind

import (
	"net/http"

	"github.com/go-kit/kit/log"
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {

}
