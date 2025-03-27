alter table stations rename column host_port to telnet_port;
alter table stations add column broadcast_port text not null default '';

---- create above / drop below ----

alter table stations drop column broadcast_port;
alter table stations rename column telnet_port to host_port;
