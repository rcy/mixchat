--! Previous: sha1:b27d41ff8d74ef86b9440b3b5a46d7e4fb502cf9
--! Hash: sha1:6e22a654986ec51a45fad20327537a7f043c27be

-- Enter migration here
drop view if exists recently_added;
create view recently_added as
select 
    tracks.id,
    tracks.filename,
    tracks.created_at as added_at,
    tracks.event_id,
    events.name,
    events.data
from tracks
join events on tracks.event_id = events.id
order by tracks.created_at desc
limit 10;
