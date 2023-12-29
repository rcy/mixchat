gen: schema
	`go env GOPATH`/bin/sqlc generate

dump:
	docker exec -it mixchat-postgres-1 pg_dump --dbname=mixchat_development --user=postgres --table=public.* --schema-only

schema:
	${MAKE} -s dump > schema.sql
