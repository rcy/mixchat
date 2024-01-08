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
       thumbnail text not null,
       title text not null,
       uploader text not null,
       duration float not null,
       views float not null
);

---- create above / drop below ----
drop table results;
drop table searches;
