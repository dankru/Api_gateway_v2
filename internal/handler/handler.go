package handler

import (
	"github.com/dankru/Api_gateway_v2/internal/apperr"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/dankru/Api_gateway_v2/internal/usecase"
	"github.com/dankru/Api_gateway_v2/internal/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type Handler struct {
	userUC usecase.UserProvider
	tracer trace.Tracer
}

func NewHandler(userUC *usecase.UserUsecase, tracer trace.Tracer) *Handler {
	return &Handler{userUC: userUC, tracer: tracer}
}

func (h *Handler) GetUser(ctx *fiber.Ctx) error {
	spanCtx, span := h.tracer.Start(ctx.UserContext(), "Handler.GetUser")
	defer span.End()
	span.SetAttributes(
		attribute.String("id", ctx.Params("id")), // Параметр запроса (id)
	)

	id := ctx.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		span.SetStatus(codes.Error, "invalid uuid")
		return fiber.NewError(http.StatusBadRequest, "invalid uuid")
	}

	user, err := h.userUC.GetUser(spanCtx, id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			log.Err(err).Msgf("user with id: %s not found", id)
			span.SetStatus(codes.Error, "user not found")
			return fiber.NewError(http.StatusNotFound)
		}
		log.Err(err).Msg("failed to get user from db")
		span.SetStatus(codes.Error, "failed to get user")
		return errors.Wrap(err, "failed to get user")
	}

	response := h.mapUserToResponse(user)
	return ctx.JSON(fiber.Map{"data": response})
}

func (h *Handler) CreateUser(ctx *fiber.Ctx) error {
	spanCtx, span := h.tracer.Start(ctx.UserContext(), "Handler.CreateUser")
	defer span.End()

	var userReq models.UserRequest
	if err := ctx.BodyParser(&userReq); err != nil {
		log.Err(err).Msg("failed to parse user input")
		return fiber.NewError(http.StatusBadRequest, "invalid input")
	}

	if err := validation.Validate(userReq); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid input")
	}

	id, err := h.userUC.CreateUser(spanCtx, userReq)
	if err != nil {
		log.Err(err).Msg("failed to create user")
		return errors.Wrap(err, "failed to insert into users")
	}

	return ctx.JSON(fiber.Map{"id": id})
}

func (h *Handler) ReplaceUser(ctx *fiber.Ctx) error {
	spanCtx, span := h.tracer.Start(ctx.UserContext(), "Handler.ReplaceUser")
	defer span.End()
	span.SetAttributes(
		attribute.String("id", ctx.Params("id")), // Параметр запроса (id)
	)

	id := ctx.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	var userReq models.UserRequest
	if err := ctx.BodyParser(&userReq); err != nil {
		log.Err(err).Msg("failed to parse user input")
		return fiber.NewError(http.StatusBadRequest, "invalid input")
	}

	if err := validation.Validate(userReq); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userUC.UpdateUser(spanCtx, id, userReq)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			log.Err(err).Msgf("user with id: %s not found", id)
			return fiber.NewError(http.StatusNotFound)
		}
		log.Err(err).Msgf("failed to update user by id %s", id)
		return errors.Wrap(err, "failed to update user")
	}

	response := h.mapUserToResponse(user)
	return ctx.JSON(fiber.Map{"data": response})
}

func (h *Handler) DeleteUser(ctx *fiber.Ctx) error {
	spanCtx, span := h.tracer.Start(ctx.UserContext(), "Handler.DeleteUser")
	defer span.End()
	span.SetAttributes(
		attribute.String("id", ctx.Params("id")), // Параметр запроса (id)
	)

	id := ctx.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		return fiber.NewError(http.StatusBadRequest, "invalid uuid")
	}

	if err := h.userUC.DeleteUser(spanCtx, id); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			log.Err(err).Msgf("user with id: %s not found", id)
			return fiber.NewError(http.StatusNotFound)
		}
		log.Err(err).Msgf("failed to delete user by id %s", id)
		return errors.Wrap(err, "failed to delete user")
	}

	return ctx.SendStatus(http.StatusNoContent)
}

func (h *Handler) mapUserToResponse(u *models.User) models.UserResponse {
	return models.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Age:       u.Age,
		Anonymous: u.Anonymous,
	}
}
