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

var commandHandlers = map[Command]map[Facing]commandResult{
	"L": map[Facing]commandResult{
		"N": commandResult{"W", locationDelta{}},
		"E": commandResult{"N", locationDelta{}},
		"S": commandResult{"E", locationDelta{}},
		"W": commandResult{"S", locationDelta{}},
	},
	"R": map[Facing]commandResult{
		"N": commandResult{"E", locationDelta{}},
		"E": commandResult{"S", locationDelta{}},
		"S": commandResult{"W", locationDelta{}},
		"W": commandResult{"N", locationDelta{}},
	},
	"M": map[Facing]commandResult{
		"N": commandResult{"", locationDelta{0, 1}},
		"E": commandResult{"", locationDelta{1, 0}},
		"S": commandResult{"", locationDelta{0, -1}},
		"W": commandResult{"", locationDelta{-1, 0}},
	},
}

// Command represents a command sent to a zergling
type Command string

// Location represents the location of a zergling
type Location struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type commandResult struct {
	newFacing     Facing
	locationDelta locationDelta
}

type locationDelta Location

// Facing represents the direction the zergling is facing
type Facing string

// Zergling represents the simplest form of a zerg mutation
type Zergling struct {
	ID             string    `json:"id"`
	Location       Location  `json:"location"`
	Facing         Facing    `json:"facing"`
	CommandHistory []Command `json:"commandHistory,omitempty"`
	Revision       string    `json:"_rev"`
}

func (z *Zergling) receiveCommand(command Command) error {
	commandHandler, ok := commandHandlers[command]
	if !ok {
		return ErrInvalidCommand
	}

	commandResult, ok := commandHandler[z.Facing]
	if !ok {
		return ErrUnknownFacing
	}

	z.Location.X += commandResult.locationDelta.X
	z.Location.Y += commandResult.locationDelta.Y

	if commandResult.newFacing != "" {
		z.Facing = commandResult.newFacing
	}

	return nil
}

// Service is an interface for the Overmind service
type Service interface {
	GetHealth(ctx context.Context) (Health, error)
	GetZerglings(ctx context.Context) ([]Zergling, error)
	PostZerglings(ctx context.Context) (Zergling, error)
	GetZerglingByID(ctx context.Context, id string) (Zergling, error)
	PostZerglingCommand(ctx context.Context, id string, command Command) (Zergling, error)
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
	zergling := Zergling{
		ID:       id,
		Location: Location{0, 0},
		Facing:   "N",
	}
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

func (s *overmindService) PostZerglingCommand(ctx context.Context, id string, command Command) (Zergling, error) {
	zergling, err := s.GetZerglingByID(ctx, id)
	if err != nil {
		return zergling, err
	}

	if zergling.CommandHistory == nil {
		zergling.CommandHistory = []Command{}
	}

	zergling.CommandHistory = append(zergling.CommandHistory, command)
	zergling.receiveCommand(command)
	// TODO: Persist
	return zergling, err
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
