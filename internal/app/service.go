package app

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// getRouteRating - функция, возвращяющая рейтинг по маршруту.
func (app *App) getRouteRating(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting route ratings in api")
	id, err := strconv.Atoi(c.Params("route"))
	if err != nil {
		return app.renderErr(c, fiber.StatusBadRequest, errors.Join(wrapErr, err))
	}

	ratings, err := app.db.GetRouteRatings(id)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}
	ratingsData := make([]struct {
		Id     int
		Name   string
		Length string
		Steps  string
	}, len(ratings))
	for i := 0; i < len(ratings); i++ {
		sprint, _ := app.db.GetSprint(ratings[i].SprintId) // TODO get route ratings directly from db
		data := app.getFullSprintData(sprint)

		ratingsData[i].Id = ratings[i].SprintId

		s := ratings[i].SprintLengthTime / 1000
		ms := ratings[i].SprintLengthTime % 1000
		min := s / 60
		s = min % 60
		ratingsData[i].Length = fmt.Sprintf("%d min, %d s, %d ms", min, s, ms)

		ratingsData[i].Steps = strconv.Itoa(data.Steps)

		usr, err := app.db.GetUserById(ratings[i].UserId)
		if err != nil {
			ratingsData[i].Name = fmt.Sprintf("User with id:%d", ratings[i].UserId)
		} else {
			ratingsData[i].Name = usr.Name
		}
	}

	var b bytes.Buffer
	q := `{{range .}}<tr hx-get={{printf "/sprint/%d" .Id }} hx-target="body">
	<td>{{.Name}}</td>
	<td>{{.Length}}</td>
	<td>{{.Steps}}</td>
	</tr>{{end}}`
	t := template.Must(template.New("").Parse(q))
	if err := t.Execute(&b, ratingsData); err != nil {
		return app.renderErr(c, fiber.StatusInternalServerError, errors.Join(wrapErr, err))
	}

	return c.SendString(b.String())
}

// getTourRating - функция, возвращяющая рейтинг по соревнованию.
func (app *App) getTourRating(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting tournament ratings in api")
	id, err := strconv.Atoi(c.Params("tour"))
	if err != nil {
		return app.renderErr(c, fiber.StatusBadRequest, errors.Join(wrapErr, err))
	}

	ratings, err := app.db.GetTournamentRatings(id)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}

	return app.renderSimpleRating(c, ratings, wrapErr)
}

// getRating - функция, возвращяющая общий рейтинг.
func (app *App) getRating(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting ratings in api")

	ratings, err := app.db.GetRatings()
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}

	return app.renderSimpleRating(c, ratings, wrapErr)
}

// updateUser - функция, обновляющая данные пользователя.
func (app *App) updateUser(c *fiber.Ctx) error {
	wrapErr := errors.New("error while updating user")
	user, _ := app.getUser(c, wrapErr)

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
			return app.errToResult(c, errors.Join(wrapErr, err))
		}

		user.Password = string(hashedPassword)
	}
	if creds.Email != "" {
		user.Email = creds.Email
	}

	err := app.db.UpdateUser(user)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	return app.renderSettings(c)
}

// participateViaId - функция, добавляющая пользователя в соревнование по id.
func (app *App) participateViaId(c *fiber.Ctx) error {
	wrapErr := errors.New("error while adding user to tour")
	user, _ := app.getUser(c, wrapErr)

	tourId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	err = app.db.AddUserToTour(tourId, user.Id)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	return c.Next()
}

// quitViaId - функция, удаляющая пользователя из соревнования по id.
func (app *App) quitViaId(c *fiber.Ctx) error {
	wrapErr := errors.New("error while removing user from the tour")
	user, _ := app.getUser(c, wrapErr)

	tourId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	err = app.db.RemoveUserFromTour(tourId, user.Id)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	return c.Next()
}

