alter table stations add column user_id text references users not null;
---- create above / drop below ----
alter table stations drop column user_id;
