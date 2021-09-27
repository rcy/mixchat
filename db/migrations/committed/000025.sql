--! Previous: sha1:ee89cfe256c8e066146d095a03ad703cbffa8949
--! Hash: sha1:f48adc4c411e977728772e407c3bf9d1f828f3b5

-- add station irc channels

drop table if exists irc_channels;
create table irc_channels (
       id serial primary key,
       station_id integer references stations not null unique,
       created_at timestamptz not null default now(),
       server text not null default 'irc.libera.chat',
       channel text not null unique
);
