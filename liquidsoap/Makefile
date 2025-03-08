# include ../.env.dev

# start: export API_BASE := http://host.docker.internal:5500
# start: export ICECAST_HOST := localhost
# start: export ICECAST_PORT := 8000
# start: export ICECAST_SOURCE_PASSWORD := ${ICECAST_SOURCE_PASSWORD}
# start: export LIQUIDSOAP_BROADCAST_PASSWORD := ${LIQUIDSOAP_BROADCAST_PASSWORD}
# start: export STATION_SLUG := dummy
# start: export TELNET_PORT := 1234
# start:
# 	liquidsoap --force-start ./start.liq

docker-build:
	docker buildx build --platform linux/amd64 --tag rcy0/mixchat-liquidsoap .

docker-push:
	docker push rcy0/mixchat-liquidsoap

dstart:
	docker run -p 1234:1234 --env-file=../.env.dev rcy0/mixchat-liquidsoap
