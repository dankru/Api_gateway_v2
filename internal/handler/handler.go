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
	"net/http"
)

type Handler struct {
	userUC usecase.UserProvider
}

func NewHandler(userUC *usecase.UserUsecase) *Handler {
	return &Handler{userUC: userUC}
}

func (h *Handler) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		return fiber.NewError(http.StatusBadRequest, "invalid uuid")
	}

	user, err := h.userUC.GetUser(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			log.Err(err).Msgf("user with id: %s not found", id)
			return fiber.NewError(http.StatusNotFound)
		}
		log.Err(err).Msg("failed to get user from db")
		return errors.Wrap(err, "failed to get user")
	}

	response := h.mapUserToResponse(user)
	return ctx.JSON(fiber.Map{"data": response})
}

func (h *Handler) CreateUser(ctx *fiber.Ctx) error {

	var userReq models.UserRequest
	if err := ctx.BodyParser(&userReq); err != nil {
		log.Err(err).Msg("failed to parse user input")
		return fiber.NewError(http.StatusBadRequest, "invalid input")
	}

	if err := validation.Validate(userReq); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid input")
	}

	id, err := h.userUC.CreateUser(ctx.Context(), userReq)
	if err != nil {
		log.Err(err).Msg("failed to create user")
		return errors.Wrap(err, "failed to insert into users")
	}

	return ctx.JSON(fiber.Map{"id": id})
}

func (h *Handler) ReplaceUser(ctx *fiber.Ctx) error {
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

	user, err := h.userUC.UpdateUser(ctx.Context(), id, userReq)
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
	id := ctx.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		return fiber.NewError(http.StatusBadRequest, "invalid uuid")
	}

	if err := h.userUC.DeleteUser(ctx.Context(), id); err != nil {
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
