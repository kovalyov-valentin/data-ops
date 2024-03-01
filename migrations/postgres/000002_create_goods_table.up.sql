CREATE TABLE IF NOT EXISTS goods(
    id SERIAL PRIMARY KEY,
    projects_id INTEGER REFERENCES projects(id),
    name VARCHAR(255) NOT NUll,
    description VARCHAR(255) NOT NULL DEFAULT '',
    priority SERIAL,/* max +1*/
    removed BOOLEAN DEFAULT false,
    created_at TIMESTAMP    NOT NULL default now()
);
create index ON goods using btree(name);
insert into projects(name) values ('Первая запись');
