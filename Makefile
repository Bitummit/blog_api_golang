PHONY: migrate
migrate:
	goose -dir migrations/ postgres "postgresql://postgres:postgres@127.0.0.1:5432/blog?sslmode=disable" up
