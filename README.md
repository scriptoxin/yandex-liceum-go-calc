# yandex-liceum-go-calc

Это API-сервис для выполнения арифметических вычислений с использованием математических выражений. Программа поддерживает базовые операции: сложение, вычитание, умножение и деление.

## Описание

API принимает арифметическое выражение и возвращает результат. Если выражение некорректно или произошла ошибка на сервере, будет возвращен соответствующий код ошибки и сообщение.

## Установка

### Клонирование репозитория

```bash
git clone https://github.com/scriptoxin/yandex-liceum-go-calc.git
cd yandex-liceum-go-calc
```

### Настройка порта

Вы можете настроить порт через переменную окружения **PORT**. По умолчанию используется порт 8080. Пример для использования другого порта:

```bash
export PORT=8081
```

## Запуск приложения

Запустите сервер командой:

```bash
go run ./cmd/calc_service/main.go
```

## Примеры использования

### Успешный запрос

Отправка POST-запроса с математическим выражением:

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```

Ответ

```json
{
  "result": "6"
}
```

### Ошибка 422: Некорректное выражение

Отправка некорректного выражения:

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2++2"
}'
```

Ответ

```json
{
  "error": "Expression is not valid"
}
```

### Ошибка 500: Внутренняя ошибка сервера

Если произошла ошибка на сервере:

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "1/0"
}'
```

Ответ

```json
{
  "error": "Internal server error"
}
```
