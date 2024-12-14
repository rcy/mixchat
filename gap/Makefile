watch:
	DATABASE_URL=postgres://postgres:postgres@localhost:6543/mixchat_development PORT=5500 air

build:
	go build -o main main.go

run:
	go run .

deploy:
	fly deploy

up:
	PGPORT=6543 docker compose up -d

stop:
	docker compose stop

down:
	docker compose down

sql:
	psql postgres://postgres:postgres@localhost:6543/mixchat_development
