package scoresvcdemo

import (
	"errors"
	"sync"
	"golang.org/x/net/context"
)

// Score Service
type Service interface {
	PostScore(ctx context.Context, s Score) (Score, error)
	GetScore(ctx context.Context, id string) (Score, error)
}

// Data model
type Score struct {
	Id string `json:"id"`
	Value int `json:"highscore"`
}

// in-memory mock persistence
type inmemService struct {
	sync.RWMutex
	data map[string]Score
}

// Error for missing score in data map
var ErrNotFound = errors.New("not found")

func NewInmemService() Service {
	return &inmemService {
		data: map[string]Score{},
	}
}

// Implements Service.PostScore
func (svc *inmemService) PostScore(ctx context.Context, s Score) (Score, error) {
	svc.Lock()
	defer svc.Unlock()

	score, ok := svc.data[s.Id]
	if ok && s.Value > score.Value {
		svc.data[s.Id].Value = s.Value
	} else if !ok {
		svc.data[s.Id] = &Score{ Id: s.Id, Value: s.Value }
	}

	return svc.data[s.Id], nil
}

// Implements Service.GetScore
func (svc *inmemService) GetScore(ctx context.Context, id string) (Score, error) {
	svc.RLock()
	defer svc.RUnlock()
	score, ok := svc.data[id]
	if !ok {
		return Score{}, ErrNotFound
	}
	return score, nil
}