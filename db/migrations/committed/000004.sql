--! Previous: sha1:c992241f4e0d8ceb48b96eb411b01e0b97e66586
--! Hash: sha1:5a8d080ff2cd48e28b31bbf50bb753ec9810062b

-- Enter migration here

delete from results;
alter table results drop column if exists event_id;
alter table results add column event_id integer not null;
alter table results add constraint fk_results_event foreign key(event_id) references events;
