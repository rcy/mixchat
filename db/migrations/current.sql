drop table if exists messages;
create table messages (
    id integer NOT NULL,
    station_id integer not null,
    message text,
    nick text
);
