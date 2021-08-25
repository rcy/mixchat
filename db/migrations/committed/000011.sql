--! Previous: sha1:3575d240fbf0dd7c9406e6b2c13b01a63ffa5ae5
--! Hash: sha1:73a6e7bf9f8d6a100afc232ffb2d77ed46920da4

-- add event type to plays

alter table plays drop column if exists action;
alter table plays add column action text;

update plays set action = 'queued' where action is null;

alter table plays alter column action set not null;
