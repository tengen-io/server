CREATE TABLE players (
    id bigint PRIMARY KEY,
    user_id bigint REFERENCES users,
    game_id bigint REFERENCES games,
    status character varying(255) NOT NULL,
    color character varying(255) NOT NULL,
    stats jsonb,
    has_passed boolean DEFAULT false,
    inserted_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE SEQUENCE players_id_seq START 1;
