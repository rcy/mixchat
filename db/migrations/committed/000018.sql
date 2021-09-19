--! Previous: sha1:7adf795257f3d4d0ac1cc3a0e788d529ecd7a04e
--! Hash: sha1:ff3678a44a4af503d938625fa29b9517315c2713

drop view if exists skips;
create view skips as
  select id, track_id, created_at as ts
    from track_events
  where
    action = 'skipped';
