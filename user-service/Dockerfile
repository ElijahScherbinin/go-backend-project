# Инициализируем образ builder
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

# Создаем директорию для исходных файлов и копируем их туда
RUN mkdir -p $GOPATH/src/user-service
COPY ./cmd/ $GOPATH/src/user-service/cmd/
COPY ./internal/ $GOPATH/src/user-service/internal/
COPY ./pkg/ $GOPATH/src/user-service/pkg/

# Создаем директорию для бинарных файлов 
RUN mkdir -p /bin/user-service
# Копируем содержимое директории config
COPY ./config/ /bin/user-service/config/

# Переходим в директорию с исходными файлами
WORKDIR $GOPATH/src/user-service
# Инициализируем проект
RUN go mod init user-service
# Обновляем зависимости
RUN go mod tidy
# Загружаем зависимости
RUN go mod download
# Запускаем автотесты
RUN go test -v ./...
# Собираем приложение в директорию с бинарными файлами
RUN go build -o /bin/user-service/app ./cmd/app/main.go

# Инициализируем образ runner
FROM alpine:latest AS runner

# Создаем директорию с бинарными файлами
RUN mkdir -p /bin/user-service
# Копируем директорию приложения с бинарными файлами из builder в runner
COPY --from=builder /bin/user-service/ /bin/user-service/

# Открываем порт
EXPOSE 8000

# Переходим в директорию с бинарными файлами
WORKDIR /bin/user-service

# Инициализируем приложение
CMD ["/bin/user-service/app"]