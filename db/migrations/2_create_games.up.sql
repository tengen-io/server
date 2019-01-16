CREATE TABLE games (
    id bigint PRIMARY KEY,
    status character varying(255) NOT NULL,
    player_turn_id character varying(255),
    board json,
    inserted_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE SEQUENCE games_id_seq START 1;
