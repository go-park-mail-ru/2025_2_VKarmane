<div align="center">

# VKармане

Бэкенд проекта Дзен-мани 💰

</div>

## Запуск проекта

```bash
go run cmd/api/main.go
```

## Тестирование

### Запуск всех тестов
```bash
go test ./...
```

### Покрытие тестами
```bash
# Общее покрытие
go test -cover ./...

# Покрытие без моков (только бизнес-логика)
go test -cover ./internal/... | grep -v "mock.go"

# Детальное покрытие по функциям
go tool cover -func=coverage.out
```

## Разработка

### Линтинг и проверка кода
```bash
# Установка линтеров
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Запуск линтеров
golangci-lint run --config .golangci.yml
```

### Генерация моков
```bash
# Установка mockery
go install github.com/vektra/mockery/v2@latest

# Генерация моков
mockery --config .mockery.yaml
```
