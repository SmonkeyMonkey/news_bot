.PHONY: up
up:
	cd internal/storage/migrations && \
	goose postgres "postgres://postgres:postgres@localhost:5432/news_bot?sslmode=disable" up

.PHONY: down
down:
	cd internal/storage/migrations && \
	goose postgres "postgres://postgres:postgres@localhost:5432/news_bot?sslmode=disable" down && \
	goose postgres "postgres://postgres:postgres@localhost:5432/news_bot?sslmode=disable" down