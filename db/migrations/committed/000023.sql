--! Previous: sha1:73787d80468c2d9623ce5fc0a0d2559d2543eab2
--! Hash: sha1:a98242df51f2930db72ffe94984c7097503a27af

alter table events alter column station_id drop not null;
update events set station_id = null;

alter table results alter column station_id drop not null;
update results set station_id = null;

alter table track_events alter column station_id drop not null;
update track_events set station_id = null;

alter table tracks alter column station_id drop not null;
update tracks set station_id = null;

delete from stations;
insert into stations (id, slug) values (1, 'emb');

update events set station_id = 1;
alter table events alter column station_id set not null;
alter table events alter column station_id set default 1;

update results set station_id = 1;
alter table results alter column station_id set not null;
alter table results alter column station_id set default 1;

update track_events set station_id = 1;
alter table track_events alter column station_id set not null;
alter table track_events alter column station_id set default 1;

update tracks set station_id = 1;
alter table tracks alter column station_id set not null;
alter table tracks alter column station_id set default 1;
