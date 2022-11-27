--! Previous: sha1:6371b5fc85326bd78c55e19119042c4f721bf56b
--! Hash: sha1:59d7fba647b78599bd4dffc78ffdfbe9c1808aca

drop table if exists messages;
create table messages (
    id serial primary key,
    station_id integer references stations not null,
    body text,
    nick text,
    created_at timestamptz default now() not null
);

drop trigger if exists insert_message_notify on messages;

drop function if exists trigger_notify_insert_station_relation_row();
create function trigger_notify_insert_station_relation_row() returns trigger as $$
begin
  perform pg_notify(
    'postgraphile:station:' || new.station_id || ':' || tg_argv[0],
    json_build_object(
      '__node__',
      json_build_array(tg_argv[0], new.id)
    )::text
  );
  return new;
end
$$ language plpgsql volatile;

create trigger insert_message_notify
  after insert on messages
  for each row
  execute procedure trigger_notify_insert_station_relation_row('messages');
