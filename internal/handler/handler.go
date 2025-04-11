package handler

import (
	"encoding/json"
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
	useCase usecase.UserProvider
}

func NewHandler(useCase *usecase.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

// Вопрос: как нам закрыть интерфейсом handler, если fiber удовлетворяет error?
// Я хочу чтобы мы могли явно требовать userResponse
func (h *Handler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		return fiber.NewError(http.StatusBadRequest, "invalid uuid")
	}

	user, err := h.useCase.GetUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			log.Err(err).Msgf("user with id: %s not found", id)
			return fiber.NewError(http.StatusNotFound)
		}
		log.Err(err).Msg("failed to get user from db")
		return errors.Wrap(err, "failed to get user")
	}

	response := h.mapUserToResponse(user)
	return c.JSON(fiber.Map{"data": response})
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	body := c.Body()

	var userReq models.UserRequest
	if err := json.Unmarshal(body, &userReq); err != nil {
		log.Err(err).Msg("failed to unmarshall user input")
		return fiber.NewError(http.StatusBadRequest, "invalid input")
	}

	if err := validation.Validate(userReq); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid input")
	}

	id, err := h.useCase.CreateUser(c.Context(), userReq)
	if err != nil {
		log.Err(err).Msg("failed to create user")
		return errors.Wrap(err, "failed to insert into users")
	}

	return c.JSON(fiber.Map{"id": id})
}

// Вопрос: как нам закрыть интерфейсом handler, если fiber удовлетворяет error?
// Я хочу чтобы мы могли явно требовать userResponse
func (h *Handler) ReplaceUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	body := c.Body()

	var userReq models.UserRequest
	if err := json.Unmarshal(body, &userReq); err != nil {
		log.Err(err).Msg("failed to unmarshall user input")
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if err := validation.Validate(userReq); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	user, err := h.useCase.UpdateUser(c.Context(), id, userReq)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			log.Err(err).Msgf("user with id: %s not found", id)
			return fiber.NewError(http.StatusNotFound)
		}
		log.Err(err).Msgf("failed to update user by id %s", id)
		return errors.Wrap(err, "failed to update user")
	}

	response := h.mapUserToResponse(user)
	return c.JSON(fiber.Map{"data": response})
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		return fiber.NewError(http.StatusBadRequest, "invalid uuid")
	}

	if err := h.useCase.DeleteUser(c.Context(), id); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			log.Err(err).Msgf("user with id: %s not found", id)
			return fiber.NewError(http.StatusNotFound)
		}
		log.Err(err).Msgf("failed to delete user by id %s", id)
		return errors.Wrap(err, "failed to delete user")
	}

	return c.SendStatus(http.StatusNoContent)
}

func (h *Handler) mapUserToResponse(u models.User) models.UserResponse {
	return models.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Age:       u.Age,
		Anonymous: u.Anonymous,
	}
}
