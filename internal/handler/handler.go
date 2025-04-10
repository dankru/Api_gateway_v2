package handler

import (
	"encoding/json"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/dankru/Api_gateway_v2/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler struct {
	useCase *usecase.UseCase
}

func init() {
}

func NewHandler(useCase *usecase.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}

	user, err := h.useCase.GetUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Err(err).Msgf("user with id: %s not found", id)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
		}
		log.Err(err).Msg("failed to get user from db")
		return errors.Wrap(err, "failed to get user")
	}

	return c.JSON(fiber.Map{"data": user})
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	body := c.Body()

	var input models.UserRequest
	if err := json.Unmarshal(body, &input); err != nil {
		log.Err(err).Msg("failed to unmarshall user input")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid input"})
	}

	if err := input.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid input"})
	}

	id, err := h.useCase.CreateUser(c.Context(), input)
	if err != nil {
		log.Err(err).Msg("failed to create user")
		return errors.Wrap(err, "failed to insert into users")
	}

	return c.JSON(fiber.Map{"id": id})
}

func (h *Handler) ReplaceUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}

	body := c.Body()

	var input models.UserRequest
	if err := json.Unmarshal(body, &input); err != nil {
		log.Err(err).Msg("failed to unmarshall user input")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid input"})

	}

	if err := input.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid input"})
	}

	response, err := h.useCase.UpdateUser(c.Context(), id, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Err(err).Msgf("user with id: %s not found", id)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
		}
		log.Err(err).Msgf("failed to update user by id %s", id)
		return errors.Wrap(err, "failed to update user")
	}
	return c.JSON(fiber.Map{"data": response})
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(err).Msg("validation failed")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}

	if err := h.useCase.DeleteUser(c.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Err(err).Msgf("user with id: %s not found", id)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
		}
		log.Err(err).Msgf("failed to delete user by id %s", id)
		return errors.Wrap(err, "failed to delete user")
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{})
}
