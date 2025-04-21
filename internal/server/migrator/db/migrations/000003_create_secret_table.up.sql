CREATE TABLE secret_info (
   id serial primary key,
   identifier varchar not null,
   value varchar not null,
   client_user integer NOT NULL REFERENCES users(id),
   createddate TIMESTAMP default now()
);

COMMIT;

CREATE INDEX secret_info__cu__identifier__index
ON secret_info (client_user,identifier);
