# Simple Makefile for a Go project

DB_USER ?= go_with_tools_auth_user
DB_NAME ?= go_with_tools_db
SCHEMA  := internal/database/schema.sql

db-schema:
	pg_dump -U $(DB_USER) -d $(DB_NAME) --schema-only -f $(SCHEMA)
	sed -i '/^\\restrict\b/d; /^\\unrestrict\b/d' $(SCHEMA)

sqlc: db-schema
	sqlc generate

proto:
	rm -r ./gen/* && protoc --go_out=./gen --go_opt=paths=source_relative --go-grpc_out=./gen --go-grpc_opt=paths=source_relative auth.proto

.PHONY: all build run test clean watch up logs down prod-up prod-logs prod-down itest db-schema sqlc swag migrate-in-docker proto
