package handler

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

type user struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Anonymous bool      `json:"anonymous"`
}

type userInput struct {
	Name      string `json:"name" validate:"required,gte=2"`
	Age       int    `json:"age" validate:"required"`
	Anonymous bool   `json:"anonymous"`
}

type Handler struct {
	Connection *pgxpool.Pool
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func NewHandler(connection *pgxpool.Pool) *Handler {
	return &Handler{Connection: connection}
}

func (h *Handler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var userData user
	err := h.Connection.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE id = $1", id).
		Scan(&userData.ID,
			&userData.Name,
			&userData.Age,
			&userData.Anonymous)

	if err != nil {
		if err == pgx.ErrNoRows {
			//TODO: Log
			c.Status(http.StatusNotFound)
			return nil
		}
		//TODO: Log
		return err
	}

	err = c.JSON(userData)
	if err != nil {
		//TODO: Log
		return err
	}
	return nil
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	body := c.Body()

	var input userInput
	err := json.Unmarshal(body, &input)
	if err != nil {
		//TODO: Log
		c.Status(http.StatusInternalServerError)
	}

	if err := input.Validate(); err != nil {
		c.Status(http.StatusBadRequest)
		err := c.SendString(err.Error())
		if err != nil {
			//TODO: Log
			return err
		}
		return nil
	}

	var userId uuid.UUID
	err = h.Connection.QueryRow(
		context.Background(),
		"INSERT INTO users (name, age, anonymous) VALUES ($1, $2, $3) RETURNING id",
		input.Name,
		input.Age,
		input.Anonymous,
	).Scan(&userId)

	if err != nil {
		//TODO: Log
		return err
	}

	err = c.Send([]byte(userId.String()))
	if err != nil {
		//TODO: Log
		return err
	}

	return nil
}

func (h *Handler) ReplaceUser(c *fiber.Ctx) error {
	id := c.Params("id")

	body := c.Body()

	var input userInput
	err := json.Unmarshal(body, &input)
	if err != nil {
		//TODO: Log
		c.Status(http.StatusInternalServerError)
		return err
	}

	if err := input.Validate(); err != nil {
		c.Status(http.StatusBadRequest)
		err := c.SendString(err.Error())
		if err != nil {
			//TODO: Log
			return err
		}
		return nil
	}

	var userData user
	err = h.Connection.QueryRow(
		context.Background(),
		"UPDATE users SET name = $1, age = $2, anonymous = $3 WHERE id = $4 RETURNING *",
		input.Name,
		input.Age,
		input.Anonymous,
		id,
	).Scan(&userData.ID, &userData.Name, &userData.Age, &userData.Anonymous)
	if err != nil {
		//TODO: Log
		return err
	}

	err = c.JSON(userData)
	if err != nil {
		//TODO: Log
		return err
	}

	return nil
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := h.Connection.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		//TODO: Log
		return err
	}
	if result.RowsAffected() == 0 {
		//TODO: Log
		c.Status(http.StatusNotFound)
		return nil
	}
	c.Status(http.StatusNoContent)
	return nil
}

func (i userInput) Validate() error {
	return validate.Struct(i)
}
