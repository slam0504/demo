package cache

import (
	"context"
	"encoding/json"
	"time"

	"demo/internal/domain/card"
	"github.com/redis/go-redis/v9"
)

// RedisRepository wraps another card.Repository and caches cards in Redis.
type RedisRepository struct {
	Repo  card.Repository
	Redis *redis.Client
}

// NewRedisRepository creates a Redis-backed caching repository.
func NewRedisRepository(repo card.Repository, addr string) *RedisRepository {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisRepository{Repo: repo, Redis: rdb}
}

func key(id string) string { return "card:" + id }

// Save persists events and updates the cache based on those events.
func (r *RedisRepository) Save(ctx context.Context, events []interface{}) error {
	if err := r.Repo.Save(ctx, events); err != nil {
		return err
	}
	for _, evt := range events {
		switch e := evt.(type) {
		case card.CardCreated:
			c := card.Card(e)
			data, _ := json.Marshal(c)
			r.Redis.Set(ctx, key(e.ID.String()), data, time.Hour)
		case card.CardUpdated:
			c := card.Card(e)
			data, _ := json.Marshal(c)
			r.Redis.Set(ctx, key(e.ID.String()), data, time.Hour)
		}
	}
	return nil
}

// Load first checks Redis and falls back to the underlying repository.
func (r *RedisRepository) Load(ctx context.Context, id string) (*card.Card, error) {
	val, err := r.Redis.Get(ctx, key(id)).Result()
	if err == nil {
		var c card.Card
		if err := json.Unmarshal([]byte(val), &c); err == nil {
			return &c, nil
		}
	}
	c, err := r.Repo.Load(ctx, id)
	if err != nil || c == nil {
		return c, err
	}
	data, _ := json.Marshal(c)
	r.Redis.Set(ctx, key(id), data, time.Hour)
	return c, nil
}

// Search delegates to the underlying repository.
func (r *RedisRepository) Search(ctx context.Context, name string, cost int, faction, category, sub string) ([]*card.Card, error) {
	return r.Repo.Search(ctx, name, cost, faction, category, sub)
}
