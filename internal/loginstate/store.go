package loginstate

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

const (
	tokenBytes      = 16
	cleanupInterval = 5 * time.Minute
	storeTTL        = 15 * time.Minute
)

type entry struct {
	challenge string
	expiresAt time.Time
}

type store struct {
	mu    sync.Mutex
	items map[string]entry
	ttl   time.Duration
}

var challengeStore = newStore(storeTTL)

func SaveChallenge(challenge string) (string, error) {
	return challengeStore.save(challenge)
}

func ConsumeChallenge(token string) (string, bool) {
	return challengeStore.consume(token)
}

func newStore(ttl time.Duration) *store {
	s := &store{
		items: make(map[string]entry),
		ttl:   ttl,
	}
	go s.cleanupLoop()
	return s
}

func (s *store) save(challenge string) (string, error) {
	for i := 0; i < 5; i++ {
		token, err := randomToken()
		if err != nil {
			return "", err
		}
		now := time.Now()
		s.mu.Lock()
		if _, exists := s.items[token]; !exists {
			s.items[token] = entry{
				challenge: challenge,
				expiresAt: now.Add(s.ttl),
			}
			s.mu.Unlock()
			return token, nil
		}
		s.mu.Unlock()
	}
	return "", errors.New("failed to allocate login state token")
}

func (s *store) consume(token string) (string, bool) {
	s.mu.Lock()
	entry, ok := s.items[token]
	if ok {
		delete(s.items, token)
	}
	s.mu.Unlock()
	if !ok || time.Now().After(entry.expiresAt) {
		return "", false
	}
	return entry.challenge, true
}

func (s *store) cleanupLoop() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		s.cleanup()
	}
}

func (s *store) cleanup() {
	now := time.Now()
	s.mu.Lock()
	for token, entry := range s.items {
		if now.After(entry.expiresAt) {
			delete(s.items, token)
		}
	}
	s.mu.Unlock()
}

func randomToken() (string, error) {
	buf := make([]byte, tokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
