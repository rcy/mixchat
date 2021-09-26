--! Previous: sha1:6e22a654986ec51a45fad20327537a7f043c27be
--! Hash: sha1:68c0fc8d2e01aae007a7f1f56cc433d824898148

drop table if exists stations;
create table stations (
       id serial primary key,
       slug text not null unique,
       name text
);
