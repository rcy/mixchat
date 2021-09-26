--! Previous: sha1:68c0fc8d2e01aae007a7f1f56cc433d824898148
--! Hash: sha1:73787d80468c2d9623ce5fc0a0d2559d2543eab2

-- add station_id to tables

alter table events drop constraint if exists events_station_id_fkey;
alter table events drop column if exists station_id;
alter table events add column station_id int;
alter table events add constraint events_station_id_fkey foreign key (station_id) references stations;

alter table results drop constraint if exists results_station_id_fkey;
alter table results drop column if exists station_id;
alter table results add column station_id int;
alter table results add constraint results_station_id_fkey foreign key (station_id) references stations;

alter table track_events drop constraint if exists track_events_station_id_fkey;
alter table track_events drop column if exists station_id;
alter table track_events add column station_id int;
alter table track_events add constraint track_events_station_id_fkey foreign key (station_id) references stations;

alter table tracks drop constraint if exists tracks_station_id_fkey;
alter table tracks drop column if exists station_id;
alter table tracks add column station_id int;
alter table tracks add constraint tracks_station_id_fkey foreign key (station_id) references stations;
