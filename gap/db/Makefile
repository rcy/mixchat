gen: schema
	`go env GOPATH`/bin/sqlc generate

dump:
	docker exec -it gap-pg-1 pg_dump --dbname=mixchat_development --user=postgres --table=public.* --schema-only

schema:
	${MAKE} -s dump > schema.sql

migrate:
	${MAKE} -C migrations migrate

rollback:
	${MAKE} -C migrations rollback
