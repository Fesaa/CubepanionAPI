# Stats service

Provides information about the player count of CubeCraft, and games currently in the lobby. Data is read by the [Laby addon](https://github.com/Fesaa/Cubepanion).

This service is served under `/cubepanion/stats`


## Endpoints

### GET /game/:game

Returns the latest `GameStat` object related to the pasted game. 

### GET /games

Returns an array of the latest `GameStat` objects for each game.


## Models

#### GameStat

```json
{
    "game": "",
    "player_count": 0,
    "unix_time_stamp": 0
}
```

## Notes

The `GameStat` with `game` equal to `Main Lobby`, does not represent the total amount of players in all lobbies, but rather all players on CubeCraft. It is read from the scoreboard.
Other player counts are read from the game compass.
