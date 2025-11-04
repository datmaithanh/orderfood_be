postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14.19-alpine3.21

createdb: 
	docker exec -it postgres14 createdb --username=root --owner=root orderfood

dropdb:
	docker exec -it postgres14 dropdb orderfood

createschema:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgresql://neondb_owner:npg_DRjtm1F7ThoI@ep-round-violet-a1nlcl0m.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&options=endpoint%3Dep-round-violet-a1nlcl0m" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://neondb_owner:npg_DRjtm1F7ThoI@ep-round-violet-a1nlcl0m.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&options=endpoint%3Dep-round-violet-a1nlcl0m" -verbose down

sqlc:
	sqlc generate

run:
	go run ./cmd/main/main.go

build:
	./build_linux_amd64.sh

PHONY: postgres createdb dropdb createschema migrateup migratedown sqlc run build