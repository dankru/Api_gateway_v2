package repository

import (
	"context"
	"github.com/dankru/Api_gateway_v2/config"
	"github.com/dankru/Api_gateway_v2/internal/apperr"
	"github.com/dankru/Api_gateway_v2/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{conn: conn}
}

func (u *UserRepository) GetUser(ctx context.Context, id string) (*models.User, error) {
	tracer := otel.Tracer(config.AppName)
	_, span := tracer.Start(ctx, "UserRepository.GetByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.query", "SELECT id, name, age, anonymous FROM users WHERE id = $1"),
		attribute.String("db.params.id", id),
		attribute.String("db.system", "postgres"),
	)

	userData := &models.User{}

	start := time.Now()
	err := u.conn.QueryRow(
		ctx,
		"SELECT id, name, age, anonymous FROM users WHERE id = $1", id).
		Scan(userData.ID,
			userData.Name,
			userData.Age,
			userData.Anonymous)

	duration := time.Since(start)

	span.SetAttributes(
		attribute.Int64("db.duration_ms", duration.Milliseconds()),
		attribute.Bool("db.success", err == nil),
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}

	return userData, err
}

func (u *UserRepository) CreateUser(ctx context.Context, userReq models.UserRequest) (uuid.UUID, error) {
	tracer := otel.Tracer(config.AppName)
	_, span := tracer.Start(ctx, "UserRepository.CreateUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.query", "INSERT INTO users (name, age, anonymous) VALUES ($1, $2, $3) RETURNING id"),
		attribute.String("db.params.name", userReq.Name),
		attribute.Int("db.params.age", userReq.Age),
		attribute.Bool("db.params.anonymous", userReq.Anonymous),
		attribute.String("db.system", "postgres"),
	)

	var userId uuid.UUID

	start := time.Now()
	err := u.conn.QueryRow(
		ctx,
		"INSERT INTO users (name, age, anonymous) VALUES ($1, $2, $3) RETURNING id",
		userReq.Name,
		userReq.Age,
		userReq.Anonymous,
	).Scan(&userId)

	duration := time.Since(start)

	span.SetAttributes(
		attribute.Int64("db.duration_ms", duration.Milliseconds()),
		attribute.Bool("db.success", err == nil),
	)

	return userId, err
}

func (u *UserRepository) UpdateUser(ctx context.Context, id string, userReq models.UserRequest) (*models.User, error) {
	tracer := otel.Tracer(config.AppName)
	_, span := tracer.Start(ctx, "UserRepository.UpdateUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.query", "UPDATE users SET name = $1, age = $2, anonymous = $3 WHERE id = $4 RETURNING id, name, age, anonymous"),
		attribute.String("db.params.name", userReq.Name),
		attribute.Int("db.params.age", userReq.Age),
		attribute.Bool("db.params.anonymous", userReq.Anonymous),
		attribute.String("db.system", "postgres"),
	)

	userData := &models.User{}

	start := time.Now()
	err := u.conn.QueryRow(
		ctx,
		"UPDATE users SET name = $1, age = $2, anonymous = $3 WHERE id = $4 RETURNING id, name, age, anonymous",
		userReq.Name,
		userReq.Age,
		userReq.Anonymous,
		id).
		Scan(userData.ID, userData.Name, userData.Age, userData.Anonymous)

	duration := time.Since(start)

	span.SetAttributes(
		attribute.Int64("db.duration_ms", duration.Milliseconds()),
		attribute.Bool("db.success", err == nil),
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}

	return userData, err
}

func (u *UserRepository) DeleteUser(ctx context.Context, id string) error {
	tracer := otel.Tracer(config.AppName)
	_, span := tracer.Start(ctx, "UserRepository.DeleteUser")
	defer span.End()
	span.SetAttributes(
		attribute.String("db.query", "DELETE FROM users WHERE id = $1"),
		attribute.String("db.params.id", id),
		attribute.String("db.system", "postgres"),
	)

	start := time.Now()
	result, err := u.conn.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Err(err).Msg("failed to delete user")
		return errors.Wrap(err, "failed to delete user")
	}

	duration := time.Since(start)

	span.SetAttributes(
		attribute.Int64("db.duration_ms", duration.Milliseconds()),
		attribute.Bool("db.success", err == nil),
	)

	if result.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}

	return nil
}
