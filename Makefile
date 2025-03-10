dev:
	godotenv -f .env go run cmd/server/main.go

docker:
	godotenv -f .env docker compose up -d --build

parse:
	godotenv -f .env go run cmd/parser/main.go