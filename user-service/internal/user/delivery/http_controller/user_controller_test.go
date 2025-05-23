package http_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"

	"user-service/internal/dto"
	"user-service/internal/mapper"
	"user-service/internal/middleware"
	"user-service/internal/model"
	"user-service/internal/user"
	"user-service/internal/user/mock"
	"user-service/pkg/jwt"

	"github.com/go-playground/assert/v2"
	"github.com/gorilla/mux"
)

const jsonRequestTemplate string = "{\"username\": \"%s\", \"password\": \"%s\"}"

var (
	mockExistOneUserModel *model.UserModel = &model.UserModel{
		Id:            1,
		Username:      "user1",
		Password_Hash: "12345678",
	}
	mockExistTwoUserModel *model.UserModel = &model.UserModel{
		Id:            2,
		Username:      "user2",
		Password_Hash: "12345678",
	}
	mockFreeUserModel *model.UserModel = &model.UserModel{
		Id:            0,
		Username:      "free",
		Password_Hash: "12345678",
	}
)

var mockExistModels = []*model.UserModel{
	mockExistOneUserModel,
	mockExistTwoUserModel,
}

var mockService *mock.MockUserService = &mock.MockUserService{
	InsertFunc: func(userModel *model.UserModel) (*model.UserModel, error) {
		return userModel, nil
	},
	GetOneFunc: func(id int) (*model.UserModel, error) {
		for _, userModel := range mockExistModels {
			if userModel.Id == id {
				return userModel, nil
			}
		}
		return nil, nil
	},
	GetAllFunc: func(page, limit int) ([]*model.UserModel, error) {
		var offset int
		if page == 0 {
			offset = 0
		} else {
			offset = limit * page
		}
		var count int
		var result []*model.UserModel
		for index, value := range mockExistModels {
			if index < offset {
				continue
			}
			if count >= limit {
				break
			}
			result = append(result, value)
			count++
		}
		return result, nil
	},
	UpdateFunc: func(id int, userModel *model.UserModel) (*model.UserModel, error) {
		return &model.UserModel{
			Id:            id,
			Username:      userModel.Username,
			Password_Hash: userModel.Password_Hash,
		}, nil
	},
	DeleteFunc: func(id int) error {
		return nil
	},
	ExistByIdFunc: func(id int) (bool, error) {
		for _, userModel := range mockExistModels {
			if userModel.Id == id {
				return true, nil
			}
		}
		return false, nil
	},
	ExistByUsernameFunc: func(username string) (bool, error) {
		for _, userModel := range mockExistModels {
			if userModel.Username == username {
				return true, nil
			}
		}
		return false, nil
	},
	ExistByUsernameAndNotIdFunc: func(username string, id int) (bool, error) {
		for _, userModel := range mockExistModels {
			if userModel.Username == username && userModel.Id != id {
				return true, nil
			}
		}
		return false, nil
	},
}

var userController user.UserController = NewUserController(mockService)
var userRouter *mux.Router = mux.NewRouter()
var userMapper mapper.UserMapper

var token string

func init() {
	userRouter = user.SetupUserRoutes(userRouter, userController)
	jwtCoder := jwt.NewJWTCoder(middleware.Alg, middleware.Secret, middleware.Issuer, middleware.Audience, middleware.ExpirationTimeDuration)
	tokenInstance := jwtCoder.NewToken("username", "admin", []string{"view", "update", "delete"}...)
	strToken, err := jwtCoder.Encode(*tokenInstance)
	if err != nil {
		panic(err)
	}
	token = *strToken
	fmt.Printf("token: '%s'\n", token)
}

func TestCreateHandler(t *testing.T) {
	t.Run("Empty Body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusBadRequest)
	})

	t.Run("JSON to UserRequest", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("{}"))
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusBadRequest)
	})

	t.Run("Check username", func(t *testing.T) {
		jsonBody := fmt.Sprintf(
			jsonRequestTemplate,
			mockExistOneUserModel.Username,
			mockFreeUserModel.Password_Hash,
		)
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(jsonBody))
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusConflict)
	})

	t.Run("User created", func(t *testing.T) {
		jsonBody := fmt.Sprintf(
			jsonRequestTemplate,
			mockFreeUserModel.Username,
			mockFreeUserModel.Password_Hash,
		)
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(jsonBody))
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusCreated)

		var dtoResponse dto.UserResponse
		err := json.NewDecoder(res.Body).Decode(&dtoResponse)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, dtoResponse.Username, mockFreeUserModel.Username)
	})

}

