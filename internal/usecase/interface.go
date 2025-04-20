package usecase

import (
	"context"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/google/uuid"
)

type UserProvider interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, userReq models.UserRequest) (uuid.UUID, error)
	UpdateUser(ctx context.Context, id string, userReq models.UserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}
