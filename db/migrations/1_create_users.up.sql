CREATE TABLE identities (
  id serial PRIMARY KEY,
  email text NOT NULL UNIQUE,
  password_hash text NOT NULL,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL
);

CREATE TABLE users (
    id serial PRIMARY KEY,
    identity_id integer REFERENCES identities(id) UNIQUE,
    name text NOT NULL UNIQUE,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
