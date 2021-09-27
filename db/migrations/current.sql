-- Enter migration here

create or replace function current_bucket(_station_id integer) returns integer as $$
  select coalesce(( select bucket from tracks where station_id = _station_id order by bucket asc limit 1), 0);
$$ language sql;

alter table tracks
      alter column bucket
      drop default;
