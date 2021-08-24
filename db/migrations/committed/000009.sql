--! Previous: sha1:9a9e7a0fe44b4573acf01d116b29aac9b7b49a3e
--! Hash: sha1:79b18c70c69fbcc9831c8a615526fadc0ea04b7b

-- timestamp events

alter table events drop column if exists created_at;
alter table events add column created_at timestamptz default now();

update events set created_at = '2020-08-01T00:00Z';

alter table events alter column created_at set not null;
