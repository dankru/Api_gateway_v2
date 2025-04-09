package usecase

import "github.com/dankru/Api_gateway_v2/internal/repository"

type UseCase struct {
	repo *repository.UserRepository
}

func NewUseCase(repo *repository.UserRepository) *UseCase {
	return &UseCase{repo: repo}
}
