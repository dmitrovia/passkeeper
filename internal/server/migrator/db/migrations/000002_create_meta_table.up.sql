CREATE TABLE meta (
   id serial primary key,
   file_name varchar not null unique,
   hash_md varchar not null,
   index_number int,
   client_user integer NOT NULL REFERENCES users(id),
   createddate TIMESTAMP default now()
);

COMMIT;

CREATE INDEX meta__client_user__index
ON meta (client_user);
