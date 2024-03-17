DROP TABLE IF EXISTS actor;
CREATE TABLE IF NOT EXISTS actor (
    id          SERIAL NOT NULL PRIMARY KEY,
    name        TEXT NOT NULL,
    gen         TEXT NOT NULL,
    birthdate   DATE NOT NULL
);

DROP TABLE IF EXISTS film;
CREATE TABLE IF NOT EXISTS film (
    id              SERIAL NOT NULL PRIMARY KEY,
    title           TEXT   NOT NULL,
    info            TEXT   NOT NULL,
    release_date    DATE NOT NULL,
    rating          FLOAT NOT NULL
);

DROP TABLE IF EXISTS actor_in_film;
CREATE TABLE IF NOT EXISTS actor_in_film(
    id_film SERIAL NOT NULL REFERENCES film(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
    id_actor SERIAL NOT NULL REFERENCES actor(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    PRIMARY KEY(id_actor, id_film)
);

DROP TABLE IF EXISTS profile;
CREATE TABLE IF NOT EXISTS profile (
   id SERIAL NOT NULL PRIMARY KEY,
   login TEXT NOT NULL UNIQUE DEFAULT '',
   password TEXT NOT NULL DEFAULT '',
   role TEXT NOT NULL DEFAULT 'user'
);