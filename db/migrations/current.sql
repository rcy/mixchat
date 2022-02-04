drop table if exists messages;
create table messages (
    id serial,
    station_id integer references stations not null,
    body text,
    nick text
);
