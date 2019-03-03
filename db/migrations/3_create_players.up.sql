CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    user_id int REFERENCES users NOT NULL,
    game_id int REFERENCES games NOT NULL,
    status character varying(255) NOT NULL,
    color character varying(255) NOT NULL,
    stats json,
    has_passed boolean DEFAULT false,
    prisoners int DEFAULT 0,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
