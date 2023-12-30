-- create table track_states (
--        state text primary key
-- );

-- insert into track_states(state) values('new');
-- insert into track_states(state) values('ready');
-- insert into track_states(state) values('queued');
-- insert into track_states(state) values('playing');

create table tracks(
       track_id text primary key,
       station_id text not null,
       created_at timestamptz not null default now(),
       artist text not null,
       title text not null,
       raw_metadata jsonb not null,
       rotation integer not null,
       queues integer not null default 0,
       plays integer not null default 0,
       skips integer not null default 0,
       playing bool not null default false
);
---- create above / drop below ----
drop table tracks;
--drop table track_states;
