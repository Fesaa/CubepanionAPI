CREATE TABLE IF NOT EXISTS submissions (
    uuid VARCHAR NOT NULL,
    unix_time_stamp BIGINT NOT NULL,
    game VARCHAR NOT NULL,
    valid BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (unix_time_stamp)
);

CREATE TABLE IF NOT EXISTS leaderboards (
    game VARCHAR NOT NULL,
    player VARCHAR NOT NULL,
    position INT NOT NULL,
    score INT NOT NULL,
    unix_time_stamp BIGINT NOT NULL,
    FOREIGN KEY (unix_time_stamp)
        REFERENCES submissions(unix_time_stamp)
);

CREATE TABLE IF NOT EXISTS ban (
    uuid VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS seasons (
    season_name VARCHAR NOT NULL,
    running BOOLEAN NOT NULL,
    PRIMARY KEY (season_name)
);

CREATE TABLE IF NOT EXISTS chest_locations (
    season_name VARCHAR NOT NULL,
    x INT NOT NULL,
    y INT NOT NULL,
    z INT NOT NULL,
    FOREIGN KEY (season_name)
        REFERENCES seasons(season_name)
);

CREATE TYPE maptype AS ENUM ('cross', 'doublecross', 'square');
CREATE TYPE genlocation AS ENUM ('base', 'semimiddle', 'middle');
CREATE TYPE gentype AS ENUM ('iron', 'gold', 'diamond');

CREATE TABLE IF NOT EXISTS eggwars_maps (
    unique_name VARCHAR NOT NULL,
    map_name VARCHAR NOT NULL,
    team_size INT NOT NULL,
    build_limit INT NOT NULL,
    colours VARCHAR NOT NULL,
    PRIMARY KEY (unique_name)
);


CREATE TABLE IF NOT EXISTS generators (
    unique_name VARCHAR NOT NULL,
    ordering INT NOT NULL,
    gen_type gentype NOT NULL,
    gen_location genlocation NOT NULL,
    level INT NOT NULL,
    count INT NOT NULL,
    PRIMARY KEY (unique_name, ordering),
    FOREIGN KEY (unique_name)
        REFERENCES eggwars_maps(unique_name)
);

CREATE TABLE IF NOT EXISTS games (
    game VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    aliases VARCHAR(255) NOT NULL DEFAULT '',
    active BOOLEAN NOT NULL,
    score_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (game)
);

CREATE TABLE IF NOT EXISTS player_locations (
    uuid VARCHAR NOT NULL,
    previous VARCHAR NOT NULL,
    current VARCHAR NOT NULL,
    in_pre_lobby BOOLEAN NOT NULL,
    version INT,
    PRIMARY KEY (uuid)
);
