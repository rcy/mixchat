up:
	docker-compose build
	docker-compose --env-file .env.dev up icecast liquidsoap postgres

deploy: export DOCKER_HOST=ssh://ubuntu@djfullmoon.com
deploy:
	echo ${DOCKER_HOST}
	docker-compose --env-file .env.prod build
	docker-compose --env-file .env.prod up -d icecast liquidsoap
	docker-compose --env-file .env.prod up -d app
	docker-compose --env-file .env.prod up -d worker

psql:
	psql postgres://djfm:djfm@localhost:5432/djfullmoon_development
