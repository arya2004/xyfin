DB_URL=postgresql://root:secret@localhost:5432/xyfin?sslmode=disable

network:
	docker network create xyfin-network

postgres:
	docker run --name postgres --network xyfin-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

run:
	docker start a8a08f8c7f44 

mysql:
	docker run --name mysql8 -p 3306:3306  -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

createdb:
	docker exec -it postgres createdb --username=root --owner=root xyfin

dropdb:
	docker exec -it postgres dropdb xyfin

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb  -destination db/mock/store.go  github.com/arya2004/xyfin/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/arya2004/xyfin/worker TaskDistributor

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: network postgres run createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc test server mock proto evans redis