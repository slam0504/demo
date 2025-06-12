package cache

import (
	"context"
	"errors"
	"testing"

	"demo/internal/domain/card"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type mockRepo struct {
	SaveFn func(ctx context.Context, evts []interface{}) error
	LoadFn func(ctx context.Context, id string) (*card.Card, error)
}

func (m *mockRepo) Save(ctx context.Context, evts []interface{}) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, evts)
	}
	return nil
}
func (m *mockRepo) Load(ctx context.Context, id string) (*card.Card, error) {
	if m.LoadFn != nil {
		return m.LoadFn(ctx, id)
	}
	return nil, nil
}
func (m *mockRepo) Search(ctx context.Context, name string, cost int, faction, category, sub string) ([]*card.Card, error) {
	return nil, nil
}

func TestRedisRepoSaveError(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	repo := &mockRepo{SaveFn: func(ctx context.Context, evts []interface{}) error { return errors.New("fail") }}
	r := &RedisRepository{Repo: repo, Redis: rdb}
	err := r.Save(context.Background(), []interface{}{card.CardCreated{ID: uuid.New()}})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRedisRepoLoadCacheMiss(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	repo := &mockRepo{LoadFn: func(ctx context.Context, id string) (*card.Card, error) { return &card.Card{Name: "N"}, nil }}
	r := &RedisRepository{Repo: repo, Redis: rdb}
	c, err := r.Load(context.Background(), "x")
	if err != nil || c == nil || c.Name != "N" {
		t.Fatalf("unexpected %v %v", c, err)
	}
}

func TestRedisRepoLoadCached(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	repo := &mockRepo{LoadFn: func(ctx context.Context, id string) (*card.Card, error) { return &card.Card{Name: "N"}, nil }}
	r := &RedisRepository{Repo: repo, Redis: rdb}
	// prime cache
	_ = r.Save(context.Background(), []interface{}{card.CardCreated{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "N"}})
	// manually set invalid json to ensure it falls back to repo
	rdb.Set(context.Background(), key("x"), "bad", 0)
	c, err := r.Load(context.Background(), "x")
	if err != nil || c == nil || c.Name != "N" {
		t.Fatalf("unexpected %v %v", c, err)
	}
}
