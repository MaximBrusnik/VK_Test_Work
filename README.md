# PubSub Сервис

Сервис публикации/подписки (PubSub) с поддержкой gRPC, реализованный на Go.

## Особенности

- Чистая архитектура
- gRPC API
- Поддержка множественных подписчиков
- Асинхронная обработка сообщений
- Graceful shutdown
- Валидация данных
- Логирование

## Структура проекта

```
.
├── cmd/                    # Точки входа в приложение
│   └── server/            # gRPC сервер
├── internal/              # Внутренний код приложения
│   ├── domain/           # Доменная логика
│   │   ├── entity/       # Доменные сущности
│   │   ├── repository/   # Интерфейсы репозиториев
│   │   └── service/      # Доменные сервисы
│   └── pubsub/           # Реализация PubSub
│       ├── delivery/     # Транспортный слой
│       │   └── grpc/     # gRPC обработчики
│       └── usecase/      # Сценарии использования
├── pkg/                   # Публичные пакеты
│   ├── logger/           # Логирование
│   ├── proto/            # Proto файлы
│   └── validator/        # Валидация
└── scripts/              # Скрипты
    └── generate_proto.sh # Генерация proto файлов
    └── generate_proto.bat
    └── generate_proto.ps1
```

## Требования

- Go 1.21 или выше
- Protocol Buffers
- gRPC

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/MaximBrusnik/VK_Test_Work
cd pubsub
```

2. Установите зависимости:
```bash
go mod download
```

3. Сгенерируйте proto файлы:
```bash
./scripts/generate_proto.sh
```
для скрипта generate_proto.sh не проверял. Только для generate_proto.bat

## Запуск

```bash
go run cmd/server/main.go
```

## API

### gRPC

Сервис предоставляет следующие gRPC методы:

- `Subscribe` - подписка на события
- `Publish` - публикация события

## Тестирование

```bash
go test ./...
```
