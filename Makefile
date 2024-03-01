compose:
	docker compose up -d
upmigrate:
	migrate -path migrations/postgres -database 'postgres://hezzl:password@localhost:5040/hezzldb?sslmode=disable' up
downmigrate:
	migrate -path migrations/postgres -database 'postgres://hezzl:password@localhost:5040/hezzldb?sslmode=disable' down
run:
	go run cmd/data-ops/main.go
stop:
	docker stop hezzl-service
	docker stop redis
	docker stop nats
	docker stop clickhouse
clean:
	docker rm hezzl-service
	docker rm redis
	docker rm nats
	docker rm clickhouse

.PHONY: compose upmigrate downmigrate stop clean run