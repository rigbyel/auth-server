# REST API сервис аутентификации

Этот репозиторий содержит REST API сервис для регистрации и аутентификации пользователей

## Конечные точки

1. **Авторизация пользователя**
   - Конечная точка: `/authorize`
   - Метод: `POST`
   - Тело запроса: JSON с полями `email` и `password`
   - Полученный токен необходимо передавать в хедере `Authorization-access`

2. **Регистрация пользователя**
   - Конечная точка: `/register`
   - Метод: `POST`
   - Тело запроса: JSON с полями `email` и `password`
   - Пароль должен содержать минимум 8 символов и хотя бы одну цифру
   - Пароль имеет статус `perfect` если содержит спецсимволы, строчные и заглавные буквы

3. **Проверка авторизации**
   - Конечная точка: `/feed`
   - Метод: `GET`
   - Access-token необходимо передавать в хедере `Authorization-access`

## Запуск Сервиса

### Использование Docker

1. Вы можете загрузить готовый Docker образ с Docker Hub
```bash
    docker pull yasminworks/authserver
```

2. Создание Docker образа
```bash
   docker build . -t authserver:latest
``` 

3. Создание Docker volume и запуск контейнера
```bash
   docker volume create auth-server
   docker run -d -it -p 8082:8082 -v auth-server:/app/storage authserver
```
   
### Компиляция Исходного Кода
0. Перейдите в корневую папку проекта

1. Установка зависимостей:
```bash
  go mod download
 ```

2. Подготовка базы данных:
```bash
   go run ./cmd/migrator --storage-path=./storage/storage.db --migrations-path=./migrations
 ```  
3. Компиляция и запуск:
 ```bash 
    go build -o auth-server ./cmd/auth-server/main.go
    ./auth-server
 ```  
### Использование Утилиты Task

Если у вас установлена утилита Task, можно запустить сервис командой
```bash
    task build
```
