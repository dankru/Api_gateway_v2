package repository

import (
	"context"
	"github.com/dankru/Api_gateway_v2/internal/apperr"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{conn: conn}
}

func (u *UserRepository) GetUser(c context.Context, id string) (models.User, error) {
	var userData models.User

	err := u.conn.QueryRow(
		c,
		"SELECT id, name, age, anonymous FROM users WHERE id = $1", id).
		Scan(&userData.ID,
			&userData.Name,
			&userData.Age,
			&userData.Anonymous)
	return userData, err
}

func (u *UserRepository) CreateUser(c context.Context, userReq models.UserRequest) (uuid.UUID, error) {
	var userId uuid.UUID

	err := u.conn.QueryRow(
		c,
		"INSERT INTO users (name, age, anonymous) VALUES ($1, $2, $3) RETURNING id",
		userReq.Name,
		userReq.Age,
		userReq.Anonymous,
	).Scan(&userId)
	return userId, err
}

func (u *UserRepository) UpdateUser(c context.Context, id string, userReq models.UserRequest) (models.User, error) {
	var userData models.User
	err := u.conn.QueryRow(
		c,
		"UPDATE users SET name = $1, age = $2, anonymous = $3 WHERE id = $4 RETURNING id, name, age, anonymous",
		userReq.Name,
		userReq.Age,
		userReq.Anonymous,
		id).
		Scan(&userData.ID, &userData.Name, &userData.Age, &userData.Anonymous)
	return userData, err
}

func (u *UserRepository) DeleteUser(c context.Context, id string) error {
	result, err := u.conn.Exec(c, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Err(err).Msg("failed to delete user")
		return errors.Wrap(err, "failed to delete user")
	}

	if result.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	return nil
}
