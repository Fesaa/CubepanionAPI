package errors

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

var (
	DBError = errors.New("internal database error, try again later. Or contact the developer")
)

func AsFiberMap(err error) fiber.Map {
	return fiber.Map{
		"error": err.Error(),
	}
}
