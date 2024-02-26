package models

type GamesProvider interface {
	GetGame(game string) *Game
	GetGameDisplayName(s string) string
}
