CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    status character varying(255) NOT NULL,
    player_turn_id character varying(255),
    board_size int,
    last_taker json,
    inserted_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
