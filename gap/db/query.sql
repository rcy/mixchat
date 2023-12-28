-- name: ActiveStations :many
select * from stations where active = true;

-- name: Station :one
select * from stations where slug = $1;

-- name: StationMessages :many
select * from messages where station_id = $1 order by id desc limit 100;

-- name: RecentPlays :many
select
        t.filename,
        te.created_at,
        coalesce(t.metadata->'common'->>'title', '')::text title,
        coalesce(t.metadata->'common'->>'artist', '')::text artist
from track_events te
join tracks t on t.id = te.track_id
where te.station_id = $1
  and te.action = 'played'
  and te.created_at > $2;

-- name: CurrentTrack :one
select
        coalesce(t.metadata->'common'->>'title', '')::text title,
        coalesce(t.metadata->'common'->>'artist', '')::text artist
from track_events te
join tracks t on t.id = te.track_id
where te.station_id = $1
  and te.action = 'played'
order by te.id desc
limit 1;
