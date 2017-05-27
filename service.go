package overmind

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	couchdb "github.com/fjl/go-couchdb"
	"github.com/google/uuid"
)

// Health represents the health of the service
type Health struct {
	Version string `json:"version"`
	Brain   string `json:"brain"`
}

type Command string

// Zergling represents the simplest form of a zerg mutation
type Zergling struct {
	ID             string    `json:"id"`
	CommandHistory []Command `json:"commandHistory,omitempty"`
}

// Service is an interface for the Overmind service
type Service interface {
	GetHealth(ctx context.Context) (Health, error)
	GetZerglings(ctx context.Context) ([]Zergling, error)
	PostZerglings(ctx context.Context) (Zergling, error)
	GetZerglingByID(ctx context.Context, id string) (Zergling, error)
}

type overmindService struct {
	brain *couchdb.Client
}

func (s *overmindService) GetHealth(ctx context.Context) (Health, error) {
	var brainStatus string
	err := s.brain.Ping()
	if err != nil {
		brainStatus = "damaged"
	} else {
		brainStatus = "ok"
	}
	return Health{Version, brainStatus}, nil
}

type allDocsResult struct {
	TotalRows int `json:"total_rows"`
	Offset    int
	Rows      []Zergling
}

func (s *overmindService) GetZerglings(ctx context.Context) ([]Zergling, error) {
	db := s.brain.DB("zerglings")
	if db == nil {
		return nil, ErrDatabaseNotFound
	}

	var result allDocsResult
	err := db.AllDocs(&result, couchdb.Options{})
	if err != nil {
		return nil, err
	}

	return result.Rows, nil
}

func (s *overmindService) PostZerglings(ctx context.Context) (Zergling, error) {
	id := uuid.New().String()
	db := s.brain.DB("zerglings")
	zergling := Zergling{ID: id}
	_, err := db.Put(id, zergling, "")
	if err != nil {
		return Zergling{}, err
	}
	return zergling, nil
}

func (s *overmindService) GetZerglingByID(ctx context.Context, id string) (Zergling, error) {
	db := s.brain.DB("zerglings")
	var zergling Zergling
	err := db.Get(id, &zergling, couchdb.Options{})
	if err != nil {
		return Zergling{}, err
	}
	return zergling, nil
}

func newCouchDBClient() *couchdb.Client {
	couchDBUserName := os.Getenv("COUCHDB_USERNAME")
	couchDBPassword := os.Getenv("COUCHDB_PASSWORD")
	couchDBServiceHost := os.Getenv("COUCHDB_SERVICE_HOST")
	couchDBServicePort := os.Getenv("COUCHDB_SERVICE_PORT")

	Logger.Log(
		"COUCHDB_USERNAME", couchDBUserName,
		"COUCHDB_PASSWORD", "***",
		"COUCHDB_SERVICE_HOST", couchDBServiceHost,
		"COUCHDB_SERVICE_PORT", couchDBServicePort)

	rawurl := fmt.Sprintf("http://%s:%s@%s:%s", couchDBUserName, couchDBPassword, couchDBServiceHost, couchDBServicePort)
	client, err := couchdb.NewClient(rawurl, &http.Transport{
		ResponseHeaderTimeout: 1 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return client
}

// NewOvermindService constructs a new instance of the Overmind service
func NewOvermindService() Service {
	return &overmindService{brain: newCouchDBClient()}
}
