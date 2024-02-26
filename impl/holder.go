package impl

import (
	"github.com/Fesaa/CubepanionAPI/models"
	"github.com/gofiber/fiber/v2"
)

type holderImpl struct {
	databaseProvider models.DatabaseProvider
}

func NewHolder(dbURL string) (models.Holder, error) {
	db, err := newDatabaseProvider(dbURL)
	if err != nil {
		return nil, err
	}

	return &holderImpl{
		databaseProvider: db,
	}, nil
}

func SetHolderMiddelware(h models.Holder) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals(models.HOLDER_KEY, h)
		return c.Next()
	}
}

func (h *holderImpl) GetDatabaseProvider() models.DatabaseProvider {
	return h.databaseProvider
}
