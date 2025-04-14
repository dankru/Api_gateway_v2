package usecase

import (
	"context"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/dankru/Api_gateway_v2/internal/repository"
	"github.com/google/uuid"
)

type UserUsecase struct {
	repo repository.UserProvider
}

func NewUserUsecase(repo repository.UserProvider) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUser(ctx context.Context, id string) (models.User, error) {
	user, err := u.repo.GetUser(ctx, id)
	return user, err
}

func (u *UserUsecase) CreateUser(ctx context.Context, userReq models.UserRequest) (uuid.UUID, error) {
	return u.repo.CreateUser(ctx, userReq)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id string, userReq models.UserRequest) (models.User, error) {
	user, err := u.repo.UpdateUser(ctx, id, userReq)
	return user, err
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	return u.repo.DeleteUser(ctx, id)
}
