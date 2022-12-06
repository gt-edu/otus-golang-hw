CREATE TABLE IF NOT EXISTS events(
    id serial primary key,
    owner_id bigint,
    title text,
    descr text,
    start_date date not null,
    start_time time,
    end_date date not null,
    end_time time,
    notification_period text
);

create index owner_idx on events (owner_id);
create index start_idx on events using btree (start_date, start_time);
