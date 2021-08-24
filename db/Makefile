export DATABASE_URL?=postgres://djfm:djfm@localhost:54322/djfullmoon_development
export SHADOW_DATABASE_URL?=postgres://djfm:djfm@localhost:54322/djfullmoon_shadow_development
export ROOT_DATABASE_URL?=postgres://postgres:postgres@localhost:54322/postgres

watch:
	npx graphile-migrate watch

commit:
	npx graphile-migrate commit

uncommit:
	npx graphile-migrate uncommit

migrate:
	npx graphile-migrate migrate

create:
	psql postgres://postgres:postgres@localhost:5432/postgres -f create.sql

init:
	npx graphile-migrate init

