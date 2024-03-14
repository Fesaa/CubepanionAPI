package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Fesaa/CubepanionAPI/chests-service/database"
	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/models"
)

var season string

type MockDatabase struct{}

var allChests = []models.ChestLocation{
	{
		Season: "winter_2024",
		X:      0,
		Y:      0,
		Z:      0,
	},
	{
		Season: "summer_2024",
		X:      0,
		Y:      0,
		Z:      0,
	},
}

func (m *MockDatabase) GetSeasons(active bool) ([]models.Season, error) {
	return []models.Season{"winter_2024", "summer_2024"}, nil
}

func (m *MockDatabase) GetChests(season string) ([]models.ChestLocation, error) {
	chest := make([]models.ChestLocation, 0)
	s := models.Season(season)
	for _, c := range allChests {
		if c.Season == s {
			chest = append(chest, c)
		}
	}

	return chest, nil
}

func (m *MockDatabase) GetCurrentChests() ([]models.ChestLocation, error) {
	return m.GetChests("winter_2024")
}

var ms core.MicroService[core.MicroServiceConfig, database.Database]

func mock(config core.DatabaseConfig) (database.Database, error) {
	return &MockDatabase{}, nil
}

func init() {
	var err error
	ms, err = core.NewMicroService(core.LoadDefaultConfigFromEnv(), mock)

	if err != nil {
		panic(err)
	}

	ms.Get("/", CurrentChestLocations)
	ms.Get("/:season", ChestLocations)
	ms.Get("/seasons/:active", Seasons)
}

func TestSeasons(t *testing.T) {
	req := httptest.NewRequest("GET", "/seasons/true", nil)

	resp, _ := ms.App().Test(req, -1)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var seasons []models.Season
	json.Unmarshal(body, &seasons)

	if len(seasons) != 2 {
		t.Errorf("Expected 1 season, got %d", len(seasons))
		t.FailNow()
	}

}

func TestCurrentChestLocations(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	resp, _ := ms.App().Test(req, -1)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var chests []models.ChestLocation
	json.Unmarshal(body, &chests)

	if len(chests) != 1 {
		t.Errorf("Expected 1 chest, got %d", len(chests))
		t.FailNow()
	}

	if chests[0].Season != "winter_2024" {
		t.Errorf("Expected winter_2024, got %s", chests[0].Season)
	}
}

func TestChestLocations(t *testing.T) {
	req := httptest.NewRequest("GET", "/summer_2024", nil)

	resp, _ := ms.App().Test(req, -1)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var chests []models.ChestLocation
	json.Unmarshal(body, &chests)

	if len(chests) != 1 {
		t.Errorf("Expected 1 chest, got %d", len(chests))
		t.FailNow()
	}

	if chests[0].Season != "summer_2024" {
		t.Errorf("Expected summer_2024, got %s", chests[0].Season)
	}
}
