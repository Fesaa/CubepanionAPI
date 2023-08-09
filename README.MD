<div align="right">
  <img src="assets/cube_logo.png" alt="Logo" width="100">
</div>

# <center>Cubepanion API</center>


Created as backend to make [Cubepanion](https://github.com/Fesaa/Cubepanion) more dynamically updatable; but publicly available for useless information about CubeCraft Java. 

# Usage

The API is freely available on `http://ameliah.art:7070/`, with a high rate limit. 

## Key features
- Leaderboard API
  - Query a player or games
  - Returning an array of LeaderboardRows
  ```json
  {
    "game": "",
    "player": "",
    "position": 0,
    "score": 0,
    "unix_time_stamp": 0
  }
  ```
- Chest API
  - Query locations for the lobby chests
  - Returning an array of ChestLocations
  ```json
  {
    "season_name": "",
    "x": 0,
    "y": 0,
    "z": 0
  }
  ```
  - Query seasons (lobbies in events have different coords)
  - Returning an array of strings
- EggWarsMap API
  - Query a specific or all eggwars maps
  - Returning (an array of) EggWarsMap
  ```json
  {
    "unique_name": "",
    "map_name": "",
    "team_size": 0,
    "build_limit": 0,
    "colours": "",
    "layout": "",
    "generators": [
      {
        "unique_name": "",
        "ordering": 0,
        "gen_type": "",
        "gen_location": "",
        "level": 0,
        "count": 0
      }
    ]
  }
  ```

## Weave
A java library is available, with no internal thrown errors 🙌, with [docs](https://fesaa.github.io/CubepanionAPI) and [example](https://github.com/Fesaa/Cubepanion) usage.

Other languages may be made if I I'm bored enough for.

# Contribution & Bugs
Always welcome to open an issue or fix the bug yourself and opening a pull request. For features, please first open a ticket or contact me on discord (ameliahh)

<div align="center">
  <img src="assets/background.png" alt="background">
</div>