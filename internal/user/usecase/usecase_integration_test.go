package usecase_test

import (
	"encoding/json"
	"testing"
	"time"

	"avanpost-show/internal/entity"
	"avanpost-show/internal/user/repository/postgres"
	"avanpost-show/internal/user/usecase"
	"avanpost-show/pkg/publisher"
	"avanpost-show/pkg/test_suite"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UseCaseIntegrationTestSuite struct {
	test_suite.DBTestSuite
	repo      usecase.Repository
	publisher usecase.Publisher
	uc        *usecase.UseCase
	userID    uuid.UUID
}

func TestUseCaseIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &UseCaseIntegrationTestSuite{
		DBTestSuite: *test_suite.NewDBTestSuite("../../../migration/postgres"),
	})
}

func (s *UseCaseIntegrationTestSuite) SetupTest() {
	s.repo = postgres.NewRepository(s.DBPool)
	s.publisher = publisher.NewPublisher(s.Nc)
	s.uc = usecase.NewUseCase(s.repo, s.publisher)
	// Clean up tables before each test
	s.CleanupTables([]string{"users"})
}

func (s *UseCaseIntegrationTestSuite) TearDownTest() {

}

func (s *UseCaseIntegrationTestSuite) TestCreateUserIntegration() {
	eventCh := make(chan entity.Event, 1)

	// Subscribe to user.created events
	_, err := s.Nc.Subscribe("user.created", func(msg *nats.Msg) {
		var event entity.Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			s.T().Errorf("Failed to unmarshal event: %v", err)
			return
		}
		eventCh <- event
	})
	assert.NoError(s.T(), err)

	userEdit := &entity.UserEdit{
		Username:  "integrationuser",
		Email:     "integration@example.com",
		Firstname: "Integration",
		Lastname:  "User",
	}

	createdUser, err := s.uc.CreateUser(s.Ctx, userEdit)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), createdUser)
	assert.Equal(s.T(), userEdit.Username, createdUser.Username)
	assert.Equal(s.T(), userEdit.Email, createdUser.Email)
	assert.Equal(s.T(), userEdit.Firstname, createdUser.Firstname)
	assert.Equal(s.T(), userEdit.Lastname, createdUser.Lastname)
	assert.True(s.T(), createdUser.IsActive)
	assert.NotZero(s.T(), createdUser.CreatedAt)
	assert.NotZero(s.T(), createdUser.UpdatedAt)

	// Wait for the user.created event
	select {
	case event := <-eventCh:
		assert.Equal(s.T(), createdUser.ID, event.ID)
	case <-time.After(1 * time.Second):
		s.T().Error("Expected user.created event not received within timeout")
	}
}

func (s *UseCaseIntegrationTestSuite) TestGetUsersPagedIntegration() {
	// Create multiple users
	users := []*entity.UserEdit{
		{Username: "intuser1", Email: "intuser1@example.com", Firstname: "Int", Lastname: "User1"},
		{Username: "intuser2", Email: "intuser2@example.com", Firstname: "Int", Lastname: "User2"},
		{Username: "intuser3", Email: "intuser3@example.com", Firstname: "Int", Lastname: "User3"},
	}

	for _, u := range users {
		_, err := s.uc.CreateUser(s.Ctx, u)
		assert.NoError(s.T(), err)
	}

	// Test pagination
	filter := &entity.UserFilterQuery{
		PageParam: entity.PageParam{PageIndex: 1, PageSize: 2},
	}

	userList, total, err := s.uc.GetUsersPaged(s.Ctx, filter)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 3, total)
	assert.Len(s.T(), userList, 2)
}

func (s *UseCaseIntegrationTestSuite) TestGetUsersPagedWithSearchIntegration() {
	// Create users
	users := []*entity.UserEdit{
		{Username: "search_john", Email: "john@example.com", Firstname: "John", Lastname: "Search"},
		{Username: "search_jane", Email: "jane@example.com", Firstname: "Jane", Lastname: "Search"},
		{Username: "search_bob", Email: "bob@example.com", Firstname: "Bob", Lastname: "Search"},
	}

	for _, u := range users {
		_, err := s.uc.CreateUser(s.Ctx, u)
		assert.NoError(s.T(), err)
	}

	// Search for "search"
	filter := &entity.UserFilterQuery{
		PageParam: entity.PageParam{PageIndex: 1, PageSize: 10},
		Search:    "search",
	}

	userList, total, err := s.uc.GetUsersPaged(s.Ctx, filter)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 3, total)
	assert.Len(s.T(), userList, 3)
}

func (s *UseCaseIntegrationTestSuite) TestUpdateUserIntegration() {
	// Create a user
	userEdit := &entity.UserEdit{
		Username:  "updateuser",
		Email:     "update@example.com",
		Firstname: "Update",
		Lastname:  "User",
	}

	createdUser, err := s.uc.CreateUser(s.Ctx, userEdit)
	assert.NoError(s.T(), err)

	userID, err := uuid.Parse(createdUser.ID)
	assert.NoError(s.T(), err)

	// Update the user
	updateData := &entity.UserEdit{
		Username:  "updateduser",
		Email:     "updated@example.com",
		Firstname: "Updated",
		Lastname:  "User",
	}

	err = s.uc.UpdateUser(s.Ctx, userID, updateData)
	assert.NoError(s.T(), err)

	// Verify update
	retrievedUser, err := s.uc.GetUserByID(s.Ctx, userID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), retrievedUser)
	assert.Equal(s.T(), updateData.Username, retrievedUser.Username)
	assert.Equal(s.T(), updateData.Email, retrievedUser.Email)
	assert.Equal(s.T(), updateData.Firstname, retrievedUser.Firstname)
	assert.Equal(s.T(), updateData.Lastname, retrievedUser.Lastname)
}

func (s *UseCaseIntegrationTestSuite) TestDeleteUserIntegration() {
	// Create a user
	userEdit := &entity.UserEdit{
		Username:  "deleteuser",
		Email:     "delete@example.com",
		Firstname: "Delete",
		Lastname:  "User",
	}

	createdUser, err := s.uc.CreateUser(s.Ctx, userEdit)
	assert.NoError(s.T(), err)

	userID, err := uuid.Parse(createdUser.ID)
	assert.NoError(s.T(), err)

	// Delete the user
	err = s.uc.DeleteUser(s.Ctx, userID)
	assert.NoError(s.T(), err)

	// Try to get the user - should return nil
	retrievedUser, err := s.uc.GetUserByID(s.Ctx, userID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), retrievedUser)
}

func (s *UseCaseIntegrationTestSuite) TestGetUserByIDIntegration() {
	// Create a user
	userEdit := &entity.UserEdit{
		Username:  "getuser",
		Email:     "get@example.com",
		Firstname: "Get",
		Lastname:  "User",
	}

	createdUser, err := s.uc.CreateUser(s.Ctx, userEdit)
	assert.NoError(s.T(), err)

	userID, err := uuid.Parse(createdUser.ID)
	assert.NoError(s.T(), err)

	// Get the user by ID
	retrievedUser, err := s.uc.GetUserByID(s.Ctx, userID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), retrievedUser)
	assert.Equal(s.T(), createdUser.ID, retrievedUser.ID)
	assert.Equal(s.T(), createdUser.Username, retrievedUser.Username)
}
