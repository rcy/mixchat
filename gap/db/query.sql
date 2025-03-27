-- name: CreateGuestUser :one
insert into users(guest, user_id) values(true, $1) returning *;

-- name: CreateUser :one
insert into users(user_id, username) values(@user_id, @username) returning *;

-- name: SessionUser :one
select users.* from sessions
join users on users.user_id = sessions.user_id
where sessions.expires_at > now()
and session_id = $1;

-- name: User :one
select * from users where user_id = $1;

-- name: UserByUsername :one
select * from users where username = $1;

-- name: CreateSession :one
insert into sessions(session_id, user_id) values($1, $2) returning session_id;

-- name: InsertEvent :one
insert into events(event_id, event_type, payload) values ($1, $2, $3) returning *;

-- name: Event :one
select * from events where event_id = $1;

-- name: PublicActiveStations :many
select * from stations where active = true and is_public = true;

-- name: Station :one
select * from stations where slug = $1;

-- name: CreateStation :one
insert into stations(station_id, slug, user_id, active) values($1, $2, $3, $4) returning *;

-- name: SetStationCurrentTrack :exec
update stations set current_track_id = $1 where station_id = $2;

-- name: SetStationHostPort :exec
update stations set host_port = $1 where station_id = $2;

-- name: StationCurrentTrack :one
select tracks.* from stations join tracks on stations.current_track_id = tracks.track_id where stations.station_id = $1;

-- name: CreateStationMessage :one
insert into station_messages(station_message_id, type, station_id, nick, body, parent_id) values($1, $2, $3, $4, $5, $6) returning *;

-- name: StationMessages :many
select * from station_messages where station_id = $1 and is_hidden = false order by station_message_id desc limit 500;

-- name: UpdateStationMessage :exec
update station_messages set type = $1, body = $2 where station_message_id = $3;

-- name: FindLastStationMessage :one
select * from station_messages where station_id = $1 order by created_at desc limit 1;

-- name: HideStationMessage :exec
update station_messages set is_hidden = true where station_message_id = $1;

-- name: TrackRequestStationMessage :one
select * from station_messages
where station_id = $1
and type = 'TrackRequested'
and parent_id = $2;

-- name: CreateTrack :one
insert into tracks(track_id, station_id, artist, title, raw_metadata, rotation)
values($1,$2,$3,$4,$5, (coalesce((select min(rotation) from tracks where station_id = $2), 0)))
returning *;

-- name: Track :one
select * from tracks where track_id = $1;

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

-- name: SetSearchStatusCompleted :exec
update searches set status = 'completed' where search_id = $1;

-- name: SetSearchStatusFailed :exec
update searches set status = 'failed' where search_id = $1;

-- name: CreateResult :exec
insert into results(result_id, search_id, station_id, extern_id, url, thumbnail, title, uploader, duration, views) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);

-- name: Search :one
select * from searches where search_id = $1;

-- name: Results :many
select * from results where search_id = $1;
