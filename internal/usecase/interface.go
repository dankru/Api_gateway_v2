package usecase

import (
	"context"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/google/uuid"
)

type UserProvider interface {
	GetUser(c context.Context, id string) (models.User, error)
	CreateUser(c context.Context, userReq models.UserRequest) (uuid.UUID, error)
	UpdateUser(c context.Context, id string, userReq models.UserRequest) (models.User, error)
	DeleteUser(c context.Context, id string) error
}
