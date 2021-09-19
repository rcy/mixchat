--! Previous: sha1:ff3678a44a4af503d938625fa29b9517315c2713
--! Hash: sha1:b27d41ff8d74ef86b9440b3b5a46d7e4fb502cf9

drop view if exists recently_played;
create view recently_played as
  select track_events.id,
         track_id,
         tracks.created_at,
         track_events.created_at as played_at,
         tracks.filename,
         event_id,
         action
  from track_events
  join tracks on track_events.track_id = tracks.id
  where action = 'played'
        and track_events.created_at > now() - interval '1' hour
  order by track_events.created_at desc
  limit 100;
