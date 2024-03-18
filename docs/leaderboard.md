# Leaderboard service

Provides information about leaderboards on the CubeCraft Games network. Data is read by the [Laby addon](https://github.com/Fesaa/Cubepanion).

This service is served under `/cubepanion/leaderboard`

## Endpoints

### GET /player/:name

Returns an array `LeaderboardRow`s for the specified player. Each row in game.

### GET /game/:game

Returns an array `LeaderboardRow`s for the specified game. The array will be 200 long.

### GET /game/:game/bounded

Returns an array `LeaderboardRow`s for the specified game. The array will be 200 long. You must provide a `start` and `end` query parameter.

## Models

### LeaderboardRow

```json
{
  "game": "",
  "player": "",
  "position": 0,
  "score": 0,
  "unix_time_stamp": 0
}
```

The score type should be loaded from the games-service. Load these once, on application startup and then periodically if needed. Don't make a new request each time. They hardly update.
