#postgres: docker-compose --env-file .env.dev up postgres
#icecast: docker-compose --env-file .env.dev up icecast
#liquidsoap: docker-compose --env-file .env.dev up liquidsoap
app: make -C app start
worker: make -C worker start
migrate: make -C db watch
frontend: make -C frontend start
