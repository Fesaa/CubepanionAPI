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
