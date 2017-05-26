package overmind

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose an overmind service.
type Endpoints struct {
	GetHealthEndpoint endpoint.Endpoint
}

type getHealthResponse struct {
	Health Health `json:"health,omitempty"`
	Err    error  `json:"err,omitempty"`
}

// MakeGetHealthEndpoint returns a GetHealth endpoint
func MakeGetHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		health, e := s.GetHealth(ctx)
		return getHealthResponse{
			Health: health,
			Err:    e,
		}, nil
	}
}

// MakeServerEndpoints returns an Endpoint struct
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetHealthEndpoint: MakeGetHealthEndpoint(s),
	}
}
