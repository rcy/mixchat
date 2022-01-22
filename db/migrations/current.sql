-- Enter migration here

-- there were two unique filename constraints
alter table tracks drop constraint if exists tracks_filename_key;
alter table tracks drop constraint if exists tracks_unique_filename;

alter table tracks add constraint tracks_unique_station_filename unique (station_id, filename);