// participateViaPassword - функция, добавляющая пользователя в соревнование по паролю.
func (app *App) participateViaPassword(c *fiber.Ctx) error {
	wrapErr := errors.New("error while adding user to tour")
	user, _ := app.getUser(c, wrapErr)

	pswd := struct {
		Password string
	}{}
	if err := c.BodyParser(&pswd); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	id, err := app.db.CheckTournamentPassword(pswd.Password)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	err = app.db.AddUserToTour(id, user.Id)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	return c.SendString(fmt.Sprintf("You are entered tour #%d", id))
}

// createRoute - функция, создающая маршрут.
func (app *App) createRoute(c *fiber.Ctx) error {
	route := models.Route{}
	wrapErr := errors.New("error while creating route")

	user, _ := app.getUser(c, wrapErr)

	if err := c.BodyParser(&route); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	if r := regexp.MustCompile(`.*wikipedia\.org\/wiki\/[^\s"]+`); !(r.Match([]byte(route.Start)) && r.Match([]byte(route.Finish))) {
		return app.errToResult(c, errors.Join(wrapErr, errors.New("input must be a wikipedia article link")))
	}

	route.CreatorId = user.Id
	id, err := app.db.AddRoute(route)

	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	return c.Redirect(fmt.Sprintf("/route/%d", id))
}

// createTour - функция, создающая соревнование.
func (app *App) createTour(c *fiber.Ctx) error {
	var pswd strings.Builder
	getRand := func(out *strings.Builder) {
		charSet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
		for i := 0; i < 32; i++ {
			random := rand.Intn(len(charSet))
			randomChar := charSet[random]
			out.WriteRune(randomChar)
		}
	}

	for {
		getRand(&pswd)
		_, err := app.db.CheckTournamentPassword(pswd.String())
		if err != nil {
			break
		}
		pswd.Reset()
	}

	t := models.Tournament{
		StartTime: time.Now(),
		EndTime:   time.Now().AddDate(0, 0, 7),
		Private:   true,
		Pswd:      pswd.String(),
	}

	wrapErr := errors.New("error while creating tour")

	user, _ := app.getUser(c, wrapErr)

	id, err := app.db.AddTournament(t, user.Id)
	if err != nil {
		return app.renderErr(c, fiber.StatusInternalServerError, errors.Join(wrapErr, err))
	}

	return c.Redirect(fmt.Sprintf("/tournament/edit/%d", id))
}

// updateTour - функция, обновляющая соревнование.
func (app *App) updateTour(c *fiber.Ctx) error {
	wrapErr := errors.New("error while editing the tour")

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	user, _ := app.getUser(c, wrapErr)
	ok, err := app.db.CheckTournamentCreator(id, user.Id)
	if !ok || err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	times := struct {
		Begin string
		End   string
	}{}
	if err := c.BodyParser(&times); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	tour, err := app.db.GetTournament(id)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}
	if t, err := time.Parse("2006-01-02T15:04:00Z", times.Begin+":00Z"); err == nil && !t.IsZero() {
		tour.StartTime = t
	}
	if t, err := time.Parse("2006-01-02T15:04:00Z", times.End+":00Z"); err == nil && !t.IsZero() {
		tour.EndTime = t
	}

	if err := app.db.UpdateTournament(tour, user.Id); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	return c.Redirect(fmt.Sprintf("/tournament/edit/%d", id))
}

// toggleTourPrivacy - функция изменяющая значение приватности соревнования.
func (app *App) toggleTourPrivace(c *fiber.Ctx) error {
	wrapErr := errors.New("error while toggling tour privacy")
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}
	user, _ := app.getUser(c, wrapErr)
	ok, err := app.db.CheckTournamentCreator(id, user.Id)
	if !ok || err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}
	tour, err := app.db.GetTournament(id)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}
	tour.Private = !tour.Private
	if err := app.db.UpdateTournament(tour, user.Id); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	return c.Redirect(fmt.Sprintf("/tournament/edit/%d", id))
}

