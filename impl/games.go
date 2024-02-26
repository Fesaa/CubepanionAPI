package impl

import (
	"sync"

	"github.com/Fesaa/CubepanionAPI/models"
)

type GamesImpl struct {
	games     map[string]*models.Game
	gamesLock sync.RWMutex

	aliases     map[string]string
	aliasesLock sync.RWMutex
}

func newGamesImpl(d models.DatabaseProvider) (models.GamesProvider, error) {
	games := make(map[string]*models.Game)
	aliases := make(map[string]string)

	allGames, err := d.GetGames(false)
	if err != nil {
		return nil, err
	}

	for _, game := range allGames {
		games[game.DisplayName] = &game
		for _, alias := range game.Aliases {
			aliases[alias] = game.DisplayName
		}
	}

	return &GamesImpl{games: games, aliases: aliases}, nil
}

func (g *GamesImpl) GetGame(displayName string) *models.Game {
	g.gamesLock.RLock()
	defer g.gamesLock.RUnlock()
	game, ok := g.games[displayName]
	if !ok {
		return nil
	}
	return game
}

func (g *GamesImpl) GetGameDisplayName(s string) string {
	g.aliasesLock.RLock()
	defer g.aliasesLock.RUnlock()
	d, ok := g.aliases[s]
	if !ok {
		return ""
	}
	return d
}