func TestGetAllHandler(t *testing.T) {
	t.Run("Get all: no param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/all", nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusOK)

		var dtoResponses []dto.UserResponse
		err := json.NewDecoder(res.Body).Decode(&dtoResponses)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, len(dtoResponses), 2)

		var dtoResponse dto.UserResponse

		dtoResponse = *userMapper.ToDto(*mockExistOneUserModel)
		assert.Equal(t, slices.Contains(dtoResponses, dtoResponse), true)

		dtoResponse = *userMapper.ToDto(*mockExistTwoUserModel)
		assert.Equal(t, slices.Contains(dtoResponses, dtoResponse), true)
	})

	t.Run("Get all: page 1 limit 1", func(t *testing.T) {
		params := url.Values{}
		params.Set("page", "1")
		params.Set("limit", "1")
		req := httptest.NewRequest(http.MethodGet, "/users/all?"+params.Encode(), nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusOK)

		var dtoResponses []dto.UserResponse
		err := json.NewDecoder(res.Body).Decode(&dtoResponses)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, len(dtoResponses), 1)

		var dtoResponse dto.UserResponse

		dtoResponse = *userMapper.ToDto(*mockExistOneUserModel)
		assert.Equal(t, slices.Contains(dtoResponses, dtoResponse), true)

		dtoResponse = *userMapper.ToDto(*mockExistTwoUserModel)
		assert.Equal(t, slices.Contains(dtoResponses, dtoResponse), false)
	})

	t.Run("Get all: page 2 limit 1", func(t *testing.T) {
		params := url.Values{}
		params.Set("page", "2")
		params.Set("limit", "1")
		req := httptest.NewRequest(http.MethodGet, "/users/all?"+params.Encode(), nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusOK)

		var dtoResponses []dto.UserResponse
		err := json.NewDecoder(res.Body).Decode(&dtoResponses)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, len(dtoResponses), 1)

		var dtoResponse dto.UserResponse

		dtoResponse = *userMapper.ToDto(*mockExistOneUserModel)
		assert.Equal(t, slices.Contains(dtoResponses, dtoResponse), false)

		dtoResponse = *userMapper.ToDto(*mockExistTwoUserModel)
		assert.Equal(t, slices.Contains(dtoResponses, dtoResponse), true)
	})

	t.Run("Get all: page 1 limit 2", func(t *testing.T) {
		params := url.Values{}
		params.Set("page", "1")
		params.Set("limit", "2")
		req := httptest.NewRequest(http.MethodGet, "/users/all?"+params.Encode(), nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusOK)

		var dtoResponses []dto.UserResponse
		err := json.NewDecoder(res.Body).Decode(&dtoResponses)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, len(dtoResponses), 2)

		var dtoResponse dto.UserResponse

		dtoResponse = *userMapper.ToDto(*mockExistOneUserModel)
		assert.Equal(t, slices.Contains(dtoResponses, dtoResponse), true)

		dtoResponse = *userMapper.ToDto(*mockExistTwoUserModel)
		assert.Equal(t, slices.Contains(dtoResponses, dtoResponse), true)
	})

	t.Run("Get all: page 2 limit 2", func(t *testing.T) {
		params := url.Values{}
		params.Set("page", "2")
		params.Set("limit", "2")
		req := httptest.NewRequest(http.MethodGet, "/users/all?"+params.Encode(), nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusOK)

		var dtoResponses []dto.UserResponse
		err := json.NewDecoder(res.Body).Decode(&dtoResponses)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, len(dtoResponses), 0)

		var dtoResponse dto.UserResponse

		dtoResponse = *userMapper.ToDto(*mockExistOneUserModel)
		assert.Equal(t, slices.Contains(dtoResponses, dtoResponse), false)

		dtoResponse = *userMapper.ToDto(*mockExistTwoUserModel)
		assert.Equal(t, slices.Contains(dtoResponses, dtoResponse), false)
	})
}

func TestGetOneHandler(t *testing.T) {
	t.Run("Empty id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/", nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	t.Run("User not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprint("/users/", mockFreeUserModel.Id), nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	t.Run("User found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprint("/users/", mockExistOneUserModel.Id), nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusOK)

		var dtoResponse dto.UserResponse
		err := json.NewDecoder(res.Body).Decode(&dtoResponse)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, dtoResponse.Id, mockExistOneUserModel.Id)
		assert.Equal(t, dtoResponse.Username, mockExistOneUserModel.Username)
	})
}

func TestUpdateHandler(t *testing.T) {
	t.Run("Empty id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/", nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	t.Run("User not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, fmt.Sprint("/users/", mockFreeUserModel.Id), nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	validPath := fmt.Sprint("/users/", mockExistOneUserModel.Id)

	t.Run("Empty Body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, validPath, nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusBadRequest)
	})

	t.Run("JSON to UserRequest", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, validPath, strings.NewReader("{}"))
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusBadRequest)
	})

	t.Run("Check username", func(t *testing.T) {
		jsonBody := fmt.Sprintf(
			jsonRequestTemplate,
			mockExistTwoUserModel.Username,
			mockExistOneUserModel.Password_Hash,
		)
		req := httptest.NewRequest(http.MethodPut, validPath, strings.NewReader(jsonBody))
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusConflict)
	})

	t.Run("User updated", func(t *testing.T) {
		jsonBody := fmt.Sprintf(
			jsonRequestTemplate,
			mockFreeUserModel.Username,
			mockExistOneUserModel.Password_Hash,
		)
		req := httptest.NewRequest(http.MethodPut, validPath, strings.NewReader(jsonBody))
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusOK)

		var dtoResponse dto.UserResponse
		err := json.NewDecoder(res.Body).Decode(&dtoResponse)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, dtoResponse.Id, mockExistOneUserModel.Id)
		assert.Equal(t, dtoResponse.Username, mockFreeUserModel.Username)
	})
}

func TestDeleteHandler(t *testing.T) {
	auth := "Bearer " + token
	t.Run("No token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/0", nil)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusUnauthorized)
	})

	t.Run("Empty id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/", nil)
		req.Header.Add("Authorization", auth)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	t.Run("User not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprint("/users/", mockFreeUserModel.Id), nil)
		req.Header.Add("Authorization", auth)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusNotFound)
	})

	t.Run("User deleted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprint("/users/", mockExistOneUserModel.Id), nil)
		req.Header.Add("Authorization", auth)
		res := httptest.NewRecorder()
		userRouter.ServeHTTP(res, req)
		assert.Equal(t, res.Code, http.StatusNoContent)
	})
}
