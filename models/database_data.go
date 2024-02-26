package models

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
	Game          string `json:"game"`
	Player        string `json:"player"`
	Position      int    `json:"position"`
	Score         int    `json:"score"`
	UnixTimeStamp int    `json:"unix_time_stamp"`
}
