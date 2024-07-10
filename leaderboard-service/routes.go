package main

import (
	"fmt"
	"github.com/Fesaa/CubepanionAPI/core/errors"
	"github.com/Fesaa/CubepanionAPI/core/log"
	"regexp"
	"strconv"

	"github.com/Fesaa/CubepanionAPI/core"
	"github.com/Fesaa/CubepanionAPI/core/models"
	"github.com/Fesaa/CubepanionAPI/leaderboard-service/database"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/validator.v2"
)

var playerRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)
var gameRegex = regexp.MustCompile(`^[a-zA-Z0-9_ ]`)

func PlayerLeaderboard(ms core.MicroService[LeaderboardServiceConfig, database.Database], c *fiber.Ctx) error {
	player := c.Params("name")
	if !playerRegex.MatchString(player) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid player name",
		})
	}

	leaderboard, err := ms.DB().GetLeaderboardForPlayer(player)
	if err != nil {
		log.Error("Error getting leaderboard for player", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}

	return c.JSON(leaderboard)
}

func GameLeaderboardBounded(ms core.MicroService[LeaderboardServiceConfig, database.Database], c *fiber.Ctx) error {
	gameDisplayName, err := convertGame(ms, c.Params("game"))
	if err != nil {
		log.Error("Could not convert game", "error", err)
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("Could not convert game %s", c.Params("game")),
		})
	}

	startS := c.Query("lower")
	endS := c.Query("upper")

	if startS == "" || endS == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "start and end are required",
		})
	}

	start, err := strconv.Atoi(startS)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("start is not a number: %v", err),
		})
	}

	end, err := strconv.Atoi(endS)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("end is not a number: %v", err),
		})
	}
	if start > end || start < 0 || end > models.LEADERBOARD_SIZE || start == end {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid range",
		})
	}

	leaderboard, err := ms.DB().GetLeaderboardBounded(gameDisplayName, start, end)
	if err != nil {
		log.Error("Error getting bounded leaderboard", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}

	return c.JSON(leaderboard)
}

func GameLeaderboard(ms core.MicroService[LeaderboardServiceConfig, database.Database], c *fiber.Ctx) error {
	gameDisplayName, err := convertGame(ms, c.Params("game"))
	if err != nil {
		log.Error("Could not convert game", "error", err)
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("Could not convert game %s", c.Params("game")),
		})
	}

	leaderboard, err := ms.DB().GetLeaderboard(gameDisplayName)
	if err != nil {
		log.Error("Error getting leaderboard", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}

	return c.JSON(leaderboard)
}

func Submit(ms core.MicroService[LeaderboardServiceConfig, database.Database], c *fiber.Ctx) error {
	var submission models.LeaderboardSubmission
	err := c.BodyParser(&submission)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("error parsing submission: %v", err),
		})
	}

	err = validator.Validate(submission)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("error validating submission: %v", err),
		})
	}

	game, err := getGame(ms, submission.Game)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("game %s does not exist", submission.Game),
		})
	}

	submission.Game = game

	err = ms.DB().SubmitLeaderboard(submission)
	if err != nil {
		log.Error("Error submitting leaderboard", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}

	return c.SendStatus(202)
}

func convertGame(ms core.MicroService[LeaderboardServiceConfig, database.Database], game string) (string, error) {
	if game == "" {
		return "", fmt.Errorf("game name is required")
	}

	if !gameRegex.MatchString(game) {
		return "", fmt.Errorf("game name is invalid")
	}

	gameDisplayName, err := getGame(ms, game)
	if err != nil {
		return "", fmt.Errorf("game %s does not exist", game)
	}

	return gameDisplayName, nil
}

func BatchPlayerLeaderboard(ms core.MicroService[LeaderboardServiceConfig, database.Database], c *fiber.Ctx) error {
	var req models.GamePlayersRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("error parsing request: %v", err),
		})
	}

	err = validator.Validate(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("error validating request: %v", err),
		})
	}

	gameDisplayName, err := convertGame(ms, req.Game)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("Could not convert game %s", req.Game),
		})
	}

	req.Game = gameDisplayName
	rows, err := ms.DB().GetLeaderboardForPlayers(req)
	if err != nil {
		log.Error("Error getting leaderboard for players", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}

	return c.JSON(rows)
}

func GetAllPlayers(ms core.MicroService[LeaderboardServiceConfig, database.Database], c *fiber.Ctx) error {
	players, err := ms.DB().GetAllPlayers()
	if err != nil {
		log.Error("Error getting all players", "error", err)
		return c.Status(500).JSON(errors.AsFiberMap(errors.DBError))
	}

	return c.JSON(players)
}
