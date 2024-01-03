create sequence guest_username_counter;

create table users(
       user_id text primary key,
       created_at timestamptz not null default now(),
       username text not null unique default 'FirstTimeCaller' || nextval('guest_username_counter'),
       guest bool not null default true
);

create table sessions(
       session_id text primary key,
       created_at timestamptz not null default now(),
       expires_at timestamptz not null default now() + interval '365 days',
       user_id text references users not null
);

---- create above / drop below ----
drop table sessions;
drop table users;
drop sequence guest_username_counter;
