package usecase

import (
	"context"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/dankru/Api_gateway_v2/internal/repository"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type UserUsecase struct {
	repo   repository.UserProvider
	tracer trace.Tracer
}

func NewUserUsecase(repo repository.UserProvider, tracer trace.Tracer) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		tracer: tracer,
	}
}

func (u *UserUsecase) GetUser(ctx context.Context, id string) (*models.User, error) {
	ctx, span := u.tracer.Start(ctx, "UserService.GetUser")
	defer span.End()
	user, err := u.repo.GetUser(ctx, id)
	return user, err
}

func (u *UserUsecase) CreateUser(ctx context.Context, userReq models.UserRequest) (uuid.UUID, error) {
	ctx, span := u.tracer.Start(ctx, "UserService.CreateUser")
	defer span.End()
	return u.repo.CreateUser(ctx, userReq)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id string, userReq models.UserRequest) (*models.User, error) {
	ctx, span := u.tracer.Start(ctx, "UserService.UpdateUser")
	defer span.End()
	user, err := u.repo.UpdateUser(ctx, id, userReq)
	return user, err
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	ctx, span := u.tracer.Start(ctx, "UserService.DeleteUser")
	defer span.End()
	return u.repo.DeleteUser(ctx, id)
}
