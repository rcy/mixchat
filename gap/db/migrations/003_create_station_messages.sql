create table station_messages(
  station_message_id text primary key,
  created_at timestamptz not null default now(),
  type text not null,
  station_id text not null,
  parent_id text not null,
  nick text not null,
  body text not null
);

---- create above / drop below ----

drop table station_messages;
