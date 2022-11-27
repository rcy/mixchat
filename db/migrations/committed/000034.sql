--! Previous: sha1:f27f94ee8f6071acc52775f13571b0710ef51369
--! Hash: sha1:98e0a7223e516cebcaf38f1504c37da16b9202af

drop trigger if exists broadcast_track_events on tracks;



drop trigger if exists broadcast_track on tracks;
create trigger broadcast_track
  after insert on tracks
  for each row
  execute procedure trigger_job('broadcast_track');
