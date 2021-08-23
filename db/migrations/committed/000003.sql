--! Previous: sha1:cf4e125b5109d676dcad46f5f1976555bb703adb
--! Hash: sha1:c992241f4e0d8ceb48b96eb411b01e0b97e66586

-- Enter migration here
alter table events drop constraint if exists events_pkey;
alter table events add primary key (id);
