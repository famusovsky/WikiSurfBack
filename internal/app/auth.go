package app

import (
	"errors"

	"github.com/badoux/checkmail"
	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// signUp - функция, регистрирующая нового пользователя.
func (app *App) signUp(c *fiber.Ctx) error {
	c.Accepts("json")
	var user models.User
	wrapErr := errors.New("error while signing up user in api")

	if err := c.BodyParser(&user); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	user.Password = string(hashedPassword)

	_, err = app.db.AddUser(user)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	app.ch.Set(c, user.Email)

	app.infoLog.Printf("user %s signed up\n", user.Name)

	c.Set("HX-Location", "/")
	return c.SendString("OK")
}

// signIn - функция, авторизирующая пользователя.
func (app *App) signIn(c *fiber.Ctx) error {
	wrapErr := errors.New("error while signing in user in api")
	creds := struct {
		Email, Password string
	}{}
	if err := c.BodyParser(&creds); err != nil || creds.Email == "" || creds.Password == "" {
		return app.errToResult(c, errors.Join(wrapErr, errors.Join(wrapErr, errors.New(`request's body is wrong`))))
	}

	user, err := app.db.GetUser(creds.Email)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	app.ch.Set(c, user.Email)

	app.infoLog.Printf("user %s signed in\n", user.Name)

	c.Set("HX-Location", "/")
	return c.SendString("OK")
}

// signIn - функция, позволяющая пользователю выйти из аккаунта.
func (app *App) signOut(c *fiber.Ctx) error {
	app.ch.Remove(c)

	c.Set("HX-Refresh", "true")
	return c.SendString("")
}
