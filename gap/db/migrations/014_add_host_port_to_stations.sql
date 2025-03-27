alter table stations add column host_port text not null default '';

---- create above / drop below ----

alter table stations drop column host_port;
