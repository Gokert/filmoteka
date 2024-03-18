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
    title           TEXT   NOT NULL DEFAULT '',
    info            TEXT   NOT NULL DEFAULT '',
    release_date    DATE NOT NULL DEFAULT CURRENT_DATE,
    rating          FLOAT NOT NULL DEFAULT 0,
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

DROP TABLE IF EXISTS profile_role;
CREATE TABLE IF NOT EXISTS profile_role(
    id_profile SERIAL NOT NULL REFERENCES profile(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
    id_role SERIAL NOT NULL REFERENCES role(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    PRIMARY KEY(id_profile, id_role)
    );


DROP TABLE IF EXISTS profile;
CREATE TABLE IF NOT EXISTS profile (
               id SERIAL NOT NULL PRIMARY KEY,
               login TEXT NOT NULL UNIQUE DEFAULT '',
               id_password SERIAL NOT NULL REFERENCES password(id)
                ON DELETE CASCADE
                ON UPDATE CASCADE,
);

DROP TABLE IF EXISTS role;
CREATE TABLE IF NOT EXISTS role (
               id SERIAL NOT NULL PRIMARY KEY,
               value TEXT NOT NULL DEFAULT '',
);



DROP TABLE IF EXISTS password;
CREATE TABLE IF NOT EXISTS password (
            id SERIAL NOT NULL PRIMARY KEY,
            value TEXT NOT NULL DEFAULT '',
);