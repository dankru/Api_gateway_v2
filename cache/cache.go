package cache

import (
	"context"
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

	ttl time.Duration
}

func (cache *CacheDecorator) get(id string) (models.WrapUser, bool) {
	cache.mu.RLock()
	wrap, exists := cache.users[id]
	cache.mu.RUnlock()
	return wrap, exists
}

func (cache *CacheDecorator) set(user models.User, id string) {
	// move ttl to conf
	wrap := models.WrapUser{user, time.Now().Add(time.Minute * 60)}

	cache.mu.Lock()
	cache.users[id] = wrap
	cache.mu.Unlock()
}

func (cache *CacheDecorator) GetUser(c context.Context, id string) (models.User, error) {
	wrap, exists := cache.get(id)
	if exists {
		return wrap.User, nil
	}

	user, err := cache.repo.GetUser(c, id)
	cache.set(user, id)

	return user, err
}

func (cache *CacheDecorator) CreateUser(c context.Context, userReq models.UserRequest) (uuid.UUID, error) {
	return cache.repo.CreateUser(c, userReq)
}

func (cache *CacheDecorator) UpdateUser(c context.Context, id string, userReq models.UserRequest) (models.User, error) {
	return cache.repo.UpdateUser(c, id, userReq)
}

func (cache *CacheDecorator) DeleteUser(c context.Context, id string) error {

	if err := cache.repo.DeleteUser(c, id); err != nil {
		return err
	}

	cache.mu.Lock()
	delete(cache.users, id)
	cache.mu.Unlock()

	return nil
}
