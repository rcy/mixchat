--! Previous: sha1:7daa13a33f501f0c7f046569b9fca32ca6d46427
--! Hash: sha1:7adf795257f3d4d0ac1cc3a0e788d529ecd7a04e

drop view if exists plays;
create view plays as
  select * from track_events;
