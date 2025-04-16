#### SQL Мигратор


Как запускать: 
1) бд можно запустить с помощью docker compose
`docker compose -f deployments/docker-compose.yaml up`
2) запустить саму утилиту
`go run ./cmd/gomigrator --config=./configs/config.yaml`
