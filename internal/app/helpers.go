package app

import (
	"errors"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (app *App) getUser(c *fiber.Ctx, wrapErr error) (models.User, bool) {
	email, err := app.ch.Read(c)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return models.User{}, false
	}

	user, err := app.db.GetUser(email)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return models.User{}, false
	}

	return user, true
}
