CREATE TYPE game_type AS ENUM ('STANDARD');

CREATE TYPE game_state AS ENUM ('INVITATION', 'IN_PROGRESS', 'FINISHED' );

CREATE TYPE game_user_type AS ENUM ('OWNER');

CREATE TABLE games (
    id serial PRIMARY KEY,
    type game_type NOT NULL,
    state game_state NOT NULL,
    board_size int NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE game_user (
  game_id integer REFERENCES games(id) NOT NULL,
  user_id integer REFERENCES users(id) NOT NULL,
  type game_user_type NOT NULL,
  created_at timestamp without time zone NOT NULL,
  updated_at timestamp without time zone NOT NULL,
  PRIMARY KEY (game_id, user_id)
);

CREATE INDEX game_user_user_game_idx ON game_user (user_id, game_id);
