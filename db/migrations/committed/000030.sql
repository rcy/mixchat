--! Previous: sha1:7399968c0d242af7c14e1c6d2a8c04023848f356
--! Hash: sha1:bcbfb6043066ee52c2ce44e0faafeea2919556ad

-- Enter migration here

-- there were two unique filename constraints
alter table tracks drop constraint if exists tracks_filename_key;
alter table tracks drop constraint if exists tracks_unique_filename;

-- replace with this one, so it's per station
alter table tracks drop constraint if exists tracks_unique_station_filename;
alter table tracks add constraint tracks_unique_station_filename unique (station_id, filename);
