--! Previous: sha1:f48adc4c411e977728772e407c3bf9d1f828f3b5
--! Hash: sha1:a4bc7feda91572c8e8a384399bea7480159064dd

-- Enter migration here
alter table events alter column station_id drop default;
alter table track_events alter column station_id drop default;
alter table tracks alter column station_id drop default;
alter table results alter column station_id drop default;
