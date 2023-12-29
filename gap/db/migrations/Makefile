include ../../.env

TERNCMD=`go env GOPATH`/bin/tern

echo:
	env

new:
	${TERNCMD} new "migration"

migrate:
	${TERNCMD} migrate

rollback:
	${TERNCMD} migrate --destination -1

.PHONY: new migrate rollback
