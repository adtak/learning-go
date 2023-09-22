CREATE TABLE path_url (
    id serial primary key,
    from_path text not null,
    to_url text not null,
    created timestamp with time zone not null default CURRENT_TIMESTAMP,
    updated timestamp with time zone not null default CURRENT_TIMESTAMP
);

INSERT INTO path_url (from_path, to_url) VALUES ('/search', 'https://www.google.co.jp');