package cache

import (
	"context"
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserProvider := mocks.NewMockUserProvider(ctrl) // создаём мок
			// Говорим что мы ожидаем вызов CreateUser с таким return
			mockUserProvider.EXPECT().
				CreateUser(gomock.Any(), tc.userRequest).
				Return(uuid.New(), nil) // Возвращаем uuid

			t.Logf("initializing cache decorator, ttl: %s\n", tc.cacheTTL)
			cache := NewCacheDecorator(mockUserProvider, tc.cacheTTL)

			t.Log("creating user through cache decorator")
			id, err := cache.CreateUser(context.Background(), tc.userRequest)
			require.NoError(t, err)
			// Проверяем, что uuid не пустой
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserProvider := mocks.NewMockUserProvider(ctrl) // создаём мок
			mockUserProvider.EXPECT().
				GetUser(gomock.Any(), tc.ID.String()).
				Return(&models.User{
					ID:        tc.ID,
					Name:      "Daniel",
					Age:       30,
					Anonymous: false,
				}, nil) // Возвращаем юзера

			t.Logf("initializing cache decorator, ttl: %s\n", tc.cacheTTL)
			cache := NewCacheDecorator(mockUserProvider, tc.cacheTTL)

			t.Log("creating user through cache decorator\n")
			user, err := cache.GetUser(context.Background(), tc.ID.String())
			require.NoError(t, err)

			t.Log("get user from cache\n")
			cached, ok := cache.users[user.ID.String()]
			require.Truef(t, ok, "пользователь должен быть в кэше \n")
			require.Equal(t, tc.ID, cached.user.ID)
		})
	}
}
