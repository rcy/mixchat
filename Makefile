DC:=docker compose

start:
	foreman start

install:
	cd frontend && npm install
	cd worker && npm install
	cd app && npm install
	cd db && npm install

up: SERVICES=icecast postgres liquidsoap
up:
	${DC} --env-file .env.dev build ${SERVICES}
	${DC} --env-file .env.dev up ${SERVICES}

down:
	${DC} --env-file .env.dev down ${SERVICES}

stop:
	${DC} --env-file .env.dev stop ${SERVICES}

deploy: export DOCKER_HOST=ssh://ubuntu@stream.djfullmoon.com
deploy:
	echo ${DOCKER_HOST}
	${DC} --env-file .env.prod build
	${DC} --env-file .env.prod up -d icecast liquidsoap
	${DC} --env-file .env.prod up -d app
	${DC} --env-file .env.prod up -d worker

psql:
	psql postgres://djfm:djfm@localhost:54322/djfullmoon_development
