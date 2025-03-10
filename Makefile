dev:
	godotenv -f .env go run cmd/server/main.go

parse:
	godotenv -f .env go run cmd/parser/main.go