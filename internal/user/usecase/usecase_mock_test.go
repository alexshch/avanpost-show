package usecase_test

import (
	"context"
	"errors"
	"testing"

	"avanpost-show/internal/entity"
	"avanpost-show/internal/user/usecase"
	"avanpost-show/internal/user/usecase/mock"
	"avanpost-show/pkg/apierror"

	"github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

func TestUseCase_GetUsersPaged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockRepository(ctrl)
	filter := &entity.UserFilterQuery{
		PageParam: entity.PageParam{PageIndex: 2, PageSize: 5},
		Search:    "john",
	}
	users := []*entity.UserShort{{ID: "1", Username: "john", FullName: "John Doe", IsActive: true}}
	repoMock.EXPECT().GetUsersPaged(gomock.Any(), filter).Return(users, 1, nil)

	p := mock.NewMockPublisher()
	u := usecase.NewUseCase(repoMock, p)
	got, total, err := u.GetUsersPaged(context.Background(), filter)
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if len(got) != 1 || got[0].ID != "1" {
		t.Fatalf("unexpected users returned: %#v", got)
	}
}

func TestUseCase_CreateUser_RetrievesCreatedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockRepository(ctrl)
	input := &entity.UserEdit{Username: "jdoe", Firstname: "John", Lastname: "Doe", Email: "john@example.com"}
	createdUser := &entity.User{ID: "abc", Username: "jdoe", Firstname: "John", Lastname: "Doe", Email: "john@example.com"}

	repoMock.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any(), input).Return(nil)
	repoMock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(createdUser, nil)

	p := mock.NewMockPublisher()
	u := usecase.NewUseCase(repoMock, p)
	got, err := u.CreateUser(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "abc" || got.Username != "jdoe" {
		t.Fatalf("got unexpected user: %#v", got)
	}
}

func TestUseCase_UpdateUser_ReturnsEntityNotFound_WhenUserMissing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockRepository(ctrl)
	id := uuid.New()
	repoMock.EXPECT().GetUserByID(gomock.Any(), id).Return(nil, nil)

	p := mock.NewMockPublisher()
	u := usecase.NewUseCase(repoMock, p)
	err := u.UpdateUser(context.Background(), id, &entity.UserEdit{Username: "jdoe"})
	if !errors.Is(err, apierror.EntityNotFound) {
		t.Fatalf("expected EntityNotFound, got %v", err)
	}
}
