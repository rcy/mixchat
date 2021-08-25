--! Previous: sha1:8c7fe3d973334366575783a998a74f29bfda5f42
--! Hash: sha1:2d484815c3608861a900e0ef14882b33ef523279

-- only update the track due time for a track when its queued, not on every play action

drop trigger if exists track_due_tg on plays;

create trigger track_due_tg
  before insert on plays
  for each row
  when (new.action = 'queued')
  execute procedure set_track_due();
