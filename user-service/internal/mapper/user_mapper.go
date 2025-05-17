package mapper

import (
	"user-service/internal/dto"
	"user-service/internal/model"

	"github.com/go-playground/validator"
)

type UserMapper struct{}

func (userMapper UserMapper) ToModel(userRequest dto.UserRequest) (*model.UserModel, error) {
	if err := validator.New().Struct(userRequest); err != nil {
		return nil, err
	}
	return &model.UserModel{
		Username: userRequest.Username,
		Password: userRequest.Password,
	}, nil
}

func (userMapper UserMapper) ToDto(userModel model.UserModel) *dto.UserResponse {
	return &dto.UserResponse{
		Id:       userModel.Id,
		Username: userModel.Username,
	}
}
