package postgres

import (
	"avanpost-show/internal/entity"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) GetUsersPaged(
	ctx context.Context,
	filter *entity.UserFilterQuery,
) (userItems []*entity.UserShort, count int, err error) {
	filters, args := generateUserFilterCondition(filter)

	countUsersQuery := strings.Join([]string{countUsersQueryBase, filters}, " ")
	userPagedResourcesQuery := strings.Join([]string{getPagedUsersQueryBase, filters, getPagedUsersQueryTail}, " ")

	if err := r.db.QueryRow(ctx, countUsersQuery, args).Scan(&count); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, userPagedResourcesQuery, args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[userShort])
	if err != nil {
		return nil, 0, err
	}
	return toUserShortList(items), count, nil
}

func generateUserFilterCondition(params *entity.UserFilterQuery) (string, pgx.NamedArgs) {
	args := pgx.NamedArgs{
		"page_size": params.PageSize,
		"offset":    params.GetOffset(),
	}
	filters := make([]string, 0, 3)

	if params.Search != "" {
		filters = append(filters, searchFilter)
		args["search"] = strings.ToLower(params.Search)
	}
	sb := strings.Builder{}
	if len(filters) > 0 {
		sb.WriteString(" where ")
		sb.WriteString(strings.Join(filters, " and "))
	}
	return sb.String(), args
}

func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	rows, err := r.db.Query(ctx, getUserByIDQuery, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	item, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[user])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toUser(item), nil
}

func (r *repository) CreateUser(ctx context.Context, id uuid.UUID, createTime time.Time, newUser *entity.UserEdit) (err error) {
	_, err = r.db.Exec(ctx, createUserCommand, pgx.NamedArgs{
		"id":         id,
		"username":   newUser.Username,
		"firstname":  newUser.Firstname,
		"lastname":   newUser.Lastname,
		"middlename": newUser.Middlename,
		"email":      newUser.Email,
		"created_at": createTime,
		"updated_at": createTime,
	})
	return
}

func (r *repository) UpdateUser(ctx context.Context, id uuid.UUID, editTime time.Time, edit *entity.UserEdit) (err error) {
	_, err = r.db.Exec(ctx, updateUserCommand, pgx.NamedArgs{
		"id":         id,
		"username":   edit.Username,
		"firstname":  edit.Firstname,
		"lastname":   edit.Lastname,
		"middlename": edit.Middlename,
		"email":      edit.Email,
		"updated_at": editTime,
	})
	return
}

func (r *repository) DeleteUser(ctx context.Context, id uuid.UUID) (err error) {
	_, err = r.db.Exec(ctx, deleteUserCommand, pgx.NamedArgs{"id": id})
	return
}
