package app

import (
	"bytes"
	"encoding/json"
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

	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	if err := decoder.Decode(&user); err != nil {
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

	return c.JSON("OK")
}

// signIn - функция, авторизирующая пользователя.
func (app *App) signIn(c *fiber.Ctx) error {
	c.Accepts("json")
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	wrapErr := errors.New("error while signing in user in api")

	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	if err := decoder.Decode(&creds); err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.NewError(fiber.StatusBadRequest, `request's body is wrong`)
	}

	user, err := app.db.GetUser(creds.Email)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.NewError(fiber.StatusInternalServerError, "error while checking user info")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.NewError(fiber.StatusUnauthorized, "wrong password")
	}

	app.ch.Set(c, user.Email)

	app.infoLog.Printf("user %s signed in\n", user.Name)

	return c.JSON("OK")
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

// getUserHistory - функция, возвращяющая получение историю спринтов пользователя.
func (app *App) getUserHistory(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting user history in api")

	user, ok := app.getUser(c, wrapErr)
	if !ok {
		return c.Redirect("/signup", fiber.StatusUnauthorized)
	}

	history, err := app.db.GetUserHistory(user.Id)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.ErrInternalServerError
	}

	app.infoLog.Printf("history for user %s is successfully getted\n", user.Name)
	return c.JSON(history)
}

// getUserRouteHistory - функция, возвращяющая получение историю спринтов пользователя по маршруту.
func (app *App) getUserRouteHistory(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting user route history in api")
	routeId, err := strconv.Atoi(c.Params("route"))
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.ErrBadRequest
	}

	user, ok := app.getUser(c, wrapErr)
	if !ok {
		return c.Redirect("/signup", fiber.StatusUnauthorized)
	}

	history, err := app.db.GetUserRouteHistory(user.Id, routeId)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.ErrInternalServerError
	}

	app.infoLog.Printf("history for route %d user %s is successfully getted\n", routeId, user.Name)
	return c.JSON(history)
}

// getOpenTournaments - функция, возвращяющая список соревнований, открытых для вступления.
func (app *App) getOpenTournaments(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting open tounaments in api")

	tournaments, err := app.db.GetOpenTournaments()
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.ErrInternalServerError
	}

	app.infoLog.Printf("open tournaments are successfully getted\n")
	return c.JSON(tournaments)
}

// getUserTournaments - функция, возвращяющая список соревнований, в которые вступил пользователь.
func (app *App) getUserTournaments(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting user tounaments in api")

	user, ok := app.getUser(c, wrapErr)
	if !ok {
		return c.Redirect("/signup", fiber.StatusUnauthorized)
	}

	tournaments, err := app.db.GetUserTournaments(user.Id)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.ErrInternalServerError
	}

	app.infoLog.Printf("tournaments of user %s are successfully getted\n", user.Name)
	return c.JSON(tournaments)
}

// getCreatorTournaments - функция, возвращяющая список соревнований, в которых ползователь является создателем.
func (app *App) getCreatorTournaments(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting creator tounaments in api")

	user, ok := app.getUser(c, wrapErr)
	if !ok {
		return c.Redirect("/signup", fiber.StatusUnauthorized)
	}

	tournaments, err := app.db.GetCreatorTournaments(user.Id)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.ErrInternalServerError
	}

	app.infoLog.Printf("tournaments of creator %s are successfully getted\n", user.Name)
	return c.JSON(tournaments)
}
