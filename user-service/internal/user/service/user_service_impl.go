package service

import (
	"user-service/internal/model"
	"user-service/internal/user"
	"user-service/internal/user/prometheus"
)

type userServiceImpl struct {
	userRepository user.UserRepository
	userPrometheus user.UserPrometheus
}

// Create implements user.UserService.
func (userService *userServiceImpl) Create(userModel *model.UserModel) (*model.UserModel, error) {
	userModel, err := userService.userRepository.Create(userModel)
	if err == nil {
		userService.userPrometheus.New()
	}
	return userModel, err
}

// GetAll implements user.UserService.
func (userService *userServiceImpl) GetAll(page int, limit int) ([]*model.UserModel, error) {
	var offset int
	if page == 0 {
		offset = 0
	} else {
		offset = limit * page
	}
	return userService.userRepository.GetAll(offset, limit)
}

// GetOne implements user.UserService.
func (userService *userServiceImpl) GetOne(id int) (*model.UserModel, error) {
	return userService.userRepository.GetOne(id)
}

// Update implements user.UserService.
func (userService *userServiceImpl) Update(id int, userModel *model.UserModel) (*model.UserModel, error) {
	return userService.userRepository.Update(id, userModel)
}

// Delete implements user.UserService.
func (userService *userServiceImpl) Delete(id int) error {
	err := userService.userRepository.Delete(id)
	if err == nil {
		userService.userPrometheus.Delete()
	}
	return err
}

// ExistById implements user.UserService.
func (userService *userServiceImpl) ExistById(id int) (bool, error) {
	return userService.userRepository.ExistById(id)
}

// ExistByUsername implements user.UserService.
func (userService *userServiceImpl) ExistByUsername(username string) (bool, error) {
	return userService.userRepository.ExistByUsername(username)
}

// ExistByUsernameAndNotId implements user.UserService.
func (userService *userServiceImpl) ExistByUsernameAndNotId(username string, id int) (bool, error) {
	return userService.userRepository.ExistByUsernameAndNotId(username, id)
}

func NewUserService(userRepository user.UserRepository) user.UserService {
	return &userServiceImpl{
		userRepository: userRepository,
		userPrometheus: prometheus.NewUserPrometheus(),
	}
}
