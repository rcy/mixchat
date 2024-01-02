-- name: InsertEvent :one
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
select * from station_messages where station_id = $1 order by station_message_id desc limit 500;

-- name: CreateTrack :one
insert into tracks(track_id, station_id, artist, title, raw_metadata, rotation)
values($1,$2,$3,$4,$5, (coalesce((select min(rotation) from tracks where station_id = $2), 0)))
returning *;

-- name: OldestUnplayedTrack :one
select * from tracks
where tracks.station_id = $1
and plays = 0
and rotation = (select min(rotation) from tracks where station_id = $1)
order by track_id asc
limit 1;

-- name: RandomTrack :one
select * from tracks
where tracks.station_id = $1
and plays > 0
and rotation = (select min(rotation) from tracks where station_id = $1)
order by random()
limit 1;

-- name: IncrementTrackRotation :exec
update tracks set rotation = rotation + 1 where track_id = $1;

-- name: IncrementTrackPlays :exec
update tracks set plays = plays + 1 where track_id = $1;

-- name: CreateSearch :exec
insert into searches(search_id, station_id, query) values($1,$2,$3);

-- name: CreateResult :exec
insert into results(result_id, search_id, station_id, extern_id, url, thumbnail, title, uploader, duration, views) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);

-- name: Search :one
select * from searches where search_id = $1;

-- name: Results :many
select * from results where search_id = $1;
