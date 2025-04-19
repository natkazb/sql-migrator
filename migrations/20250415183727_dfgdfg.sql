-- SQL
-- Up begin
create table if not exists event (
    id serial PRIMARY KEY,
    title text NOT NULL,
    start_date timestamptz NOT NULL,
    end_date timestamptz NOT NULL,
    description text NOT NULL,
    user_id integer NOT NULL,
    notify_on integer NOT NULL
);
-- Up end

-- Down begin
drop table if exists event;
-- Down end