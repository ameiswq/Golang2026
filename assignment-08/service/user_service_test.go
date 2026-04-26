package service

import (
	"assignment-08/repository"
	"testing"
	"errors"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Sofya", Email: "sofya@mail.com"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestCreateUser(t *testing.T){
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)
	user := &repository.User{ID: 1, Name: "Sofya", Email: "sofya@mail.com"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)
	err := userService.CreateUser(user)
	assert.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	t.Run("user already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 1, Email: "sofya@mail.com"}
		mockRepo.EXPECT().GetByEmail("sofya@mail.com").Return(user, nil)
		err := userService.RegisterUser(user, "sofya@mail.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("new user success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 2, Name: "Sofya", Email: "sofya@mail.com"}
		mockRepo.EXPECT().GetByEmail("sofya@mail.com").Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(nil)
		err := userService.RegisterUser(user, "sofya@mail.com")
		assert.NoError(t, err)
	})

	t.Run("repository error on CreateUser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 2, Name: "Sofya", Email: "sofya@mail.com"}
		mockRepo.EXPECT().GetByEmail("sofya@mail.com").Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(errors.New("create failed"))
		err := userService.RegisterUser(user, "sofya@mail.com")
		assert.Error(t, err)
	})
}

func TestUpdateUserName(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)
		err := userService.UpdateUserName(1, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})

	t.Run("user not found repo error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)
		mockRepo.EXPECT().GetUserByID(1).Return(nil, errors.New("not found"))
		err := userService.UpdateUserName(1, "NewName")
		assert.Error(t, err)
	})

	t.Run("successful update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 1, Name: "OldName", Email: "sofya@mail.com"}
		mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(updated *repository.User) error {
			assert.Equal(t, "NewName", updated.Name)
			return nil
		})
		err := userService.UpdateUserName(1, "NewName")
		assert.NoError(t, err)
		assert.Equal(t, "NewName", user.Name)
	})

	t.Run("UpdateUser fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 1, Name: "OldName", Email: "sofya@mail.com"}
		mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(updated *repository.User) error {
			assert.Equal(t, "NewName", updated.Name)
			return errors.New("update failed")
		})
		err := userService.UpdateUserName(1, "NewName")
		assert.Error(t, err)
		assert.Equal(t, "NewName", user.Name)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("attempt to delete admin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)
		err := userService.DeleteUser(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not allowed to delete admin")
	})

	t.Run("successful delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)
		mockRepo.EXPECT().DeleteUser(2).Return(nil)
		err := userService.DeleteUser(2)
		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := repository.NewMockUserRepository(ctrl)
		userService := NewUserService(mockRepo)
		mockRepo.EXPECT().DeleteUser(2).Return(errors.New("delete failed"))
		err := userService.DeleteUser(2)
		assert.Error(t, err)
	})
}