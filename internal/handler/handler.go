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
		log.Err(errors.Wrap(err, "invalid uuid provided")).Msg("validation failed")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}

	user, err := h.useCase.GetUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Err(errors.Wrap(err, "user not found")).Msgf("user with id: %s not found", id)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
		}
		log.Err(errors.Wrap(err, "db get user query failed")).Msg("Failed to get user from db")
		return err
	}

	return c.JSON(fiber.Map{"data": user})
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	body := c.Body()

	var input models.UserRequest
	if err := json.Unmarshal(body, &input); err != nil {
		log.Err(errors.Wrap(err, "failed to unmarshall request body into UserRequest struct")).Msg("failed to unmarshall user input")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid input"})
	}

	if err := input.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid input"})
	}

	id, err := h.useCase.CreateUser(c.Context(), input)
	if err != nil {
		log.Err(errors.Wrap(err, "failed to insert into users")).Msg("failed to create user")
		return err
	}

	return c.JSON(fiber.Map{"id": id})
}

func (h *Handler) ReplaceUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(errors.Wrap(err, "invalid uuid provided")).Msg("validation failed")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}

	body := c.Body()

	var input models.UserRequest
	if err := json.Unmarshal(body, &input); err != nil {
		log.Err(errors.Wrap(err, "input inmarshalling failed")).Msg("failed to unmarshall user input")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid input"})

	}

	if err := input.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid input"})
	}

	response, err := h.useCase.UpdateUser(c.Context(), id, input)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Err(errors.Wrap(err, "user not found")).Msgf("user with id: %s not found", id)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
		}
		log.Err(errors.Wrap(err, "failed to update user")).Msgf("failed to update user by id %s", id)
		return err
	}
	return c.JSON(fiber.Map{"data": response})
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		log.Err(errors.Wrap(err, "invalid uuid provided")).Msg("validation failed")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "invalid uuid"})
	}

	if err := h.useCase.DeleteUser(c.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Err(errors.Wrap(err, "user not found")).Msgf("user with id: %s not found", id)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "user not found"})
		}
		log.Err(errors.Wrap(err, "failed to delete user")).Msgf("failed to delete user by id %s", id)
		return err
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{})
}
