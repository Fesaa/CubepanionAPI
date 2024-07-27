package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/models"
	"github.com/Fesaa/CubepanionAPI/maps-service/database"
)

var allMaps = []models.GameMap{
	{
		UniqueName: "castle",
		MapName:    "Castle",
		TeamSize:   2,
		Layout:     models.SQUARE,
		Colours:    "red,blue",
		BuildLimit: 40,
		Generators: []models.Generator{},
	},
}

type MockDatabase struct{}

var ms core.MicroService[core.MicroServiceConfig, database.Database]

func mock(config core.DatabaseConfig) (database.Database, error) {
	return &MockDatabase{}, nil
}

func (m *MockDatabase) GetAllMaps() ([]models.GameMap, error) {
	return allMaps, nil
}

func init() {
	var err error
	ms, err = core.NewMicroService(core.LoadDefaultConfigFromEnv(), mock)
	if err != nil {
		panic(err)
	}

	ms.Get("/", Maps)
}

func TestMaps(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	resp, _ := ms.App().Test(req, -1)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
		t.FailNow()
	}

	body, _ := io.ReadAll(resp.Body)

	var got []models.GameMap
	err := json.Unmarshal(body, &got)
	if err != nil {
		t.Errorf("Error unmarshalling response: %s", err)
		t.FailNow()
	}

	if len(got) != 1 {
		t.Errorf("Expected 1 map, got %d", len(got))
		t.FailNow()
	}

	if got[0].UniqueName != "castle" {
		t.Errorf("Expected map name to be castle, got %s", got[0].UniqueName)
		t.FailNow()
	}
}
