--! Previous: sha1:bcbfb6043066ee52c2ce44e0faafeea2919556ad
--! Hash: sha1:6371b5fc85326bd78c55e19119042c4f721bf56b

-- Enter migration here

alter table tracks drop column if exists metadata;
alter table tracks add column metadata jsonb;

drop trigger if exists update_track_metadata on tracks;
create trigger update_track_metadata
  after insert on tracks
  for each row
  execute procedure trigger_job('update_track_metadata');
