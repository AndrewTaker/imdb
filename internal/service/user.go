package service

import (
	"context"
	"fmt"
	"imdb/internal/repository"
	"imdb/internal/security"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	r *repository.UserRepository
	l *slog.Logger
}

func NewUserService(r *repository.UserRepository, l *slog.Logger) *UserService {
	return &UserService{r: r, l: l}
}

func (s *UserService) Create(ctx context.Context, email, password string) error {
	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	return s.r.Create(ctx, email, passwordHash)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*repository.User, error) {
	return s.r.GetByEmail(ctx, email)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*repository.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	return s.r.GetByID(ctx, objID)
}

func (s *UserService) GetAll(
	ctx context.Context,
	pag repository.PaginationOptions,
	sort []repository.SortOptions,
) ([]repository.User, error) {
	return s.r.GetAll(ctx, pag, sort)
}
