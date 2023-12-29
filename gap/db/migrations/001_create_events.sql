create table events(
       event_id text primary key,
       event_type text not null,
       created_at timestamptz not null default now(),
       payload jsonb not null       
);

create or replace function notify_event_insert()
returns trigger as $$
begin
  perform pg_notify('event_inserted', new.event_id::text);
  return new;
end;
$$ language plpgsql;

create trigger events_after_insert
after insert on events
for each row
execute function notify_event_insert();

---- create above / drop below ----
drop table events;
drop function notify_event_insert();
