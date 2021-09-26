--! Previous: sha1:a98242df51f2930db72ffe94984c7097503a27af
--! Hash: sha1:ee89cfe256c8e066146d095a03ad703cbffa8949

-- add timestamp to station

alter table stations drop column if exists created_at;
alter table stations add column created_at timestamptz default now();

update stations set created_at = now();

alter table stations alter column created_at set not null;
