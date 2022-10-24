create extension if not exists "uuid-ossp";
-- drop table if exists fares;
create table fares (
  id uuid default uuid_generate_v4() primary key,
  fare_object jsonb not null
);
-- create index if not exists fare_index on fares using gin (fare_object);
-- drop index if exists fare_index;
-- drop index if exists lines_index;
create index if not exists lines_index on fares using gin ((fare_object->'Lines') jsonb_path_ops);
create index if not exists scheduled_stop_points_index on fares using gin (
  (fare_object->'ScheduledStopPoints') jsonb_path_ops
);
-- explain analyze
-- select *
-- from fares
-- where fare_object->'Lines' @> '[{"PublicCode": "483", "OperatorRef": {"Ref": "noc:ROST"}}]'
--   and fare_object->'ScheduledStopPoints' @> '[{"ScheduledStopPointRef": "atco:25001425"}, {"ScheduledStopPointRef": "atco:2500IMG2914"}]';