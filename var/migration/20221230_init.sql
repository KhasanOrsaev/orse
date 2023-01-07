create table if not exists devices(
    id integer primary key,
    name text,
    address text,
    status bool,
    active bool
);