package usecase

import (
	"context"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/dankru/Api_gateway_v2/internal/repository"
	"github.com/google/uuid"
)

type UseCase struct {
	repo repository.UserProvider
}

func NewUseCase(repo repository.UserProvider) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) GetUser(ctx context.Context, id string) (models.User, error) {
	user, err := u.repo.GetUser(ctx, id)
	return user, err
}

func (u *UseCase) CreateUser(ctx context.Context, userReq models.UserRequest) (uuid.UUID, error) {
	return u.repo.CreateUser(ctx, userReq)
}

func (u *UseCase) UpdateUser(ctx context.Context, id string, userReq models.UserRequest) (models.User, error) {
	user, err := u.repo.UpdateUser(ctx, id, userReq)
	return user, err
}

func (u *UseCase) DeleteUser(ctx context.Context, id string) error {
	return u.repo.DeleteUser(ctx, id)
}
