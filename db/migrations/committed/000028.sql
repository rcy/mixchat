--! Previous: sha1:ccef02b48ac1babfff4e5ec505b850120e9a0cb9
--! Hash: sha1:c692b7f918e3cc7575bcf3887ef4cdf0d1771f68

create or replace function set_track_due_yeeted() returns trigger as $$
begin
  -- push the track forward a lot of steps, effectively taking it out of rotation
  perform update_track_bucket(new.track_id, 1000000);
  return new;
end
$$ language plpgsql volatile;

drop trigger if exists yeeted_track_due_tg on track_events;
create trigger yeeted_track_due_tg
  before insert on track_events
  for each row
  when (new.action = 'yeeted')
  execute procedure set_track_due_yeeted();
