package app

import (
	"errors"

	"github.com/badoux/checkmail"
	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// TODO change error output to normal (->result)

// signUp - функция, регистрирующая нового пользователя.
func (app *App) signUp(c *fiber.Ctx) error {
	c.Accepts("json")
	var user models.User
	wrapErr := errors.New("error while signing up user in api")

	if err := c.BodyParser(&user); err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.NewError(fiber.StatusBadRequest, `request's body is wrong`)
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.NewError(fiber.StatusBadRequest, `email is wrong`)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.NewError(fiber.StatusInternalServerError, "error while hashing the password")
	}

	user.Password = string(hashedPassword)

	_, err = app.db.AddUser(user)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.NewError(fiber.StatusInternalServerError, "error while saving user info")
	}

	app.ch.Set(c, user.Email)

	app.infoLog.Printf("user %s signed up\n", user.Name)

	return c.Redirect("/")
}

// signIn - функция, авторизирующая пользователя.
func (app *App) signIn(c *fiber.Ctx) error {
	wrapErr := errors.New("error while signing in user in api")
	email, pswd := c.Query("email"), c.Query("password")
	if email == "" || pswd == "" {
		app.errLog.Println(errors.Join(wrapErr, errors.New(`request's body is wrong`)))
		return fiber.NewError(fiber.StatusBadRequest, `request's body is wrong`)
	}

	user, err := app.db.GetUser(email)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.NewError(fiber.StatusInternalServerError, "error while checking user info")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pswd)); err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.NewError(fiber.StatusUnauthorized, "wrong password")
	}

	app.ch.Set(c, user.Email)

	app.infoLog.Printf("user %s signed in\n", user.Name)

	return c.Redirect("/")
}

// signIn - функция, позволяющая пользователю выйти из аккаунта.
func (app *App) signOut(c *fiber.Ctx) error {
	app.ch.Remove(c)

	c.Set("HX-Refresh", "true")
	return c.SendString("")
}
