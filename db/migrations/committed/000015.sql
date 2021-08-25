--! Previous: sha1:82f22d211de814e7f41d654e0119bcada0977410
--! Hash: sha1:04bfbf7f2114cf394bf1943f412533df70b08669

create or replace function set_track_due_skipped() returns trigger as $$
begin
  -- push the track forward as many times as its been skipped; every time it is skipped it gets pushed out one further
  perform update_track_bucket(
    new.track_id,
    (
      select count(1)
        from plays
      where
        track_id = new.track_id
      and
        action = 'skipped'
    )::integer
  );
  return new;
end
$$ language plpgsql volatile;
