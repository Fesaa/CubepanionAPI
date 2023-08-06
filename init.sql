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
)

CREATE TABLE IF NOT EXISTS chest_locations (
    season_name VARCHAR NOT NULL,
    x INT NOT NULL,
    y INT NOT NULL,
    z INT NOT NULL,
    FOREIGN KEY (season_name)
        REFERENCES seasons(season_name)
);