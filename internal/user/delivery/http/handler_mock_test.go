package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"avanpost-show/internal/entity"
	httpdelivery "avanpost-show/internal/user/delivery/http"
	"avanpost-show/internal/user/delivery/http/mock"
	"avanpost-show/pkg/apierror"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	gomock "go.uber.org/mock/gomock"
)

func TestGetUsers_ReturnsOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUserUseCase(ctrl)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users?pageIndex=2&pageSize=5&search=test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	filter := &entity.UserFilterQuery{
		PageParam: entity.PageParam{PageIndex: 2, PageSize: 5},
		Search:    "test",
	}
	users := []*entity.UserShort{{ID: "1", Username: "testuser", FullName: "Test User", IsActive: true}}
	mockUseCase.EXPECT().GetUsersPaged(gomock.Any(), filter).Return(users, 1, nil)

	h := httpdelivery.NewUserHandler(mockUseCase)
	if err := h.GetUsers(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var got struct {
		PageIndex int                 `json:"pageIndex"`
		PageSize  int                 `json:"pageSize"`
		Total     int                 `json:"total"`
		Items     []*entity.UserShort `json:"items"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Total != 1 || got.PageIndex != 2 || got.PageSize != 5 || len(got.Items) != 1 {
		t.Fatalf("unexpected response payload: %#v", got)
	}
}

func TestGetUserByID_ReturnsNotFound_WhenMissing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUserUseCase(ctrl)
	e := echo.New()
	id := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/users/"+id, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id)

	mockUseCase.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, nil)

	h := httpdelivery.NewUserHandler(mockUseCase)
	if err := h.GetUserByID(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestUpdateUser_ReturnsNotFound_OnEntityNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock.NewMockUserUseCase(ctrl)
	e := echo.New()
	id := uuid.New().String()
	payload, err := json.Marshal(&entity.UserEdit{Username: "jdoe"})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPut, "/users/"+id, bytes.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id)

	mockUseCase.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(apierror.EntityNotFound)

	h := httpdelivery.NewUserHandler(mockUseCase)
	if err := h.UpdateUser(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
