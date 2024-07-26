DB_URL=postgresql://root:secret@postgres12:5432/simple_bank?sslmode=disable
postgres:
	docker run --name postgres12 --network simplebank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

pgAdmin:
	docker run --name pgadmin-container -p 5050:80 -e PGADMIN_DEFAULT_EMAIL=kaitok57a01@gmail.com -e PGADMIN_DEFAULT_PASSWORD=secret -d dpage/pgadmin4

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

mock:
	mockgen -package mockdb  -destination db/mock/store.go simple-bank/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migratedown db_docs db_schema sqlc test server pgAdmin mock
