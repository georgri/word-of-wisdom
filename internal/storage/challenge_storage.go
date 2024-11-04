package storage

import (
	"sync"

	"github.com/georgri/word-of-wisdom/internal/domain"
)

type challengeStorage struct {
	mutex      sync.RWMutex
	challenges map[string]*domain.HashcashChallenge
}

// NewStorage returns a new challenge challengeStorage
func NewStorage() *challengeStorage {
	return &challengeStorage{
		challenges: make(map[string]*domain.HashcashChallenge),
	}
}

// GetChallenge returns a challenge by id
func (s *challengeStorage) GetChallenge(id string) (challenge *domain.HashcashChallenge, ok bool) {
	s.mutex.RLock()
	challenge, ok = s.challenges[id]
	s.mutex.RUnlock()
	return
}

// SetChallenge saves a challenge
func (s *challengeStorage) SetChallenge(id string, challenge *domain.HashcashChallenge) {
	s.mutex.Lock()
	s.challenges[id] = challenge
	s.mutex.Unlock()
}

// DeleteChallenge deletes a challenge by id
func (s *challengeStorage) DeleteChallenge(id string) {
	s.mutex.Lock()
	delete(s.challenges, id)
	s.mutex.Unlock()
}
