package mocking

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/mocking/mocks"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// We will use this mocked interface which contains all the functions that need to be mocked.
	mockDB := mocks.NewMockDbHandlerIf(ctrl)
	name := "some_name"
	password := "some_password"
	mockDB.EXPECT().CreateDatabase(name).Return(sqlmock.NewResult(1, 1), nil)
	mockDB.EXPECT().CreateDatabaseUser(name).Return(sqlmock.NewResult(1, 1), nil)
	mockDB.EXPECT().AssignPassword(name, password).Return(sqlmock.NewResult(1, 1), nil)
	mockDB.EXPECT().GrantPrivileges(name, name).Return(nil, errors.New("failed"))

	var userCreated bool
	userCreated = true

	// Here we are passing the mockDB handler, as the handleDatabaseCreation() accepts the db object
	handleDatabaseCreation(mockDB, name, password, &userCreated)
}
