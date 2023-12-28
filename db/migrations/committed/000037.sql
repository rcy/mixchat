--! Previous: sha1:4668e7748e3da2896ec49b4ef4bab93b6d7dc566
--! Hash: sha1:aa809886c32285dbaf8fa1d3a980a334c25f59e2

alter table stations drop column if exists active;
alter table stations add column active bool default false;
