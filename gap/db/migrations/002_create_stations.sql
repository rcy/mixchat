create table stations(
  station_id text primary key,
  created_at timestamptz not null default now(),
  slug text not null unique,
  name text not null default '',
  active boolean not null
);

---- create above / drop below ----

drop table stations;
