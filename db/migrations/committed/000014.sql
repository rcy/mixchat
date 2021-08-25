--! Previous: sha1:2d484815c3608861a900e0ef14882b33ef523279
--! Hash: sha1:82f22d211de814e7f41d654e0119bcada0977410

create or replace function update_track_bucket(track_id integer, count integer) returns integer as $$
  begin
    update tracks set bucket = bucket + count, fuzz = random() where id = track_id;
    return 1;
  end
$$ language plpgsql;

create or replace function set_track_due_played() returns trigger as $$
begin
  -- push the track forward 1 step
  perform update_track_bucket(new.track_id, 1);
  return new;
end
$$ language plpgsql volatile;

create or replace function set_track_due_skipped() returns trigger as $$
begin
  -- push the track forward as many times as its been skipped; every time it is skipped it gets pushed out one further 
  perform update_track_bucket(new.track_id, (select count(1) from plays where track_id = 3 and action = 'skipped')::integer);
  return new;
end
$$ language plpgsql volatile;

drop trigger if exists track_due_tg on plays;

drop trigger if exists queued_track_due_tg on plays;
create trigger queued_track_due_tg
  before insert on plays
  for each row
  when (new.action = 'queued')
  execute procedure set_track_due_played();

drop trigger if exists skipped_track_due_tg on plays;
create trigger skipped_track_due_tg
  before insert on plays
  for each row
  when (new.action = 'skipped')
  execute procedure set_track_due_skipped();

drop function if exists set_track_due();
