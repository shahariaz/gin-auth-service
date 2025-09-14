package lib

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenStore interface {
	Blacklist(token string, expiry time.Duration) error
	IsBlacklisted(token string) (bool, error)
}

// type InMemoryTokenStore struct {
// 	tokens map[string]time.Time
// 	mu     sync.RWMutex
// }

// func NewInMemoryTokenStore() *InMemoryTokenStore {
// 	return &InMemoryTokenStore{
// 		tokens: make(map[string]time.Time),
// 	}
// }

// func (s *InMemoryTokenStore) Blacklist(token string, expiry time.Duration) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.tokens[token] = time.Now().Add(expiry)
// 	return nil
// }

// func (s *InMemoryTokenStore) IsBlacklisted(token string) (bool, error) {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	expiry, exists := s.tokens[token]
// 	if !exists {
// 		return false, nil
// 	}
// 	if time.Now().After(expiry) {
// 		delete(s.tokens, token)
// 		return false, nil
// 	}
// 	return true, nil
// }

type RedisTokenStore struct {
	client *redis.Client
}

func NewRedisTokenStore(redisURL string) (*RedisTokenStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisTokenStore{client: client}, nil
}

func (s *RedisTokenStore) Blacklist(token string, expiry time.Duration) error {
	return s.client.Set(context.Background(), "blacklist:"+token, "1", expiry).Err()
}

func (s *RedisTokenStore) IsBlacklisted(token string) (bool, error) {
	exists, err := s.client.Exists(context.Background(), "blacklist:"+token).Result()
	return exists == 1, err
}
