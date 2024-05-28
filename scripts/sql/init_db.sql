DROP TABLE IF EXISTS actor CASCADE;
CREATE TABLE IF NOT EXISTS actor (
                                     id          SERIAL NOT NULL PRIMARY KEY,
                                     name        TEXT NOT NULL DEFAULT '',
                                     gen         TEXT NOT NULL DEFAULT '',
                                     birthdate   DATE NOT NULL DEFAULT CURRENT_DATE
);

DROP TABLE IF EXISTS film CASCADE;
CREATE TABLE IF NOT EXISTS film (
                                    id              SERIAL NOT NULL PRIMARY KEY,
                                    title           TEXT   NOT NULL DEFAULT '',
                                    info            TEXT   NOT NULL DEFAULT '',
                                    release_date    DATE NOT NULL DEFAULT CURRENT_DATE,
                                    rating          FLOAT NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS actor_in_film CASCADE;
CREATE TABLE IF NOT EXISTS actor_in_film(
                                            id_film SERIAL NOT NULL REFERENCES film(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
    id_actor SERIAL NOT NULL REFERENCES actor(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    PRIMARY KEY(id_actor, id_film)
    );

DROP TABLE IF EXISTS profile CASCADE;
CREATE TABLE IF NOT EXISTS profile (
                                       id SERIAL NOT NULL PRIMARY KEY,
                                       login TEXT NOT NULL UNIQUE DEFAULT '',
                                       password bytea NOT NULL DEFAULT '',
                                       role TEXT NOT NULL DEFAULT 'user'
);

INSERT INTO profile(login, password, role) VALUES ('admin', '\xc7ad44cbad762a5da0a452f9e854fdc1e0e7a52a38015f23f3eab1d80b931dd472634dfac71cd34ebc35d16ab7fb8a90c81f975113d6c7538dc69dd8de9077ec', 'admin');