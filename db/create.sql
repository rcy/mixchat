create role djfm with password 'djfm' login;
create database djfullmoon_development with owner djfm;
create database djfullmoon_development_shadow with owner djfm;
-- alter database djfullmoon_development owner to djfm;
-- alter database djfullmoon_development_shadow owner to djfm;

