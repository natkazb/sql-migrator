#### SQL Мигратор

Как запускать: 
1) бд можно запустить с помощью docker compose
`docker compose -f deployments/docker-compose.yaml up`
2) запустить саму утилиту с конфигурацией подключения к БД
`go run ./cmd/gomigrator --config=./configs/config.yaml`

### Создать миграцию (CREATE)
`go run ./cmd/gomigrator --config=./configs/config.yaml create name`
