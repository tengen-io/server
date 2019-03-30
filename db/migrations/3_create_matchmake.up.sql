CREATE TABLE matchmake_requests (
    id bigserial PRIMARY KEY,
    queue text NOT NULL,
    user_id integer REFERENCES users(id) NOT NULL,
    rank integer NOT NULL,
    rank_delta integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    UNIQUE (user_id)
);
