up:
	docker-compose build
	docker-compose --env-file .env.dev up icecast liquidsoap

deploy: DOCKER_HOST=ubuntu@radio.nonzerosoftware.com
deploy:
	docker-compose build
	docker-compose --env-file .env.prod up
