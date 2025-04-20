package cache

import (
	"context"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/dankru/Api_gateway_v2/internal/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
	"unsafe"
)

type wrapUser struct {
	user      *models.User
	expiredAt time.Time
}

type CacheDecorator struct {
	repo repository.UserProvider

	mu    sync.RWMutex
	users map[string]wrapUser

	sizeBytes       int
	elementCount    int
	cacheTTL        time.Duration
	cleanerInterval time.Duration
}

func NewCacheDecorator(repo repository.UserProvider, cacheTTL time.Duration) *CacheDecorator {
	return &CacheDecorator{
		repo:     repo,
		mu:       sync.RWMutex{},
		users:    make(map[string]wrapUser, 100),
		cacheTTL: cacheTTL,
	}
}

func (cache *CacheDecorator) StartCleaner(ctx context.Context, cleanerInterval time.Duration) {
	ticker := time.NewTicker(cleanerInterval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("cache cleaner shutting down...")
				return
			case t := <-ticker.C:
				cache.invalidateExpired(t)
			}
		}
	}()
}

func (cache *CacheDecorator) invalidateExpired(t time.Time) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	for id, wrap := range cache.users {
		if wrap.expiredAt.Before(t) {
			log.Info().Msgf("invalidating expired user: %s", wrap.user.ID)
			cache.decrementCacheMetrics(wrap)
			delete(cache.users, id)
		}
	}
}

func (cache *CacheDecorator) get(id string) (wrapUser, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	wrap, exists := cache.users[id]
	return wrap, exists
}

func (cache *CacheDecorator) set(user *models.User, id string) {
	wrap := wrapUser{user, time.Now().Add(cache.cacheTTL)}

	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.users[id] = wrap
	cache.incrementCacheMetrics(wrap)
}

func (cache *CacheDecorator) renewExpiredAt(id string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if wrap, exists := cache.users[id]; exists {
		wrap.expiredAt = time.Now().Add(cache.cacheTTL)
		cache.users[id] = wrap
	}
}

func (cache *CacheDecorator) GetUser(ctx context.Context, id string) (*models.User, error) {
	wrap, exists := cache.get(id)
	if exists {
		cache.renewExpiredAt(id)
		return wrap.user, nil
	}

	user, err := cache.repo.GetUser(ctx, id)
	if err != nil {
		return user, err
	}
	cache.set(user, id)

	return user, nil
}

func (cache *CacheDecorator) CreateUser(ctx context.Context, userReq models.UserRequest) (uuid.UUID, error) {
	return cache.repo.CreateUser(ctx, userReq)
}

func (cache *CacheDecorator) UpdateUser(ctx context.Context, id string, userReq models.UserRequest) (*models.User, error) {
	user, err := cache.repo.UpdateUser(ctx, id, userReq)
	if err != nil {
		return user, err
	}
	cache.set(user, id)
	return user, nil
}

func (cache *CacheDecorator) DeleteUser(ctx context.Context, id string) error {

	if err := cache.repo.DeleteUser(ctx, id); err != nil {
		return err
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.decrementCacheMetrics(cache.users[id])
	delete(cache.users, id)

	return nil
}

func (cache *CacheDecorator) incrementCacheMetrics(u wrapUser) {
	cache.elementCount++
	cache.sizeBytes += cache.getWrapUserSize(u)
}

func (cache *CacheDecorator) decrementCacheMetrics(u wrapUser) {
	cache.elementCount--
	cache.sizeBytes -= cache.getWrapUserSize(u)
}

func (cache *CacheDecorator) getWrapUserSize(u wrapUser) int {
	size := int(unsafe.Sizeof(u)) + int(unsafe.Sizeof(u.user.ID)) + int(unsafe.Sizeof(u.user.Name)) + int(unsafe.Sizeof(u.user.Age)) + int(unsafe.Sizeof(u.user.Anonymous)) + int(unsafe.Sizeof(u.user.PasswordHash))
	return size
}

func (cache *CacheDecorator) ElementCount() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.elementCount
}

func (cache *CacheDecorator) SizeBytes() int {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	return cache.sizeBytes
}
