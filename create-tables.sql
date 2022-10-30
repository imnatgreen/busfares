create extension if not exists "uuid-ossp";
drop table if exists fares;
create table fares (
  id uuid default uuid_generate_v4() primary key,
  fare_object jsonb not null
);
drop index if exists scheduled_stop_points_index;
drop index if exists lines_index;
create index if not exists lines_index on fares using gin ((fare_object->'Lines') jsonb_path_ops);
create index if not exists scheduled_stop_points_index on fares using gin (
  (fare_object->'ScheduledStopPoints') jsonb_path_ops
);