package overmind

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose an overmind service.
type Endpoints struct {
	GetHealthEndpoint    endpoint.Endpoint
	GetZerglingsEndpoint endpoint.Endpoint
}

type getHealthResponse struct {
	Health Health `json:"health,omitempty"`
	Err    error  `json:"err,omitempty"`
}

type getZerglingsResponse struct {
	Zerglings []Zergling `json:"zerglings,omitempty"`
	Err       error      `json:"err,omitempty"`
}

func makeGetHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		health, e := s.GetHealth(ctx)
		return getHealthResponse{
			Health: health,
			Err:    e,
		}, nil
	}
}

func makeGetZerglingsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		zerglings, e := s.GetZerglings(ctx)
		return getZerglingsResponse{
			Zerglings: zerglings,
			Err:       e,
		}, nil
	}
}

// MakeServerEndpoints returns an Endpoint struct
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetHealthEndpoint:    makeGetHealthEndpoint(s),
		GetZerglingsEndpoint: makeGetZerglingsEndpoint(s),
	}
}
