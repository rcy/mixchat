--! Previous: sha1:79b18c70c69fbcc9831c8a615526fadc0ea04b7b
--! Hash: sha1:3575d240fbf0dd7c9406e6b2c13b01a63ffa5ae5

-- add track_changes table
-- latest record in table represents currently playing track

drop table if exists track_changes;
create table track_changes (
  id serial primary key,
  track_id int references tracks not null,
  created_at timestamptz not null default now()
);
