package http

import (
	"avanpost-show/internal/entity"
	"avanpost-show/pkg/apierror"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserUseCase interface {
	GetUsersPaged(
		ctx context.Context,
		params *entity.UserFilterQuery,
	) ([]*entity.UserShort, int, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, input *entity.UserEdit) (*entity.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, input *entity.UserEdit) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type handler struct {
	u UserUseCase
}

func NewUserHandler(u UserUseCase) *handler {
	return &handler{u: u}
}

// GetUsers
// @Tags User
// @Summary Get user paged list
// @Description Get user paged list
// @Produce json
// @Param pageIndex query string false "page number"
// @Param pageSize query string false "number of elements"
// @Param search query string false "searching user by username, name"
// @Success 200 {object} entity.UsersPaged
// @Failure 500 {object} entity.Message "Internal Server Error"
// @Router /users [get]
func (h *handler) GetUsers(c echo.Context) error {
	qp := c.QueryParam("pageIndex")
	pageNumber, err := strconv.Atoi(qp)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}
	qs := c.QueryParam("pageSize")
	pageSize, err := strconv.Atoi(qs)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	pageInfo := &entity.UserFilterQuery{
		PageParam: entity.PageParam{
			PageIndex: pageNumber,
			PageSize:  pageSize,
		},
		Search: c.QueryParam("search"),
	}
	users, total, err := h.u.GetUsersPaged(c.Request().Context(), pageInfo)
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(
			http.StatusInternalServerError,
			entity.NewMessage(http.StatusText(http.StatusInternalServerError)),
		)
	}
	return c.JSON(http.StatusOK, entity.NewPagedItemsList(pageNumber, pageSize, total, users))
}

// GetUserByID
// @Tags User
// @Summary Get user by id
// @Description Get user by user id
// @Produce json
// @Param id path string true "user ID"
// @Success 200 {object} entity.User
// @Failure	404	{object} entity.Message "Not Found"
// @Failure 500 {object} entity.Message "Internal Server Error"
// @Router /users/{id} [get]
func (h *handler) GetUserByID(c echo.Context) error {
	idString := c.Param("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			entity.NewMessage(http.StatusText(http.StatusNotFound)),
		)
	}
	user, err := h.u.GetUserByID(c.Request().Context(), id)
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(
			http.StatusInternalServerError,
			entity.NewMessage(http.StatusText(http.StatusInternalServerError)),
		)
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, entity.NewMessage(http.StatusText(http.StatusNotFound)))
	}
	return c.JSON(http.StatusOK, user)
}

// CreateUser
// @Tags User
// @Summary Create user
// @Description Create user
// @Accept json
// @Produce json
// @Param user body entity.UserEdit true "user"
// @Success 201 {object} entity.User
// @Failure	404	{object} entity.Message "Not Found"
// @Failure 500 {object} entity.Message "Internal Server Error"
// @Router /users [post]
func (h *handler) CreateUser(c echo.Context) error {
	newUser := &entity.UserEdit{}
	if err := c.Bind(newUser); err != nil {
		slog.Error(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, entity.NewMessage(http.StatusText(http.StatusUnprocessableEntity)))
	}
	role, err := h.u.CreateUser(c.Request().Context(), newUser)
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(
			http.StatusInternalServerError,
			entity.NewMessage(http.StatusText(http.StatusInternalServerError)),
		)
	}
	return c.JSON(http.StatusCreated, role)
}

// UpdateUser
// @Tags User
// @Summary Update user
// @Description Update user
// @Accept json
// @Produce json
// @Param id path string true "user ID"
// @Param rule body entity.UserEdit true "user body"
// @Success 204
// @Failure	404	{object} entity.Message "Not Found"
// @Failure 500 {object} entity.Message "Internal Server Error"
// @Router /users/{id} [put]
func (h *handler) UpdateUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			entity.NewMessage(http.StatusText(http.StatusNotFound)),
		)
	}
	updateUser := &entity.UserEdit{}
	if err := c.Bind(updateUser); err != nil {
		slog.Error(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, entity.NewMessage(http.StatusText(http.StatusUnprocessableEntity)))
	}
	err = h.u.UpdateUser(c.Request().Context(), id, updateUser)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, apierror.EntityNotFound) {
			return c.JSON(http.StatusNotFound, entity.NewMessage(http.StatusText(http.StatusNotFound)))
		}
		return c.JSON(
			http.StatusInternalServerError,
			entity.NewMessage(http.StatusText(http.StatusInternalServerError)),
		)
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteUser
// @Tags User
// @Summary Delete user
// @Description Delete user
// @Produce json
// @Param id path string true "user ID"
// @Success 204
// @Failure 404 {object} entity.Message "Not Found"
// @Failure 500 {object} entity.Message "Internal Server Error"
// @Router /users/{id} [delete]
func (h *handler) DeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			entity.NewMessage(http.StatusText(http.StatusNotFound)),
		)
	}
	err = h.u.DeleteUser(c.Request().Context(), id)
	if err != nil {
		slog.Error(err.Error())
		if errors.Is(err, apierror.EntityNotFound) {
			return c.JSON(http.StatusNotFound, entity.NewMessage(http.StatusText(http.StatusNotFound)))
		}
		return c.JSON(
			http.StatusInternalServerError,
			entity.NewMessage(http.StatusText(http.StatusInternalServerError)),
		)
	}
	return c.NoContent(http.StatusNoContent)
}
