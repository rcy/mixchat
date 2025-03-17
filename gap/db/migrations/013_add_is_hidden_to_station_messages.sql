alter table station_messages add column is_hidden bool default false;

---- create above / drop below ----

alter table station_messages drop column is_hidden;
