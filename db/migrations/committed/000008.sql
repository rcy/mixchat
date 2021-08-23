--! Previous: sha1:b16903159624c61cd24f81ebe814905ea13a6c1c
--! Hash: sha1:9a9e7a0fe44b4573acf01d116b29aac9b7b49a3e

alter table tracks drop constraint if exists tracks_unique_filename;
alter table tracks add constraint tracks_unique_filename unique (filename);
