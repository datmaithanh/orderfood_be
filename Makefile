postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14.19-alpine3.21

createdb: 
	docker exec -it postgres14 createdb --username=root --owner=root orderfood

dropdb:
	docker exec -it postgres14 dropdb orderfood

createschema:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/orderfood?sslmode=disable" -verbose up
	
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/orderfood?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run ./cmd/main/main.go

PHONY: postgres createdb dropdb createschema migrateup migratedown sqlc server