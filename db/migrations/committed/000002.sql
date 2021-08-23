--! Previous: sha1:6cfd6e8335077fdb58e63f5a37a67b9b59705e85
--! Hash: sha1:cf4e125b5109d676dcad46f5f1976555bb703adb

-- Enter migration here

drop trigger if exists insert_result on results;

drop function if exists trigger_notify();
create function trigger_notify() returns trigger as $$
begin
  perform pg_notify(tg_argv[0], (case when tg_op = 'delete' then old.id else new.id end)::text);
  return new;
end;
$$ language plpgsql volatile;

create trigger insert_result
  after insert on results
  for each row
  execute procedure trigger_notify('result');
