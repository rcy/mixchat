create table searches(
       search_id text primary key,
       station_id text not null,
       created_at timestamptz not null default now(),
       query text not null,
       status text not null default 'pending'
);
create table results(
       result_id text primary key,
       search_id text not null,
       station_id text not null,
       created_at timestamptz not null default now(),
       extern_id text not null,
       url text not null,
       thumbnail text,
       title text,
       duration integer,
       views bigint
);

---- create above / drop below ----
drop table results;
drop table searches;
