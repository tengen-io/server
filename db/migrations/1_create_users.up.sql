CREATE TABLE identities (
  id serial PRIMARY KEY,
  email character varying(255) NOT NULL UNIQUE,
  password_hash text NOT NULL,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL
);

CREATE TABLE users (
    id serial PRIMARY KEY,
    identity_id integer REFERENCES identities(id) UNIQUE,
    name character varying(255) NOT NULL UNIQUE,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
