#include ../../.env

TERNCMD=`go env GOPATH`/bin/tern

migrate:
	${TERNCMD} migrate

rollback:
	${TERNCMD} migrate --destination -1

new:
	${TERNCMD} new "migration"

.PHONY: new migrate rollback
