#postgres: docker-compose --env-file .env.dev up postgres
#icecast: docker-compose --env-file .env.dev up icecast
pg: docker logs --follow --since 0m  djfullmoon_postgres_1
liquidsoap: make -C liquidsoap start
app: make -C app start
worker: make -C worker start
migrate: make -C db watch
frontend: make -C frontend start
