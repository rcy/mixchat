create role appuser with password 'appuser' login;
grant appuser to postgres;
create database mixchat_development with owner appuser;
