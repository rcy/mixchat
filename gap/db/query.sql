-- name: CreateEvent :one
insert into events(event_id, event_type, payload) values ($1, $2, $3) returning *;

-- name: Event :one
select * from events where event_id = $1;

-- name: ActiveStations :many
select * from stations where active = true;

-- name: Station :one
select * from stations where slug = $1;

-- name: CreateStation :one
insert into stations(station_id, slug, active) values($1, $2, $3) returning *;

-- name: CreateStationMessage :one
insert into station_messages(station_message_id, type, station_id, nick, body, parent_id) values($1, $2, $3, $4, $5, $6) returning *;

-- name: StationMessages :many
select * from station_messages where station_id = $1 order by station_message_id desc limit 5;

-- -- name: StationMessages :many
-- select * from messages where station_id = $1 order by id desc limit 100;

-- -- name: RecentPlays :many
-- select
--         t.filename,
--         te.created_at,
--         coalesce(t.metadata->'common'->>'title', '')::text title,
--         coalesce(t.metadata->'common'->>'artist', '')::text artist
-- from track_events te
-- join tracks t on t.id = te.track_id
-- where te.station_id = $1
--   and te.action = 'played'
--   and te.created_at > $2;

-- -- name: CurrentTrack :one
-- select
--         coalesce(t.metadata->'common'->>'title', '')::text title,
--         coalesce(t.metadata->'common'->>'artist', '')::text artist
-- from track_events te
-- join tracks t on t.id = te.track_id
-- where te.station_id = $1
--   and te.action = 'played'
-- order by te.id desc
-- limit 1;