// deleteTour - функция, удаляющая соревнование.
func (app *App) deleteTour(c *fiber.Ctx) error {
	wrapErr := errors.New("error while deleting the tour")

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	user, _ := app.getUser(c, wrapErr)
	ok, err := app.db.CheckTournamentCreator(id, user.Id)
	if !ok || err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	if err := app.db.DeleteTournament(id, user.Id); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	return c.Redirect("/tournaments")
}

// addRouteToTour - функция, добавляющая маршрут в соревнование.
func (app *App) addRouteToTour(c *fiber.Ctx) error {
	wrapErr := errors.New("error while adding route to the tour")

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#routesResult")
	}

	user, _ := app.getUser(c, wrapErr)
	ok, err := app.db.CheckTournamentCreator(id, user.Id)
	if !ok || err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#routesResult")
	}

	route, err := app.getOrCreateRoute(c)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#routesResult")
	}

	if err := app.db.AddRouteToTour(models.TRRelation{
		TournamentId: id,
		RouteId:      route.Id,
	}, user.Id); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#routesResult")
	}

	return c.Redirect(fmt.Sprintf("/tournament/edit/%d", id))
}

// removeRouteFromTour - функция, удаляющая маршрут из соревнования.
func (app *App) removeRouteFromTour(c *fiber.Ctx) error {
	wrapErr := errors.New("error while removing route from the tour")

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#routesResult")
	}

	user, _ := app.getUser(c, wrapErr)
	ok, err := app.db.CheckTournamentCreator(id, user.Id)
	if !ok || err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#routesResult")
	}

	route := models.Route{}
	if err := c.BodyParser(&route); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#routesResult")
	}
	route.CreatorId = user.Id

	if r, err := app.db.GetRouteByCreds(route.Start, route.Finish); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#routesResult")
	} else {
		route = r
	}

	if err := app.db.RemoveRouteFromTour(models.TRRelation{
		TournamentId: id,
		RouteId:      route.Id,
	}, user.Id); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#routesResult")
	}

	return c.Redirect(fmt.Sprintf("/tournament/edit/%d", id))
}

// addCreatorToTour - функция, добавляющая создателя в соревнование.
func (app *App) addCreatorToTour(c *fiber.Ctx) error {
	wrapErr := errors.New("error while adding creator to the tour")

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#creatorResult")
	}

	user, _ := app.getUser(c, wrapErr)
	ok, err := app.db.CheckTournamentCreator(id, user.Id)
	if !ok || err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#creatorResult")
	}

	email := struct {
		Email string
	}{}

	if err := c.BodyParser(&email); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#creatorResult")
	}
	creator, err := app.db.GetUser(email.Email)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#creatorResult")
	}

	if err := app.db.AddCreatorToTour(models.TURelation{
		TournamentId: id,
		UserId:       creator.Id,
	}, user.Id); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#creatorResult")
	}

	return c.Redirect(fmt.Sprintf("/tournament/edit/%d", id))
}

// removeCreatorFromTour - функция, удаляющая создателя из соревнования.
func (app *App) removeCreatorFromTour(c *fiber.Ctx) error {
	wrapErr := errors.New("error while adding creator to the tour")

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#creatorResult")
	}

	user, _ := app.getUser(c, wrapErr)
	ok, err := app.db.CheckTournamentCreator(id, user.Id)
	if !ok || err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#creatorResult")
	}

	email := struct {
		Email string
	}{}

	if err := c.BodyParser(&email); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#creatorResult")
	}
	creator, err := app.db.GetUser(email.Email)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#creatorResult")
	}

	if err := app.db.RemoveCreatorFromTour(models.TURelation{
		TournamentId: id,
		UserId:       creator.Id,
	}, user.Id); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err), "#creatorResult")
	}

	return c.Redirect(fmt.Sprintf("/tournament/edit/%d", id))
}
