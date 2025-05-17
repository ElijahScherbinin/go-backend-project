package user

import "user-service/internal/model"

type UserService interface {
	Create(userModel *model.UserModel) (*model.UserModel, error)
	GetAll(page, limit int) ([]*model.UserModel, error)
	GetOne(id int) (*model.UserModel, error)
	Update(id int, userModel *model.UserModel) (*model.UserModel, error)
	Delete(id int) error
	ExistById(id int) (bool, error)
	ExistByUsername(username string) (bool, error)
	ExistByUsernameAndNotId(username string, id int) (bool, error)
}
