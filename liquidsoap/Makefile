include ../.env.dev

start: export API_BASE := ${API_BASE}
start: export ICECAST_HOST := localhost
start: export ICECAST_PORT := 8010
start: export ICECAST_SOURCE_PASSWORD := ${ICECAST_SOURCE_PASSWORD}
start: export LIQUIDSOAP_BROADCAST_PASSWORD := ${LIQUIDSOAP_BROADCAST_PASSWORD}
start:
	liquidsoap ./start.liq
