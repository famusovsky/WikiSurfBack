package app

import (
	"errors"
	"strconv"

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

// TODO signout

func (app *App) getRating(c *fiber.Ctx) error {
	// TODO
	return c.SendString(`<tr><td>First</td><td>Second</td></tr>`)
}

// getRouteRating - функция, возвращяющая рейтинг по маршруту.
func (app *App) getRouteRating(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting route ratings in api")
	id, err := strconv.Atoi(c.Params("route"))
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.ErrBadRequest
	}

	ratings, err := app.db.GetRouteRatings(id)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.ErrInternalServerError
	}

	app.infoLog.Printf("ratings for route %d is successfully getted\n", id)
	return c.JSON(ratings)
}

// getTourRating - функция, возвращяющая рейтинг по соревнованию.
func (app *App) getTourRating(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting tournament ratings in api")
	id, err := strconv.Atoi(c.Params("tour"))
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.ErrBadRequest
	}

	ratings, err := app.db.GetTournamentRatings(id)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	app.infoLog.Printf("ratings for tour %d is successfully getted\n", id)
	return c.JSON(ratings)
}

func (app *App) updateUser(c *fiber.Ctx) error {
	wrapErr := errors.New("error while updating user")
	user, ok := app.getUser(c, wrapErr)
	if !ok {
		return wrapErr
	}

	creds := struct {
		Name     string
		Email    string
		Password string
	}{}
	c.BodyParser(&creds)
	if creds.Name != "" {
		user.Name = creds.Name
	}
	if creds.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
		if err != nil {
			app.errLog.Println(errors.Join(wrapErr, err))
			return errors.Join(wrapErr, err)
		}

		user.Password = string(hashedPassword)
	}
	if creds.Email != "" {
		user.Email = creds.Email
	}

	err := app.db.UpdateUser(user)
	if err != nil {
		return err
	}

	return c.SendString("Ok")
}
