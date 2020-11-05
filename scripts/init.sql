DROP TRIGGER IF EXISTS seasons_inc on seasons;
DROP TRIGGER IF EXISTS episodes_inc on episodes;
DROP TABLE IF EXISTS
    users, sessions, content, directors, content_director, actors, content_actor,
    genres, content_genre, countries, content_country, movies, tv_shows, seasons,
    episodes, rates, favourites
    CASCADE;


-- Users of the app
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    nickname varchar(64) NOT NULL,
    email varchar(64) UNIQUE NOT NULL,
    password text NOT NULL,
    avatar varchar(64) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS sessions (
    id serial PRIMARY KEY,
    value varchar(64) UNIQUE NOT NULL,
    expires timestamptz NOT NULL,
    user_id int NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- Content
CREATE TABLE IF NOT EXISTS content (
    id serial PRIMARY KEY,
    name varchar(128) NOT NULL,
    original_name varchar(128) NOT NULL,
    description text NOT NULL,
    short_description text NOT NULL,
    year smallint NOT NULL, -- если сериал, то год выхода 1 сезона
    images varchar(128) NOT NULL, -- путь к папке с постерами (/images/witcher), в которой лежит small.png и large.png
    type varchar(16) NOT NULL -- movie, tvshow

    -- TODO трейлер позже решить
);


-- Content directors
CREATE TABLE IF NOT EXISTS directors (
    id serial PRIMARY KEY,
    name varchar(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS content_director (
    content_id int NOT NULL,
    director_id int NOT NULL,

    PRIMARY KEY(content_id, director_id),
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (director_id) REFERENCES directors(id) ON DELETE CASCADE
);


-- Content actors
CREATE TABLE IF NOT EXISTS actors (
    id serial PRIMARY KEY,
    name varchar(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS content_actor (
    content_id int NOT NULL,
    actor_id int NOT NULL,
    
    PRIMARY KEY(content_id, actor_id),
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES actors(id) ON DELETE CASCADE
);


-- Сontent genres
CREATE TABLE IF NOT EXISTS genres (
    id serial PRIMARY KEY,
    name varchar(64) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS content_genre (
    content_id int NOT NULL,
    genre_id int NOT NULL,

    PRIMARY KEY(content_id, genre_id),
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
);


-- Сontent countries
CREATE TABLE IF NOT EXISTS countries (
    id serial PRIMARY KEY,
    name varchar(64) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS content_country (
    content_id int NOT NULL,
    country_id int NOT NULL,

    PRIMARY KEY(content_id, country_id),
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE
);


-- Movie that has one to one relationship with content
CREATE TABLE IF NOT EXISTS movies (
    id serial PRIMARY KEY,
    video varchar(128) NOT NULL,
    content_id int UNIQUE NOT NULL, -- one to one with content

    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
);


-- TVShow that has one to one relationship with content
CREATE TABLE IF NOT EXISTS tv_shows (
    id serial PRIMARY KEY,
    seasons int NOT NULL DEFAULT 0, -- тригер на каждое создание сезона
    content_id int UNIQUE NOT NULL,

    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
);

-- TVShow seasons
CREATE TABLE IF NOT EXISTS seasons (
    id serial PRIMARY KEY,
    number int NOT NULL,
    episodes int NOT NULL DEFAULT 0, -- тригер на каждое создание эпизода
    tv_show_id int NOT NULL,

    FOREIGN KEY (tv_show_id) REFERENCES tv_shows(id) ON DELETE CASCADE
);

-- TVShow episodes
CREATE TABLE IF NOT EXISTS episodes (
    id serial PRIMARY KEY,
    number int NOT NULL,
    name varchar(128) NOT NULL,
    video varchar(128) NOT NULL,
    description text NOT NULL,
    poster varchar(128) NOT NULL, -- путь к папке с постерами (/images/witcher/s1 /s2 ...), в которой лежит e1.png e2.png ...
    season_id int NOT NULL,

    FOREIGN KEY (season_id) REFERENCES seasons(id) ON DELETE CASCADE
);


-- Users content rating
CREATE TABLE IF NOT EXISTS rates (
    user_id int NOT NULL,
    content_id int NOT NULL,
    likes boolean NOT NULL,

    PRIMARY KEY(user_id, content_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
);


-- Users favourite content
CREATE TABLE IF NOT EXISTS favourites (
    user_id int NOT NULL,
    content_id int NOT NULL,
    created timestamptz NOT NULL,

    PRIMARY KEY(user_id, content_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE
);


-- Triger for increment seasons number in tw_show
CREATE OR REPLACE FUNCTION seasons_inc() RETURNS trigger AS
$seasons_inc$
    BEGIN
        UPDATE tv_shows
        SET seasons = seasons + 1
        WHERE id=NEW.tv_show_id;
        RETURN NEW;
    END;
$seasons_inc$
LANGUAGE plpgsql;

CREATE TRIGGER seasons_inc AFTER INSERT ON seasons
    FOR EACH ROW EXECUTE PROCEDURE seasons_inc();


-- Triger for increment episodes number in seasons
CREATE OR REPLACE FUNCTION episodes_inc() RETURNS trigger AS
$episodes_inc$
    BEGIN
        UPDATE seasons
        SET episodes = episodes + 1
        WHERE id=NEW.season_id;
        RETURN NEW;
    END;
$episodes_inc$
LANGUAGE plpgsql;

CREATE TRIGGER episodes_inc AFTER INSERT ON episodes
    FOR EACH ROW EXECUTE PROCEDURE episodes_inc();
