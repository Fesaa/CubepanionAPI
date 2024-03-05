package impl

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Fesaa/CubepanionAPI/models"
)

type GamesImpl struct {
	games     map[string]*models.Game
	gamesLock *sync.RWMutex

	aliases     map[string]string
	aliasesLock *sync.RWMutex
}

func newGamesImpl(d models.DatabaseProvider) (models.GamesProvider, error) {
	games := make(map[string]*models.Game)
	aliases := make(map[string]string)

	g := &GamesImpl{games: games, aliases: aliases, gamesLock: &sync.RWMutex{}, aliasesLock: &sync.RWMutex{}}
	go g.loadGames(d)

	return g, nil
}

func (g *GamesImpl) loadGames(d models.DatabaseProvider) {
	for range time.Tick(time.Duration(5) * time.Minute) {
		allGames, err := d.GetGames(false)
		if err != nil {
			slog.Error(fmt.Sprintf("Error loading games: %s", err.Error()))
			continue
		}
		g.gamesLock.Lock()
		g.aliasesLock.Lock()
		g.games = make(map[string]*models.Game)
		g.aliases = make(map[string]string)
		for _, game := range allGames {
			g.games[game.DisplayName] = &game
			g.aliases[game.DisplayName] = game.DisplayName
			g.aliases[game.Game] = game.DisplayName
			for _, alias := range game.Aliases {
				slog.Debug(fmt.Sprintf("Adding alias %s for game %s", alias, game.DisplayName))
				g.aliases[alias] = game.DisplayName
			}
		}
		g.aliasesLock.Unlock()
		g.gamesLock.Unlock()
	}
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
