#### SQL Мигратор
![Build Status](https://github.com/natkazb/sql-migrator/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/natkazb/sql-migrator)](https://goreportcard.com/report/github.com/natkazb/sql-migrator)
[![codecov](https://codecov.io/gh/natkazb/sql-migrator/branch/main/graph/badge.svg)](https://codecov.io/gh/natkazb/sql-migrator)

Как запускать: 
1) бд можно запустить с помощью docker compose
```
docker compose -f deployments/docker-compose.yaml up
```
2) запустить саму утилиту с конфигурацией подключения к БД
```
go run ./cmd/gomigrator --config=./configs/config.yaml
```

### Создать миграцию (CREATE)
```
go run ./cmd/gomigrator --config=./configs/config.yaml create <name>
```

### Создать GO миграцию (CREATE-GO)
```
go run ./cmd/gomigrator --config=./configs/config.yaml create-go <name>
```

### Применение всех миграций (UP)
```
go run ./cmd/gomigrator --config=./configs/config.yaml up
```

### Откат последней миграции (DOWN)
```
go run ./cmd/gomigrator --config=./configs/config.yaml down
```

### Повтор последней миграции (откат + накат) (REDO)
```
go run ./cmd/gomigrator --config=./configs/config.yaml redo
```

### Вывод версии базы (DBVERSION)
```
go run ./cmd/gomigrator --config=./configs/config.yaml dbversion
```

### Вывод статуса миграций (STATUS)
```
go run ./cmd/gomigrator --config=./configs/config.yaml status <limit>
```
