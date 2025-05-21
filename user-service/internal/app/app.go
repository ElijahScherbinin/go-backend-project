package app

import (
	"log"
	"user-service/internal/config"
	"user-service/internal/server"
	"user-service/internal/user"
	"user-service/internal/user/delivery/http_controller"
	"user-service/internal/user/repository"
	"user-service/internal/user/service"
	"user-service/pkg/db"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const configPath string = "config/config"

func Run() {
	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalln("Ошибка загрузки файла конфигурации.", err)
	}

	database, err := db.New(appConfig.GetDatabaseConfig())
	if err != nil {
		log.Fatalln("Ошибка инициализации экземпляра драйвера базы данных.", err)
	}

	var userRepository user.UserRepository = repository.NewUserRepository(database)
	var userService user.UserService = service.NewUserService(userRepository)
	var userController user.UserController = http_controller.NewUserController(userService)

	routerInstance := mux.NewRouter()
	user.SetupUserRoutes(routerInstance, userController)

	routerInstance.Handle("/metrics", promhttp.Handler())

	var serverConfig config.ServerConfig = appConfig.GetServerConfig()
	if err := server.Run(routerInstance, serverConfig.GetHost(), serverConfig.GetPort()); err != nil {
		log.Fatalln("Ошибка запуска веб-сервера.", err)
	}
}
