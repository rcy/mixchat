app: sleep 10 && make -C app start
worker: sleep 10 && make -C worker start
postgres: docker-compose --env-file .env.dev up postgres
icecast: docker-compose --env-file .env.dev up icecast
liquidsoap: docker-compose --env-file .env.dev up liquidsoap
