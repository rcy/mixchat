--! Previous: sha1:5a8d080ff2cd48e28b31bbf50bb753ec9810062b
--! Hash: sha1:9eda2afaaf05a96c0fc7ec7adb45a63599a7afbe

drop table if exists bucket;
drop table if exists plays;
drop table if exists skips;
drop table if exists tracks;

create table tracks (
       id serial primary key,
       filename text not null unique,
       created_at timestamptz not null default now(),
       bucket int not null, -- default bucket(),
       fuzz real not null default 0,
       event_id int references events not null
);

-- return minimum bucket value from tracks
create or replace function bucket() returns integer as $$
  select coalesce(( select bucket from tracks order by bucket asc limit 1), 0);
$$ language sql;

create or replace function set_track_due() returns trigger as $$
begin
  update tracks set bucket = bucket + 1, fuzz = random() where id = new.track_id;
  return new;
end
$$ language plpgsql volatile;

create table plays (
       id serial primary key,
       track_id integer references tracks not null,
       created_at timestamptz not null default now()
);
create trigger track_due_tg
  before insert on plays
  for each row
  execute procedure set_track_due();

create table skips (
       id serial primary key,
       track_id integer references tracks not null,
       created_at timestamptz not null default now()
);
create trigger track_due_tg
  before insert on skips
  for each row
  execute procedure set_track_due();
