package main

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/models"
	"github.com/Fesaa/CubepanionAPI/games-service/database"
)

type MockDatbase struct{}

var mockGames = []models.Game{
	{
		Game:        "team_eggwars",
		DisplayName: "Team EggWars",
		Aliases:     []string{"team_eggwars", "tew"},
		Active:      true,
		ScoreType:   "win",
	},
	{
		Game:        "team_skywars",
		DisplayName: "Team SkyWars",
		Aliases:     []string{"team_skywars", "tsw"},
		Active:      true,
		ScoreType:   "win",
	},
	{
		Game:        "lucky_islands",
		DisplayName: "Lucky Islands",
		Aliases:     []string{"lucky_islands", "li"},
		Active:      true,
		ScoreType:   "win",
	},
}

func (m *MockDatbase) GetGames(active bool) ([]models.Game, error) {
	return mockGames, nil
}

func (m *MockDatbase) GetGame(s string) (string, error) {
	for _, game := range mockGames {
		if game.Game == s {
			return game.DisplayName, nil
		}
	}

	return "", errors.New("Game not found")
}

func mock(d core.DatabaseConfig) (database.Database, error) {
	return &MockDatbase{}, nil
}

var ms core.MicroService[core.MicroServiceConfig, database.Database]

func init() {
	var err error
	ms, err = core.NewMicroService(core.LoadDefaultConfigFromEnv(), mock)

	if err != nil {
		panic(err)
	}

	ms.Get("/:active", games)
	ms.Get("/game/:game", game)
}

func TestGames(t *testing.T) {
	req, _ := http.NewRequest("GET", "/true", nil)

	resp, _ := ms.App().Test(req, -1)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	got := string(body)

	expected := `[{"game":"team_eggwars","display_name":"Team EggWars","aliases":["team_eggwars","tew"],"active":true,"score_type":"win"},{"game":"team_skywars","display_name":"Team SkyWars","aliases":["team_skywars","tsw"],"active":true,"score_type":"win"},{"game":"lucky_islands","display_name":"Lucky Islands","aliases":["lucky_islands","li"],"active":true,"score_type":"win"}]`
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}
}

func TestGame(t *testing.T) {
	req, _ := http.NewRequest("GET", "/game/team_eggwars", nil)

	resp, _ := ms.App().Test(req, -1)

	body, _ := io.ReadAll(resp.Body)
	got := string(body)

	expected := "Team EggWars"
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}
}
