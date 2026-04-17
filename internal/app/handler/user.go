package handler

import (
	"avanpost-show/internal/user/delivery/http"
	"avanpost-show/internal/user/repository/postgres"
	"avanpost-show/internal/user/usecase"
	"avanpost-show/pkg/publisher"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserHandler(db *pgxpool.Pool, p *publisher.Publisher) Handler {
	return http.NewUserHandler(usecase.NewUseCase(postgres.NewRepository(db), p))
}
