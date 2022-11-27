--! Previous: sha1:59d7fba647b78599bd4dffc78ffdfbe9c1808aca
--! Hash: sha1:f27f94ee8f6071acc52775f13571b0710ef51369

drop trigger if exists broadcast_track_event on track_events;

create trigger broadcast_track_event
  after insert on track_events
  for each row
--  when (new.action = 'played')
  execute procedure trigger_job('broadcast_track_event');
