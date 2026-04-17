package usecase

import (
	"avanpost-show/internal/entity"
	"avanpost-show/pkg/apierror"
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	GetUsersPaged(
		ctx context.Context,
		filter *entity.UserFilterQuery,
	) ([]*entity.UserShort, int, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, userID uuid.UUID, now time.Time, newUser *entity.UserEdit) error
	UpdateUser(ctx context.Context, id uuid.UUID, now time.Time, edit *entity.UserEdit) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
type Publisher interface {
	Publish(subj string, t any) error
}

type UseCase struct {
	r Repository
	p Publisher
}

func NewUseCase(r Repository, p Publisher) *UseCase {
	return &UseCase{r: r, p: p}
}

func (u *UseCase) GetUsersPaged(
	ctx context.Context,
	filter *entity.UserFilterQuery,
) ([]*entity.UserShort, int, error) {
	users, count, err := u.r.GetUsersPaged(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (u *UseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return u.r.GetUserByID(ctx, id)
}

func (u *UseCase) CreateUser(ctx context.Context, newUser *entity.UserEdit) (*entity.User, error) {
	id := uuid.New()
	now := time.Now()
	err := u.r.CreateUser(ctx, id, now, newUser)
	if err != nil {
		return nil, err
	}
	user, err := u.r.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	event := entity.Event{ID: user.ID}
	if err = u.p.Publish("user.created", event); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UseCase) UpdateUser(ctx context.Context, id uuid.UUID, edit *entity.UserEdit) error {
	user, err := u.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return apierror.EntityNotFound
	}
	return u.r.UpdateUser(ctx, id, time.Now(), edit)
}

func (u *UseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.r.DeleteUser(ctx, id)
}
