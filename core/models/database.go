package models

const LEADERBOARD_SIZE int = 200

type Season string

type ChestLocation struct {
	Season Season `json:"season_name"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Z      int    `json:"z"`
}

type MapLayout string

const (
	DOUBLE_CROSS MapLayout = "double_cross"
	TRIAGLE      MapLayout = "triangle"
	SQUARE       MapLayout = "square"
	CROSS        MapLayout = "cross"
)

type EggWarsMap struct {
	UniqueName string      `json:"unique_name"`
	MapName    string      `json:"map_name"`
	TeamSize   int         `json:"team_size"`
	Layout     MapLayout   `json:"layout"`
	Colours    string      `json:"colours"`
	BuildLimit int         `json:"build_limit"`
	Generators []Generator `json:"generators"`
}

type GenType string
type GenLocation string

const (
	DIAMOND GenType = "diamond"
	GOLD    GenType = "gold"
	IRON    GenType = "iron"

	MIDDLE     GenLocation = "middle"
	SEMIMIDDLE GenLocation = "semimiddle"
	BASE       GenLocation = "base"
)

type Generator struct {
	UniqueName string      `json:"unique_name"`
	Ordering   int         `json:"ordering"`
	Type       GenType     `json:"gen_type"`
	Location   GenLocation `json:"gen_location"`
	Level      int         `json:"level"`
	Count      int         `json:"count"`
}

type Game struct {
	Game        string   `json:"game"`
	DisplayName string   `json:"display_name"`
	Aliases     []string `json:"aliases"`
	Active      bool     `json:"active"`
	ScoreType   string   `json:"score_type"`
}

type LeaderboardRow struct {
	Game          string `json:"game" validate:"nonzero,regexp=[a-zA-Z0-9_ ]"`
	Player        string `json:"player" validate:"nonzero,regexp=[a-zA-Z0-9_]{3\\,16}"`
	Position      int    `json:"position" validate:"nonzero,min=1,max=200"`
	Score         int    `json:"score" validate:"nonzero"`
	UnixTimeStamp int    `json:"unix_time_stamp"`
}

type Submission struct {
	Uuid          string `json:"uuid"`
	UnixTimeStamp int    `json:"unix_time_stamp"`
	Game          string `json:"game"`
	Valid         bool   `json:"valid"`
}

type Location struct {
	Current    string `json:"current"`
	Previous   string `json:"previous"`
	InPreLobby bool   `json:"inPreLobby"`
}
