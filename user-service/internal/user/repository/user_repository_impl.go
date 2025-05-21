package repository

import (
	"fmt"
	"log"
	"user-service/internal/model"
	"user-service/internal/user"
	"user-service/pkg/db"
)

type userRepositoryImpl struct {
	db db.DB
}

// Create implements user.UserRepository.
func (userRepository *userRepositoryImpl) Create(userModel *model.UserModel) (*model.UserModel, error) {
	dbConnect, err := userRepository.db.OpenConnect()
	if err != nil {
		return nil, err
	}
	defer dbConnect.Close()

	var id int
	dbConnect.QueryRow(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) returning id",
		&userModel.Username,
		&userModel.Password_Hash,
	).Scan(&id)

	return userRepository.getOneById(id)
}

// GetAll implements user.UserRepository.
func (userRepository *userRepositoryImpl) GetAll(offset, limit int) ([]*model.UserModel, error) {
	dbConnect, err := userRepository.db.OpenConnect()
	if err != nil {
		return nil, err
	}
	defer dbConnect.Close()

	rows, err := dbConnect.Query("SELECT id, username, password_hash FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}

	var users []*model.UserModel
	for rows.Next() {
		var foundUser model.UserModel
		if err := rows.Scan(&foundUser.Id, &foundUser.Username, &foundUser.Password_Hash); err != nil {
			log.Println("Ошибка сканирования строки из запроса GetAll!")
			continue
		}
		users = append(users, &foundUser)
	}

	return users, nil
}

// GetOne implements user.UserRepository.
func (userRepository *userRepositoryImpl) GetOne(id int) (*model.UserModel, error) {
	return userRepository.getOneById(id)
}

// Update implements user.UserRepository.
func (userRepository *userRepositoryImpl) Update(id int, userModel *model.UserModel) (*model.UserModel, error) {
	dbConnect, err := userRepository.db.OpenConnect()
	if err != nil {
		return nil, err
	}
	defer dbConnect.Close()

	result, err := dbConnect.Exec("UPDATE users SET username=$1, password_hash=$2 WHERE id=$3",
		&userModel.Username,
		&userModel.Password_Hash,
		&id,
	)
	if err != nil {
		return nil, err
	}

	if count, err := result.RowsAffected(); count == 0 && err == nil {
		return nil, fmt.Errorf("не удалось обновить пользователя с 'id'=%d", id)
	}

	return userRepository.getOneById(id)
}

// Delete implements user.UserRepository.
func (userRepository *userRepositoryImpl) Delete(id int) error {
	dbConnect, err := userRepository.db.OpenConnect()
	if err != nil {
		return err
	}
	defer dbConnect.Close()

	result, err := dbConnect.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	if count, err := result.RowsAffected(); count == 0 && err == nil {
		return fmt.Errorf("не удалось удалить пользователя с 'id'=%d", id)
	}

	return nil
}

// ExistById implements user.UserRepository.
func (userRepository *userRepositoryImpl) ExistById(id int) (bool, error) {
	dbConnect, err := userRepository.db.OpenConnect()
	if err != nil {
		return false, err
	}
	defer dbConnect.Close()

	var count int
	if err = dbConnect.QueryRow("SELECT COUNT(*) FROM users WHERE id=$1", id).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistByUsername implements user.UserRepository.
func (userRepository *userRepositoryImpl) ExistByUsername(username string) (bool, error) {
	dbConnect, err := userRepository.db.OpenConnect()
	if err != nil {
		return false, err
	}
	defer dbConnect.Close()

	var count int
	if err = dbConnect.QueryRow("SELECT COUNT(*) FROM users WHERE username=$1", username).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistByUsernameAndNotId implements user.UserRepository.
func (userRepository *userRepositoryImpl) ExistByUsernameAndNotId(username string, id int) (bool, error) {
	dbConnect, err := userRepository.db.OpenConnect()
	if err != nil {
		return false, err
	}
	defer dbConnect.Close()

	var count int
	if err = dbConnect.QueryRow("SELECT COUNT(*) FROM users WHERE username=$1 and id!=$2", username, id).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (userRepository *userRepositoryImpl) getOneById(id int) (*model.UserModel, error) {
	dbConnect, err := userRepository.db.OpenConnect()
	if err != nil {
		return nil, err
	}
	defer dbConnect.Close()

	row := dbConnect.QueryRow("SELECT id, username, password_hash FROM users WHERE id=$1", id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var foundUser model.UserModel
	err = row.Scan(&foundUser.Id, &foundUser.Username, &foundUser.Password_Hash)
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}

func NewUserRepository(db db.DB) user.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}
