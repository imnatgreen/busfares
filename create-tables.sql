create extension if not exists "uuid-ossp";
drop table if exists fares;
create table fares (
  id uuid default uuid_generate_v4() primary key,
  fare jsonb not null
);
create index operators on fares using gin (data->)