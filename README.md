# Распределённый вычислитель арифметических выражений

## Описание проекта

Проект представляет собой распределённую систему для вычисления арифметических выражений. Она включает два компонента:

- **Оркестратор** – управляет выражениями, разбивает их на задачи, распределяет между агентами и сохраняет данные в SQLite.
- **Агент** – вычисляет задачи, полученные от оркестратора через gRPC.

Теперь система поддерживает регистрацию и вход пользователей. Все выражения вычисляются в контексте конкретного пользователя.

## Структура проекта

```bash
.
├── cmd/
│   ├── agent/                 # Вычислительный агент
│   │   └── main.go
│   └── calc_service/          # Оркестратор (сервер)
│       └── main.go
│
├── internal/
│   ├── evaluator/             # Логика выражений
│   │   └── evaluator.go
│   ├── handlers/              # HTTP-обработчики
│       ├── auth.go            # Регистрация и логин
│       └── calculate.go
│
│
├── pkg/
│   ├── errors/                # Кастомные ошибки
│   │   └── errors.go
│   └── auth/                  # Работа с JWT
│   │    └── jwt.go
│   └── db/               # Работа с SQLite
│       └── db.go
│
├── proto/                     # gRPC-сервисы
│   └── calculator.proto
│
├── static/                    # Frontend
├── tests/                     # Тесты
├── go.mod
├── go.sum
├── README.md
```

## Основной функционал

- Регистрация: `POST /api/v1/register`
- Вход: `POST /api/v1/login` → JWT
- Отправка выражения: `POST /api/v1/calculate`
- Список выражений: `GET /api/v1/expressions`
- Выражение по ID: `GET /api/v1/expressions/:id`
- Задача агенту (gRPC): `RequestTask`, `SubmitResult`

## Запуск

### 1. Клонируйте проект

```bash
git clone https://github.com/scriptoxin/yandex-liceum-go-calc.git
cd yandex-liceum-go-calc
```

### 2. Установите зависимости и сгенерируйте gRPC код

```bash
go mod tidy

# Сгенерировать код из proto:
# Установите buf: https://docs.buf.build/installation
buf generate proto
```

### 3. Настройте окружение

#### Оркестратор

```bash
export TIME_ADDITION_MS=5000
export TIME_SUBTRACTION_MS=5000
export TIME_MULTIPLICATION_MS=5000
export TIME_DIVISION_MS=5000
export JWT_SECRET=your-secret
export DB_PATH=./data.db

go run cmd/calc_service/main.go
```

#### Агент

```bash
export COMPUTING_POWER=4
export GRPC_ORCHESTRATOR_ADDR=localhost:9090

go run cmd/agent/main.go
```

### 4. Проверьте работу

```bash
curl -X POST http://localhost:8080/api/v1/register   -H 'Content-Type: application/json'   -d '{"login":"user1", "password":"pass123"}'

curl -X POST http://localhost:8080/api/v1/login   -H 'Content-Type: application/json'   -d '{"login":"user1", "password":"pass123"}'
# В ответ — JWT
```

## Тестирование

```bash
go test ./...
```

## Возможности интерфейса

- Ввод выражений и просмотр истории
- Обновление статусов в реальном времени
- Отображение результатов только текущего пользователя

## Масштабирование

- Несколько агентов → выше производительность
- gRPC повышает надёжность обмена
- SQLite сохраняет состояние между запусками
