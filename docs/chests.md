# Chests Service

Provides locations of the Lobby Chest Hunter locations. And seasons, these seasons are not official and only used for internal seperation.

This service is served under `/cubepanion/chests`

## Endpoints

### GET /

Returns an array of `ChestLocation`s for the seasons currently active.

### GET /:season

Returns an array of `ChestLocation`s for the specified season.

### GET /seasons/:active

Returns an array of (active) `Seaon`s.

## Models

### ChestLocation

```json
{
  "seasons_name": "",
  "x": 0,
  "y": 0,
  "z": 0
}
```

### Season

Simple string, representing the season name.
