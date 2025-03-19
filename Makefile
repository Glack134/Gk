build:
	docker-compose build message

run:
	docker-compose up message

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwe230405180405@0.0.0.0:5432'/postgres?sslmode=disable' up

swag:
	swag init -g cmd/main.go
