CREATE TYPE game_type AS ENUM ('STANDARD');

CREATE TYPE game_state AS ENUM ('INVITATION', 'IN_PROGRESS', 'FINISHED' );

CREATE TABLE games (
    id serial PRIMARY KEY,
    type game_type NOT NULL,
    state game_state NOT NULL,
    board_size int NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
