package cache

import (
	"github.com/dankru/Api_gateway_v2/internal/mocks"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCacheDecorator_CreateUser(t *testing.T) {
	testCases := []struct {
		name        string
		cacheTTL    time.Duration
		userRequest models.UserRequest
	}{
		{
			name:     "валидный тест на создание юзера",
			cacheTTL: time.Second,
			userRequest: models.UserRequest{
				Name:      "Дмитрий",
				Age:       20,
				Anonymous: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserProvider := newMockUserProvider(t)
			mockUserProvider.EXPECT().
				CreateUser(gomock.Any(), tc.userRequest).
				Return(uuid.New(), nil)

			t.Logf("initializing cache decorator, ttl: %s\n", tc.cacheTTL)
			cache := NewCacheDecorator(mockUserProvider, tc.cacheTTL)

			t.Log("creating user through cache decorator")
			id, err := cache.CreateUser(t.Context(), tc.userRequest)
			t.Log("check err")
			require.NoError(t, err)
			t.Log("check uuid not nil")
			require.NotEqual(t, uuid.Nil, id, "UUID не должен быть пустым (nil)\n")
		})
	}
}

func TestCacheDecorator_GetUser(t *testing.T) {
	testCases := []struct {
		name     string
		cacheTTL time.Duration
		ID       uuid.UUID
	}{
		{
			name:     "валидный UUID тест",
			cacheTTL: time.Second,
			ID:       uuid.New(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			cache, mockUserProvider := setupCacheAndMock(t, tc.cacheTTL)

			expectGetUserProvider(mockUserProvider, tc.ID)

			user := getUserAndCheckNoError(t, cache, tc.ID)

			assertUserCached(t, cache, tc.ID)

			cleanAndAssertUserNotCached(cache, t, tc.cacheTTL, user.ID)
		})
	}
}

func newMockUserProvider(t *testing.T) *mocks.MockUserProvider {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := mocks.NewMockUserProvider(ctrl)
	return mockUserProvider
}

func setupCacheAndMock(t *testing.T, ttl time.Duration) (*CacheDecorator, *mocks.MockUserProvider) {
	mockUserProvider := newMockUserProvider(t)
	cache := NewCacheDecorator(mockUserProvider, ttl)

	return cache, mockUserProvider
}

func expectGetUserProvider(mock *mocks.MockUserProvider, id uuid.UUID) {
	mock.EXPECT().
		GetUser(gomock.Any(), id.String()).
		Return(&models.User{
			ID:        id,
			Name:      "Daniel",
			Age:       30,
			Anonymous: false,
		}, nil)
}

func getUserAndCheckNoError(t *testing.T, cache *CacheDecorator, id uuid.UUID) *models.User {
	user, err := cache.GetUser(t.Context(), id.String())
	require.NoError(t, err)
	return user
}

func assertUserCached(t *testing.T, cache *CacheDecorator, id uuid.UUID) {
	cached, ok := cache.users[id.String()]
	require.True(t, ok, "expected user ID %s to be cached", id)
	require.Equal(t, id, cached.user.ID)
}

func cleanAndAssertUserNotCached(cache *CacheDecorator, t *testing.T, ttl time.Duration, id uuid.UUID) {
	cache.StartCleaner(t.Context(), ttl)

	time.Sleep(time.Second * 2)

	t.Log("check user invalidated\n")
	_, ok := cache.users[id.String()]
	require.False(t, ok, "user must NOT be in cache \n")
}
