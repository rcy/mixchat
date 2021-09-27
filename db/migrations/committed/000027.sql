--! Previous: sha1:a4bc7feda91572c8e8a384399bea7480159064dd
--! Hash: sha1:ccef02b48ac1babfff4e5ec505b850120e9a0cb9

-- Enter migration here

alter table tracks
      alter column bucket
      drop default;

drop function if exists bucket();

create or replace function current_bucket(_station_id integer) returns integer as $$
  select coalesce(( select bucket from tracks where station_id = _station_id order by bucket asc limit 1), 0);
$$ language sql;
