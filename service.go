package overmind

import "context"

// Health represents the health of the service
type Health struct {
	Version string `json:"version"`
}

// Service is an interface for the Overmind service
type Service interface {
	GetHealth(ctx context.Context) (Health, error)
}

type overmindService struct {
}

func (s *overmindService) GetHealth(ctx context.Context) (Health, error) {
	return Health{"1.0.0"}, nil
}

// NewOvermindService constructs a new instance of the Overmind service
func NewOvermindService() Service {
	return &overmindService{}
}
