alter table stations add column current_track_id text references tracks;
---- create above / drop below ----
alter table stations drop column current_track_id;
