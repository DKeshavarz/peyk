package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DKeshavarz/peyk/internal/domain"
	"github.com/DKeshavarz/peyk/internal/infra/cache"
)

type ConnectionCodeRepository interface {
	Save(ctx context.Context, code *domain.ConnectionCode) error
	FindByCode(ctx context.Context, code string) (*domain.ConnectionCode, error)
	Delete(ctx context.Context, code string) error
}


var _ ConnectionCodeRepository = (*connectionCodeRepository)(nil)

type connectionCodeRepository struct {
	cache *cache.Cache
}

func NewConnectionCodeRepository(cache *cache.Cache) ConnectionCodeRepository {
	return &connectionCodeRepository{
		cache: cache,
	}
}

func (r *connectionCodeRepository) Save(ctx context.Context, code *domain.ConnectionCode) error {
	if code == nil {
		return fmt.Errorf("code cannot be nil")
	}

	ttl := time.Until(code.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("code already expired")
	}

	data, err := json.Marshal(code)
	if err != nil {
		return fmt.Errorf("failed to marshal connection code: %w", err)
	}

	
	if err := r.cache.Set(code.Code, string(data), ttl); err != nil {
		return fmt.Errorf("failed to save code by code: %w", err)
	}

	return nil
}

// FindByCode retrieves a connection code by its code value
func (r *connectionCodeRepository) FindByCode(ctx context.Context, code string) (*domain.ConnectionCode, error) {

	data, err := r.cache.Get(code)
	if err != nil {
		if err == cache.ErrNotFound {
			return nil, fmt.Errorf("code not found: %s", code)
		}
		return nil, fmt.Errorf("failed to get code from cache: %w", err)
	}

	var connectionCode domain.ConnectionCode
	if err := json.Unmarshal([]byte(data), &connectionCode); err != nil {
		return nil, fmt.Errorf("failed to unmarshal connection code: %w", err)
	}

	return &connectionCode, nil
}

func (r *connectionCodeRepository) Delete(ctx context.Context, code string) error {
	_, err := r.FindByCode(ctx, code)
	if err != nil {
		return nil
	}

	
	if err := r.cache.Delete(code); err != nil {
		return fmt.Errorf("failed to delete code: %w", err)
	}

	return nil
}
