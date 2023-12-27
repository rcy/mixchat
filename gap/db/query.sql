-- name: Stations :many
select * from stations;

-- name: Station :one
select * from stations where slug = $1;
