DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS profile;

CREATE TABLE IF NOT EXISTS profile (
    id serial PRIMARY KEY,
    nickname text NOT NULL,
    email text UNIQUE NOT NULL,
    password text NOT NULL,
    avatar text NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS session (
    id serial PRIMARY KEY,
    value text UNIQUE NOT NULL,
    expires timestamptz NOT NULL,
    profile_id int NOT NULL,

    FOREIGN KEY (profile_id) REFERENCES profile(id)
);
