# Marketplace — REST API сервис объявлений

## Описание

Marketplace — это бэкенд-приложение на Go, реализующее базовую функциональность интернет-доски объявлений.
Пользователи могут регистрироваться, авторизовываться и размещать объявления с изображением и ценой.
API использует JWT для авторизации, SQLite как хранилище и построено по REST-архитектуре.

---

## Структура проекта

```
Marketplace/
├── cmd/                            
│   └── main.go                     # Точка входа в приложение. Инициализирует зависимости через fx, запускает HTTP-сервер
├── config/                         
│   └── local.yaml                  # YAML-файл конфигурации
├── internal/
│   ├── app/                    
│       └── jwt_model.go            # Структуры запросов/ответов для JWT
│       └── jwt_service_test.go     # Реализация логики генерации и валидации JWT-токенов
│       └── jwt_service.go          # Юнит-тесты для JWT-сервиса
│       └── market_interface.go     # Интерфейс для MarketService
│       └── market_model.go         # Модель объявления, параметры фильтрации, структура ответа
│       └── market_service_test.go  # Бизнес-логика работы с объявлениями
│       └── market_service.go       # Юнит-тесты для сервиса объявлений
│       └── mock_market_model.go    # Мок реализации MarketServicer для тестирования
│       └── mock_user_model.go      # Мок реализация UserRepository для тестирования
│       └── user_interface.go       # Интерфейс UserService
│       └── user_model.go           # Модель пользователя, структура регистрации
│       └── user_service_test.go    # Бизнес-логика регистрации, входа и валидации
│       └── user_service.go         # Юнит-тесты для логики работы с пользователем
│   ├── config/                     
│       └── config_model.go         # Структура конфигурации
│       └── config_service_test.go  # Парсинг конфигурационного файла и валидация
│       └── config_service.go       # Юнит-тесты загрузки и валидации конфига
│   ├── datasource/                 
│       └── tests/
│           └── market_repo_test.go # Интеграционные тесты для MarketRepo
│           └── user_repo_test.go   # Интеграционные тесты для UserRepo
│       └── db_service.go           # Инициализация SQLite-соединения
│       └── market_db.go            # Реализация репозитория объявлений
│       └── user_db.go              # Реализация репозитория пользователей
│   ├── di/                         
│       └── service.go              # Настройка зависимостей через fx
│   └── web/                        
│       └── market_handler_test.go  # Юнит-тесты эндопинтов объявлений
│       └── market_handler.go       # Реализация эндпоинтов объявлений
│       └── midlware_test.go        # Юнит-тесты middleware авторизации
│       └── midlware.go             # Middleware авторизации: обязательной и опциональной
│       └── router.go               # Настройка роутера (маршрутов), подключение middleware
│       └── user_handler_test.go    # Юнит-тесты эндпоинтов юзера
│       └── user_handler.go         # Реализация эндпоинтов юзера
└── storage/
    └── marketplace.db              # SQLite база данных

```

---

## Запуск

1. **Соберите проект:**
   ```sh
   go build -o server ./cmd
   ```

2. **Укажите путь к конфигу:**
   ```sh
   export CONFIG_PATH=./config/local.yaml
   ```

3. **Запустите сервис:**
   ```sh
   ./server
   ```

4. **Сервис будет доступен на порту, указанном в конфиге (по умолчанию 8080).**

---

## Основные возможности

- **Регистрация и авторизация пользователей (JWT)**
- **Валидация логина и пароля (по правилам из YAML)**
- **Создание объявлений (зарегестрированному пользователю)**
- **Валидация объявления (по правилам из YAML)**
- **Получение списка объявлений (с авторизацией и без)**
- **Валидация JWT для всех защищённых эндпоинтов**

---

## Примеры использования API

### 1. Регистрация и вход пользователя

```http
POST /register

Content-Type: application/json

{
  "login": "TestUser",
  "password": "Password1"
}
```

```http
POST /login

Content-Type: application/json

{
  "login": "TestUser",
  "password": "Password1"
}
```

### 2. Обновление JWT токена

```http
POST /refresh-access-token
Content-Type: application/json

{
    "refresh_token": "..."
}
```

### 3. Создание объявления

```http
POST /new-ad
Content-Type: application/json

{
	"title": "Самокат",
	"description": "Самокат зеленый пользовались 1 год",
	"image_url": "samokat.jpg",
	"price": 100000
}
```

### 4. Получение объявления (доступно с Authorization: Bearer <access_token> и без)

```http
GET /ads-list?page=1&limit=10&sort_by=price&order=asc&min_price=100&max_price=5000
```

**params:** 
`page` - int
`limit` - int
`sort_by` - price/date
`order` - asc/desc
`min_price` - int
`max_price` - int


---

## Пример конфига (`config/local.yaml`)

```yaml
env: local # or prod
http_port: 8080
db: "./storage/marketplace.db"
JWT_ACCESS_SECRET: "1234567890abcdef1234567890abcdef"
JWT_REFRESH_SECRET: "1234567890abcdef1234567890abcdef"
JWT_EXP_ACCESS_TOKEN: 15 # minutes
JWT_EXP_REFRESH_TOKEN: 24 # hours
username:
    min_length: 3
    max_length: 20
    allowed_characters: "A-Za-z0-9_-"
    case_insensitive: true
password:
    min_length: 8
    max_length: 64
    require_upper: true
    require_lower: true
    require_digit: true
ad:
    min_length_title: 3
    max_length_title: 100
    min_length_description: 10
    max_length_description: 1000
    img_type: 
        - jpg
        - jpeg
        - png
        - webm
    price_min: 0.01
```

---

## Тестирование

```sh
go test ./...
```

---

## Используемые Технологии

- **SQLite + github.com/mattn/go-sqlite3** (БД)
- **chi** (роутер)
- **zap** (логирование)
- **fx** (DI фреймворк)
- **bcrypt** (хэширование паролей)
- **uuid** (генерация UUID)