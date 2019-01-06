CREATE TABLE games (
    id bigint NOT NULL,
    status character varying(255) NOT NULL,
    player_turn_id integer,
    inserted_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE SEQUENCE games_id_seq START 1;
