package http_controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"user-service/internal/dto"
	"user-service/internal/mapper"
	"user-service/internal/user"
	"user-service/pkg/http_helper"

	"github.com/gorilla/mux"
)

const defaultPageLimit int = 10

type userControllerImpl struct {
	userService user.UserService
	userMapper  mapper.UserMapper
}

// Create implements user.UserController.
func (userController *userControllerImpl) Create() http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, requerst *http.Request) {
			log.Println("UserController.Create:", requerst.URL.Path, "from", requerst.Host)

			data, err := io.ReadAll(requerst.Body)
			if err != nil {
				errText := "Ошибка чтения тела запроса!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.Create:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			if len(data) == 0 {
				errText := "Получено пустое тело запроса!"
				http.Error(responseWriter, errText, http.StatusBadRequest)
				log.Println(
					"UserController.Create:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			var dtoRequest dto.UserRequest
			if json.Unmarshal(data, &dtoRequest) != nil {
				errText := "Ошибка конвертации JSON в dto.UserRequest!"
				http.Error(responseWriter, errText, http.StatusBadRequest)
				log.Println(
					"UserController.Create:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			userModel, err := userController.userMapper.ToModel(dtoRequest)
			if err != nil {
				errText := "Ошибка конвертации dto.UserRequest в model.UserModel!"
				http.Error(responseWriter, errText, http.StatusBadRequest)
				log.Println(
					"UserController.Create:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			if isExistUsername, err := userController.userService.ExistByUsername(userModel.Username); err != nil {
				errText := "Ошибка проверки существования пользователя с заданным именем!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.Create:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			} else if isExistUsername {
				errText := "Пользователь с заданным именем уже существует!"
				http.Error(responseWriter, errText, http.StatusConflict)
				log.Println(
					"UserController.Create:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			userModel, err = userController.userService.Create(userModel)
			if err != nil {
				errText := "Ошибка создания записи в базе данных!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.Create:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			dtoResponse := userController.userMapper.ToDto(*userModel)
			responseWriter.WriteHeader(http.StatusCreated)

			if err := json.NewEncoder(responseWriter).Encode(dtoResponse); err != nil {
				errText := "Ошибка конвертации dto.UserResponse в JSON!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.Create:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			log.Println(
				"UserController.Create:", requerst.URL.Path,
				"from", requerst.Host,
				"result:",
				dtoResponse,
			)
		},
	)
}

// GetAll implements user.UserController.
func (userController *userControllerImpl) GetAll() http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, requerst *http.Request) {
			log.Println("UserController.GetAll Handler Serving:", requerst.URL.Path, "from", requerst.Host)

			queryParams := requerst.URL.Query()

			pageNubmer, err := http_helper.GetQueryParam[int](queryParams, "page")
			if err != nil {
				switch err {
				case http_helper.ErrParamIsEmpty:
					pageNubmer = 0
				default:
					errText := "Ошибка конвертации параметра запроса id!"
					http.Error(responseWriter, errText, http.StatusInternalServerError)
					log.Println(
						"UserController.GetAll:", requerst.URL.Path,
						"from", requerst.Host,
						errText,
						err,
					)
					return
				}
			}
			if pageNubmer > 0 {
				pageNubmer--
			}

			pageLimit, err := http_helper.GetQueryParam[int](queryParams, "limit")
			if err != nil {
				switch err {
				case http_helper.ErrParamIsEmpty:
					pageLimit = defaultPageLimit
				default:
					errText := "Ошибка конвертации параметра запроса limit!"
					http.Error(responseWriter, errText, http.StatusInternalServerError)
					log.Println(
						"UserController.GetAll:", requerst.URL.Path,
						"from", requerst.Host,
						errText,
						err,
					)
					return
				}
			}

			models, err := userController.userService.GetAll(pageNubmer, pageLimit)
			if err != nil {
				errText := "Ошибка получения списка пользователей!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.GetAll:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			var dtoResponses []dto.UserResponse = make([]dto.UserResponse, 0)
			for _, value := range models {
				dtoResponses = append(dtoResponses, *userController.userMapper.ToDto(*value))
			}

			responseWriter.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(responseWriter).Encode(dtoResponses); err != nil {
				errText := "Ошибка конвертации []response.UserResponse в JSON!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.GetAll:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			log.Println(
				"UserController.GetAll:", requerst.URL.Path,
				"from", requerst.Host,
				"result:",
				dtoResponses,
			)
		},
	)
}

// GetOne implements user.UserController.
func (userController *userControllerImpl) GetOne() http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, requerst *http.Request) {
			log.Println("UserController.GetOne Handler Serving:", requerst.URL.Path, "from", requerst.Host)

			params := mux.Vars(requerst)
			id, err := http_helper.GetRouteParam[int](params, "id")
			if err != nil {
				switch err {
				case http_helper.ErrParamIsEmpty:
					errText := "Ошибка заполнения параметра запроса id!"
					http.Error(responseWriter, errText, http.StatusBadRequest)
					log.Println(
						"UserController.GetOne:", requerst.URL.Path,
						"from", requerst.Host,
						errText,
						err,
					)
					return
				default:
					errText := "Ошибка конвертации параметра запроса id!"
					http.Error(responseWriter, errText, http.StatusInternalServerError)
					log.Println(
						"UserController.GetOne:", requerst.URL.Path,
						"from", requerst.Host,
						errText,
						err,
					)
					return
				}
			}

			if isExist, err := userController.userService.ExistById(id); err != nil {
				errText := "Ошибка проверки существования пользователя с заданным id!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.GetOne:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			} else if !isExist {
				errText := "Пользователь с заданным id не существует!"
				http.Error(responseWriter, errText, http.StatusNotFound)
				log.Println(
					"UserController.GetOne:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			userModel, err := userController.userService.GetOne(id)
			if err != nil {
				errText := "Ошибка получения пользователя!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.GetOne:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			dtoResponse := userController.userMapper.ToDto(*userModel)
			responseWriter.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(responseWriter).Encode(dtoResponse); err != nil {
				errText := "Ошибка конвертации response.UserResponse в JSON!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.GetOne:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			log.Println(
				"UserController.GetOne:", requerst.URL.Path,
				"from", requerst.Host,
				"result:",
				dtoResponse,
			)
		},
	)
}

// Update implements user.UserController.
func (userController *userControllerImpl) Update() http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, requerst *http.Request) {
			log.Println("UserController.Update Handler Serving:", requerst.URL.Path, "from", requerst.Host)

			params := mux.Vars(requerst)
			id, err := http_helper.GetRouteParam[int](params, "id")
			if err != nil {
				switch err {
				case http_helper.ErrParamIsEmpty:
					errText := "Ошибка заполнения параметра запроса id!"
					http.Error(responseWriter, errText, http.StatusBadRequest)
					log.Println(
						"UserController.Update:", requerst.URL.Path,
						"from", requerst.Host,
						errText,
						err,
					)
					return
				default:
					errText := "Ошибка конвертации параметра запроса id!"
					http.Error(responseWriter, errText, http.StatusInternalServerError)
					log.Println(
						"UserController.Update:", requerst.URL.Path,
						"from", requerst.Host,
						errText,
						err,
					)
					return
				}
			}

			if isExist, err := userController.userService.ExistById(id); err != nil {
				errText := "Ошибка проверки существования пользователя с заданным id!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.Update:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			} else if !isExist {
				errText := "Пользователь с заданным id не существует!"
				http.Error(responseWriter, errText, http.StatusNotFound)
				log.Println(
					"UserController.Update:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			data, err := io.ReadAll(requerst.Body)
			if err != nil {
				errText := "Ошибка чтения тела запроса!"
				http.Error(responseWriter, errText, http.StatusBadRequest)
				log.Println(
					"UserController.Update:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			if len(data) == 0 {
				errText := "Получено пустое тело запроса!"
				http.Error(responseWriter, errText, http.StatusBadRequest)
				log.Println(
					"UserController.Update:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			var dtoRequest dto.UserRequest
			if json.Unmarshal(data, &dtoRequest) != nil {
				errText := "Ошибка конвертации JSON в dto.UserRequest!"
				http.Error(responseWriter, errText, http.StatusBadRequest)
				log.Println(
					"UserController.Update:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			userModel, err := userController.userMapper.ToModel(dtoRequest)
			if err != nil {
				errText := "Ошибка конвертации dto.UserRequest в model.UserModel!"
				http.Error(responseWriter, errText, http.StatusBadRequest)
				log.Println(
					"UserController.Update:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			if isExistUsername, err := userController.userService.ExistByUsernameAndNotId(userModel.Username, id); err != nil {
				errText := "Ошибка проверки существования пользователя с заданным именем!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.Update:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			} else if isExistUsername {
				errText := "Пользователь с заданным именем уже существует!"
				http.Error(responseWriter, errText, http.StatusConflict)
				log.Println(
					"UserController.Update:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			userModel, err = userController.userService.Update(id, userModel)
			if err != nil {
				errText := "Ошибка обновления записи в базе данных!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.Update:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			dtoResponse := userController.userMapper.ToDto(*userModel)
			responseWriter.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(responseWriter).Encode(dtoResponse); err != nil {
				errText := "Ошибка конвертации response.UserResponse в JSON!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.Update:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			log.Println(
				"UserController.Update:", requerst.URL.Path,
				"from", requerst.Host,
				"result:",
				dtoResponse,
			)
		},
	)
}

// Delete implements user.UserController.
func (userController *userControllerImpl) Delete() http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, requerst *http.Request) {
			log.Println("UserController.Delete Handler Serving:", requerst.URL.Path, "from", requerst.Host)

			params := mux.Vars(requerst)
			id, err := http_helper.GetRouteParam[int](params, "id")
			if err != nil {
				switch err {
				case http_helper.ErrParamIsEmpty:
					errText := "Ошибка заполнения параметра запроса id!"
					http.Error(responseWriter, errText, http.StatusBadRequest)
					log.Println(
						"UserController.Delete:", requerst.URL.Path,
						"from", requerst.Host,
						errText,
						err,
					)
					return
				default:
					errText := "Ошибка конвертации параметра запроса id!"
					http.Error(responseWriter, errText, http.StatusInternalServerError)
					log.Println(
						"UserController.Delete:", requerst.URL.Path,
						"from", requerst.Host,
						errText,
						err,
					)
					return
				}
			}

			if isExist, err := userController.userService.ExistById(id); err != nil {
				errText := "Ошибка проверки существования пользователя с заданным id!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.Delete:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			} else if !isExist {
				errText := "Пользователь с заданным id не существует!"
				http.Error(responseWriter, errText, http.StatusNotFound)
				log.Println(
					"UserController.Delete:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			if err := userController.userService.Delete(id); err != nil {
				errText := "Ошибка удаления пользователя!"
				http.Error(responseWriter, errText, http.StatusInternalServerError)
				log.Println(
					"UserController.Delete:", requerst.URL.Path,
					"from", requerst.Host,
					errText,
					err,
				)
				return
			}

			responseWriter.WriteHeader(http.StatusNoContent)
			log.Printf("UserController.Delete(%d) is success\n", id)
		},
	)
}

func NewUserController(userService user.UserService) user.UserController {
	return &userControllerImpl{
		userService: userService,
	}
}
