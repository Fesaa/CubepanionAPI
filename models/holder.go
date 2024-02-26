package models

const HOLDER_KEY string = "holder"

type Holder interface {
	GetDatabaseProvider() DatabaseProvider
	GetGamesProvider() GamesProvider
}
