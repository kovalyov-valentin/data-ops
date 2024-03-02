# DATA-OPS SERVICE

## Общая информация:
Реализация сервиса для 

# Как поднять:

## 1. Запуск контейнера с PostgreSQL, Nats, Clickhouse, Redis:
```
make compose
```
## 2. Миграция таблиц в поднятой Postgres:
```
make upmigrate
```
## 3. Запуск сервера:
```
make run
```
## 4. Остановка контейнеров:
```
make stop
```
## TODO:
- Исправить архитектуру (приведение к общим интерфейсам, таким образом исправить проблему с Clickhouse и Nats)
  
## Стек технологий проекта:
* GO (gin)
* Postgres (sql)
* Redis
* Nats
* Clickhouse
* Docker, docker-compose
