package mock

import "user-service/internal/model"

// Создаем мок-реализацию
type MockUserService struct {
	InsertFunc                  func(model *model.UserModel) (*model.UserModel, error)
	GetAllFunc                  func(page, limit int) ([]*model.UserModel, error)
	GetOneFunc                  func(id int) (*model.UserModel, error)
	UpdateFunc                  func(id int, model *model.UserModel) (*model.UserModel, error)
	DeleteFunc                  func(id int) error
	ExistByIdFunc               func(id int) (bool, error)
	ExistByUsernameFunc         func(username string) (bool, error)
	ExistByUsernameAndNotIdFunc func(username string, id int) (bool, error)
}

// Insert implements service.UserService.
func (m *MockUserService) Create(userModel *model.UserModel) (*model.UserModel, error) {
	return m.InsertFunc(userModel)
}

// GetAll implements service.UserService.
func (m *MockUserService) GetAll(page, limit int) ([]*model.UserModel, error) {
	return m.GetAllFunc(page, limit)
}

// GetOne implements service.UserService.
func (m *MockUserService) GetOne(id int) (*model.UserModel, error) {
	return m.GetOneFunc(id)
}

// Update implements service.UserService.
func (m *MockUserService) Update(id int, userModel *model.UserModel) (*model.UserModel, error) {
	return m.UpdateFunc(id, userModel)
}

// Delete implements service.UserService.
func (m *MockUserService) Delete(id int) error {
	return m.DeleteFunc(id)
}

// ExistById implements service.UserService.
func (m *MockUserService) ExistById(id int) (bool, error) {
	return m.ExistByIdFunc(id)
}

// ExistByUsername implements service.UserService.
func (m *MockUserService) ExistByUsername(username string) (bool, error) {
	return m.ExistByUsernameFunc(username)
}

// ExistByUsernameAndNotId implements service.UserService.
func (m *MockUserService) ExistByUsernameAndNotId(username string, id int) (bool, error) {
	return m.ExistByUsernameAndNotIdFunc(username, id)
}
