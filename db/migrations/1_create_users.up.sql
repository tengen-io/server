CREATE TABLE users (
    id int PRIMARY KEY,
    username character varying(255) NOT NULL UNIQUE,
    email character varying(255) NOT NULL UNIQUE,
    encrypted_password character varying(255) NOT NULL,
    inserted_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE SEQUENCE users_id_seq START 1;
