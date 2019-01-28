CREATE TABLE stones (
    id SERIAL PRIMARY KEY,
    game_id int REFERENCES games NOT NULL,
    x int NOT NULL,
    y int NOT NULL,
    color character varying(255) NOT NULL,
    inserted_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    UNIQUE (game_id, x, y)
);
