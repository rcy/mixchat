up:
	docker-compose build
	docker-compose --env-file .env.dev up icecast liquidsoap

deploy: export DOCKER_HOST=ssh://ubuntu@djfullmoon.com
deploy:
	echo ${DOCKER_HOST}
	docker-compose build
	docker-compose --env-file .env.prod up -d
