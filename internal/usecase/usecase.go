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

func NewUseCase(repo *repository.UserRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) GetUser(c context.Context, id string) (models.User, error) {
	user, err := u.repo.GetUser(c, id)
	return user, err
}

func (u *UseCase) CreateUser(c context.Context, userReq models.UserRequest) (uuid.UUID, error) {
	return u.repo.CreateUser(c, userReq)
}

func (u *UseCase) UpdateUser(c context.Context, id string, userReq models.UserRequest) (models.User, error) {
	user, err := u.repo.UpdateUser(c, id, userReq)
	return user, err
}

func (u *UseCase) DeleteUser(c context.Context, id string) error {
	return u.repo.DeleteUser(c, id)
}
