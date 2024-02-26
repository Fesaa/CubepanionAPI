package routes

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/Fesaa/CubepanionAPI/models"
	"github.com/gofiber/fiber/v2"
)

var playerRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)
var gameRegex = regexp.MustCompile(`^[a-zA-Z0-9_ ]`)

func LeaderboardApi(app *fiber.App) {
	g := app.Group("/leaderboard_api")
	g.Get("/leaderboard/:game", leaderboardAPI_game)
	g.Get("/leaderboard/:game/bounded", leaderboardAPI_game_bounded)
	g.Get("/player/:name", leaderboardAPI_player)
	g.Get("/games/:active", leaderboardAPI_games)
}

func leaderboardAPI_games(c *fiber.Ctx) error {
	holder, _ := c.Locals(models.HOLDER_KEY).(models.Holder)
	db := holder.GetDatabaseProvider()

	activeS := c.Query("active", "true")
	active := activeS == "true"

	games, err := db.GetGames(active)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(games)
}

func leaderboardAPI_player(c *fiber.Ctx) error {
	holder, _ := c.Locals(models.HOLDER_KEY).(models.Holder)
	db := holder.GetDatabaseProvider()

	player := c.Params("name")
	if player == "" {
		return jsonError(c, 400, "name parameter is required")
	}
	if !playerRegex.MatchString(player) {
		return jsonError(c, 400, "name must only contain letters, numbers, and underscores")
	}

	leaderboard, err := db.GetLeaderboardForPlayer(player)
	if err != nil {
		return jsonError(c, 500, err.Error())
	}

	return c.JSON(leaderboard)
}

func leaderboardAPI_game_bounded(c *fiber.Ctx) error {
	holder, _ := c.Locals(models.HOLDER_KEY).(models.Holder)
	db := holder.GetDatabaseProvider()

	game := c.Params("game")
	if game == "" {
		return jsonError(c, 400, "game parameter is required")
	}
	if !gameRegex.MatchString(game) {
		return jsonError(c, 400, "game must only contain letters, numbers, and underscores")
	}

	startS := c.Query("lower")
	endS := c.Query("upper")

	if startS == "" || endS == "" {
		return jsonError(c, 400, "lower and upper parameters are required")
	}

	start, err := strconv.Atoi(startS)
	if err != nil {
		return jsonError(c, 400, "start must be an integer")
	}

	end, err := strconv.Atoi(endS)
	if err != nil {
		return jsonError(c, 400, "end must be an integer")
	}
	if start > end {
		return jsonError(c, 400, "start must be less than end")
	}
	if start < 0 {
		return jsonError(c, 400, "start must be greater than or equal to 0")
	}
	if end > 200 {
		return jsonError(c, 400, "end must be less than or equal to 200")
	}

	leaderboard, err := db.GetLeaderboardBounded(game, start, end)
	if err != nil {
		return jsonError(c, 500, err.Error())
	}

	return c.JSON(leaderboard)
}

func leaderboardAPI_game(c *fiber.Ctx) error {
	holder, _ := c.Locals(models.HOLDER_KEY).(models.Holder)
	db := holder.GetDatabaseProvider()

	game := c.Params("game")
	if game == "" {
		return jsonError(c, 400, "game parameter is required")
	}
	if !gameRegex.MatchString(game) {
		return jsonError(c, 400, "game must only contain letters, numbers, and underscores")
	}
	leaderboard, err := db.GetLeaderboard(game)
	if err != nil {
		return jsonError(c, 500, err.Error())
	}

	return c.JSON(leaderboard)
}

func jsonError(c *fiber.Ctx, status int, error string) error {
	slog.Error(fmt.Sprintf("[%d] %s: %v", status, c.Path(), error))
	return c.Status(status).JSON(fiber.Map{"error": error})
}
