package usecase

import (
	"context"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/dankru/Api_gateway_v2/internal/repository"
	"github.com/google/uuid"
)

type UseCase struct {
	repo *repository.UserRepository
}

func NewUseCase(repo *repository.UserRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) GetUser(c context.Context, id string) (models.UserResponse, error) {
	user, err := u.repo.GetUser(c, id)
	return u.mapUserToResponse(user), err
}

func (u *UseCase) CreateUser(c context.Context, input models.UserRequest) (uuid.UUID, error) {
	return u.repo.CreateUser(c, input)
}

func (u *UseCase) UpdateUser(c context.Context, id string, input models.UserRequest) (models.UserResponse, error) {
	user, err := u.repo.UpdateUser(c, id, input)
	return u.mapUserToResponse(user), err
}

func (u *UseCase) DeleteUser(c context.Context, id string) error {
	return u.repo.DeleteUser(c, id)
}

func (u *UseCase) mapUserToResponse(user models.User) models.UserResponse {
	return models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Age:       user.Age,
		Anonymous: user.Anonymous,
	}
}
