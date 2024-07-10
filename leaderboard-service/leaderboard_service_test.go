package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/models"
	"github.com/Fesaa/CubepanionAPI/leaderboard-service/database"
)

type MockDatbase struct{}

var leaderboards = []models.LeaderboardRow{
	{
		Game:          "Team EggWars",
		Player:        "Mivke",
		Position:      1,
		Score:         100,
		UnixTimeStamp: -1,
	},
	{
		Game:          "Team EggWars",
		Player:        "Fesa",
		Position:      2,
		Score:         50,
		UnixTimeStamp: -1,
	},
}

func mock(d core.DatabaseConfig) (database.Database, error) {
	return &MockDatbase{}, nil
}

func (m *MockDatbase) GetLeaderboard(game string) ([]models.LeaderboardRow, error) {
	rows := []models.LeaderboardRow{}
	for _, row := range leaderboards {
		if row.Game == game {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func (m *MockDatbase) GetLeaderboardBounded(game string, start, end int) ([]models.LeaderboardRow, error) {
	rows := []models.LeaderboardRow{}
	for _, row := range leaderboards {
		if row.Game == game && row.Position >= start && row.Position <= end {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func (m *MockDatbase) GetLeaderboardForPlayer(player string) ([]models.LeaderboardRow, error) {
	rows := []models.LeaderboardRow{}
	for _, row := range leaderboards {
		if row.Player == player {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func (m *MockDatbase) SubmitLeaderboard(req models.LeaderboardSubmission) error {
	return nil
}

func (m *MockDatbase) GetLeaderboardForPlayers(req models.GamePlayersRequest) ([]models.LeaderboardRow, error) {
	rows := []models.LeaderboardRow{}
	for _, row := range leaderboards {
		for _, player := range req.Players {
			if row.Player == player && row.Game == req.Game {
				rows = append(rows, row)
			}
		}
	}
	return rows, nil
}

func (m *MockDatbase) GetAllPlayers() ([]string, error) {
	players := []string{}
	for _, row := range leaderboards {
		players = append(players, row.Player)
	}
	return players, nil
}

func testConfig() LeaderboardServiceConfig {
	c := core.LoadDefaultConfigFromEnv()
	return LeaderboardServiceConfig{
		YHost:            c.Host(),
		YPort:            c.Port(),
		YGamesServiceURL: os.Getenv("GAMES_SERVICE_URL"),
		YDatabase:        *c.Database().(*core.DefaultDatabaseConfig),
		YRedis:           c.Redis().(core.DefaultRedisConfig),
		YLogging:         *c.LoggingConfig().(*core.DefaultLoggingConfig),
	}
}

var ms core.MicroService[LeaderboardServiceConfig, database.Database]

func init() {
	var err error
	ms, err = core.NewMicroService(testConfig(), mock)

	if err != nil {
		panic(err)
	}

	ms.Post("/", Submit)
	ms.Get("/player/:name", PlayerLeaderboard)
	ms.Get("/game/:game", GameLeaderboard)
	ms.Get("/game/:game/bounded", GameLeaderboardBounded)
}

func TestSubmittion(t *testing.T) {
	body := models.LeaderboardSubmission{
		Uuid:          "",
		Game:          "Team EggWars",
		UnixTimeStamp: 1,
		Entries:       []models.LeaderboardRow{},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	req := httptest.NewRequest("POST", "/", &buf)

	resp, _ := ms.App().Test(req, -1)
	if resp.StatusCode != 400 {
		t.Errorf("Expected 400, got %d", resp.StatusCode)
		t.Fail()
	}
}

func TestPlayerLeaderboard(t *testing.T) {
	req := httptest.NewRequest("GET", "/player/Mivke", nil)

	resp, _ := ms.App().Test(req, -1)
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
		t.Fail()
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var rows []models.LeaderboardRow
	err := json.Unmarshal(body, &rows)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(rows) != 1 {
		t.Errorf("Expected 1, got %d", len(rows))
		t.Fail()
	}
}
