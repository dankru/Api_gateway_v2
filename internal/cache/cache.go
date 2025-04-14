package cache

import (
	"context"
	"github.com/dankru/Api_gateway_v2/config"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/dankru/Api_gateway_v2/internal/repository"
	"github.com/google/uuid"
	"sync"
	"time"
)

type CacheDecorator struct {
	repo repository.UserProvider

	mu    sync.RWMutex
	users map[string]models.WrapUser

	cfg config.Config
}

func NewCacheDecorator(repo repository.UserProvider, cfg config.Config) *CacheDecorator {
	return &CacheDecorator{
		repo:  repo,
		mu:    sync.RWMutex{},
		users: make(map[string]models.WrapUser, 100),
		cfg:   cfg,
	}
}

// Не понимаю как её правильно завершать
func (cache *CacheDecorator) StartCleaner(c context.Context) {
	ticker := time.NewTicker(cache.cfg.CleanerInterval)

	go func() {
		for {
			select {
			case <-c.Done():
				return
			case t := <-ticker.C:
				for id, wrap := range cache.users {
					if wrap.ExpiredAt.Before(t) {
						cache.invalidate(id)
					}
				}
			}
		}
	}()
}

func (cache *CacheDecorator) get(id string) (models.WrapUser, bool) {
	cache.mu.RLock()
	wrap, exists := cache.users[id]
	cache.mu.RUnlock()
	return wrap, exists
}

func (cache *CacheDecorator) set(user models.User, id string) {
	wrap := models.WrapUser{user, time.Now().Add(cache.cfg.CacheTTL)}

	cache.mu.Lock()
	cache.users[id] = wrap
	cache.mu.Unlock()
}

func (cache *CacheDecorator) invalidate(id string) {
	cache.mu.Lock()
	delete(cache.users, id)
	cache.mu.Unlock()
}

func (cache *CacheDecorator) GetUser(c context.Context, id string) (models.User, error) {
	wrap, exists := cache.get(id)
	if exists {
		return wrap.User, nil
	}

	user, err := cache.repo.GetUser(c, id)
	if err != nil {
		return user, err
	}
	cache.set(user, id)

	return user, nil
}

func (cache *CacheDecorator) CreateUser(c context.Context, userReq models.UserRequest) (uuid.UUID, error) {
	return cache.repo.CreateUser(c, userReq)
}

func (cache *CacheDecorator) UpdateUser(c context.Context, id string, userReq models.UserRequest) (models.User, error) {
	user, err := cache.repo.UpdateUser(c, id, userReq)
	if err != nil {
		return user, err
	}
	cache.set(user, id)
	return user, nil
}

func (cache *CacheDecorator) DeleteUser(c context.Context, id string) error {

	if err := cache.repo.DeleteUser(c, id); err != nil {
		return err
	}
	cache.invalidate(id)

	return nil
}
