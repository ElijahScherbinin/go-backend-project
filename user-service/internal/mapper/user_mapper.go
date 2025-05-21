package mapper

import (
	"crypto/sha256"
	"fmt"
	"user-service/internal/dto"
	"user-service/internal/model"

	"github.com/go-playground/validator"
)

type UserMapper struct{}

func (userMapper UserMapper) ToModel(userRequest dto.UserRequest) (*model.UserModel, error) {
	if err := validator.New().Struct(userRequest); err != nil {
		return nil, err
	}
	password_hash := sha256.Sum256([]byte(userRequest.Password))
	return &model.UserModel{
		Username:      userRequest.Username,
		Password_Hash: fmt.Sprintf("%x", password_hash),
	}, nil
}

func (userMapper UserMapper) ToDto(userModel model.UserModel) *dto.UserResponse {
	return &dto.UserResponse{
		Id:       userModel.Id,
		Username: userModel.Username,
	}
}
