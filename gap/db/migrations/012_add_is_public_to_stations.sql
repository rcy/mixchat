alter table stations add column is_public bool not null default false;

---- create above / drop below ----

alter table stations drop column is_public;
