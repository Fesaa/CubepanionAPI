# Welcome to the Cubepanion API documentation

Relevant GitHub repositories for the project are the [api](https://github.com/Fesaa/CubepanionAPI), and the [Laby addon](https://github.com/Fesaa/Cubepanion).

The project aims to make the addon more dynamic and less dependent on updates to the addon itself. With the extra benefit of being able to use the API for other projects.

## Services

- `chests` - Provides information about the Lobby Chest Hunter locations
- `games` - Provides information about the avaible games
- `leaderboards` - Provides information about game, and player leaderboards
- `maps` - Provides information about EggWars maps

## Project layout & deployment

Each service is deployed as a seperate docker container, and served behind a nginx reverse proxy. There is a redis cache layer present, so you may so delayed updates to the API.
