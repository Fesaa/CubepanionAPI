# Leaderboard service

Provides information about leaderboards on the CubeCraft Games network. Data is read by the [Laby addon](https://github.com/Fesaa/Cubepanion).

This service is served under `/cubepanion/leaderboard`

## Endpoints

### GET /player/:name

Returns an array `LeaderboardRow`s for the specified player. Each row in game.

### GET /game/:game

Returns an array `LeaderboardRow`s for the specified game. The array will be 200 long.

### GET /game/:game/bounded

Returns an array `LeaderboardRow`s for the specified game. You must provide a `lower` and `upper` query parameter, the array will be as long as their difference.
The endpoint above is generally faster by ~15%, just use whatever is easier for your use case.

### GET /batch

Body must be `BatchRequest`

Returns an array of `LeaderboardRow`s for the passed players, in a specific game. Not guaranteed to be ordered in any specific way

### GET /players

Returns an array of strings, containing all names of the players on the latest leaderboards

## Models

### LeaderboardRow

```json
{
  "game": "",
  "player": "",
  "position": 0,
  "score": 0,
  "texture": "",
  "unix_time_stamp": 0
}
```

The score type should be loaded from the games-service. Load these once, on application startup and then periodically if needed. Don't make a new request each time. They hardly update.

### BatchRequest

```json
{
  "game": "",
  "players": [""]
}
```

## Notes

### Texture field

The texture field in `LeaderboardRow` is a base64 encoding of a (partial) `UUID to Profile and Skin/Cape` response from the Mojang API. 
You may find the JSON structure on the [Protocol Wiki](https://wiki.vg/Mojang_API#UUID_to_Profile_and_Skin.2FCape). You'll most likely be  interested in the `textures.SKIN.url` field, keep in mind that this is an url for the complete skin.
But you may find the players UUID in this url as well.

This is a new field, not all responses will include this just yet.


