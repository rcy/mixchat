alter table stations add column background_image_url text not null default '';
---- create above / drop below ----
alter table stations drop column background_image_url;
