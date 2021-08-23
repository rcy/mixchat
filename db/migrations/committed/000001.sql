--! Previous: -
--! Hash: sha1:6cfd6e8335077fdb58e63f5a37a67b9b59705e85

-- Enter migration here

drop table if exists events;
create table events (
       id serial,
       name text,
       data jsonb
);

drop table if exists results;
create table results (
       id serial,
       name text,
       data jsonb
);

drop function if exists trigger_job();
create function trigger_job() returns trigger as $$
begin
  perform graphile_worker.add_job(tg_argv[0], json_build_object(
    'schema', tg_table_schema,
    'table', tg_table_name,
    'op', tg_op,
    'id', (case when tg_op = 'delete' then old.id else new.id end)
  ));
  return new;
end;
$$ language plpgsql volatile;


create trigger insert_event
  after insert on events
  for each row
  execute procedure trigger_job('event_created');
