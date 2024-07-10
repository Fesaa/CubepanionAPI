# Games Service

Provides information about (active) games on the CubeCraft Games network.

This service is served under `/cubepanion/games`

## Endpoints

### GET /:active

Returns an array of (active) `Game`s.

### GET /game/:game

Returns a `Game` object for the specified game. This may be the display, game, or any of the aliases. 

## Models

### Game

```json
{
  "game": "",
  "display_name": "",
  "aliases": [""],
  "active": true,
  "score_type": ""
}
```
