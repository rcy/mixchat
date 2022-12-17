#export DATABASE_URL?=postgres://postgres:postgres@localhost:5432/djfullmoon_development

start:
	npx nodemon -r dotenv/config worker.js dotenv_config_path=../.env.dev


build:
	docker build --tag quickiron6536/mixchat-worker .

push:
	docker push quickiron6536/mixchat-worker
