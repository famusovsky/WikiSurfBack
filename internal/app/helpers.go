package app

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"regexp"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/gofiber/fiber/v2"
)

// checkAuth - middleware, проверяющий авторизацию пользователя.
func (app *App) checkReg(c *fiber.Ctx) error {
	_, ok := app.getUser(c, errors.New("error while checking authorization"))
	if !ok {
		return c.Redirect("/auth")
	}

	return c.Next()
}

// checkAuthExt - middleware, проверяющий авторизацию пользователя в расширении.
func (app *App) checkRegExt(c *fiber.Ctx) error {
	_, ok := app.getUser(c, errors.New("error while checking authorization"))
	if !ok {
		return c.Redirect("/ext/auth")
	}

	return c.Next()
}

// getUser - функция, возвращающая пользователя по fiber.Ctx.
func (app *App) getUser(c *fiber.Ctx, wrapErr error) (models.User, bool) {
	email, err := app.ch.Read(c)
	if err != nil {
		// app.errLog.Println(errors.Join(wrapErr, err))
		return models.User{}, false
	}

	user, err := app.db.GetUser(email)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return models.User{}, false
	}

	return user, true
}

// sprintData - структура, хранящая полные данные о спринте.
type sprintData struct {
	Id         int
	Start      string
	Finish     string
	StartTime  string
	LengthTime string
	Steps      int
}

// getFullSprintData - функция, возвращающая полные данные о спринте.
func (app *App) getFullSprintData(sprint models.Sprint) sprintData {
	res := sprintData{}
	res.Id = sprint.Id
	res.StartTime = sprint.StartTime.Format("2006 Jan 2 15:04")
	s := sprint.LengthTime / 1000
	ms := sprint.LengthTime % 1000
	min := s / 60
	s = min % 60
	res.LengthTime = fmt.Sprintf("%d min, %d s, %d ms", min, s, ms)
	res.Steps = len(sprint.Path)
	r, err := app.db.GetRoute(sprint.RouteId)
	if err != nil {
		res.Start = "Error"
		res.Finish = "Error"
	} else {
		res.Start = r.Start
		res.Finish = r.Finish
	}

	return res
}

// getToursTable - функция, возвращающая html таблицу соревнований.
func getToursTable(tours []models.Tournament) (string, error) {
	var b bytes.Buffer
	q := `{{range .}}<tr><td hx-get={{printf "/tournament/%d" .Id }} hx-target="body">{{.Id}}</td></tr>{{end}}`

	temp := template.Must(template.New("").Parse(q))
	if err := temp.Execute(&b, tours); err != nil {
		return "", err
	}

	return b.String(), nil
}

// errToResult - функция, встраивающая ошибку в поле #result.
func (app *App) errToResult(c *fiber.Ctx, err error, name ...string) error {
	result := "#result"
	if len(name) != 0 {
		result = name[0]
	}
	app.errLog.Print(err)
	c.Set("HX-Retarget", result)
	return c.SendString(err.Error())
}

// renderErr - функция, рендерящая страницу ошибки.
func (app *App) renderErr(c *fiber.Ctx, status int, err error) error {
	app.errLog.Println(err)
	return c.Render("error", fiber.Map{
		"status":  status,
		"errText": err.Error(),
	}, "layouts/base")
}

// renderSimpleRating - функция, рендерящая простую таблицу рейтинга.
func (app *App) renderSimpleRating(c *fiber.Ctx, ratings []models.TourRating, wrapErr error) error {
	var b bytes.Buffer
	q := `{{range .}}<tr><td>{{.UserName}}</td><td>{{.Points}}</td></tr>{{end}}`
	t := template.Must(template.New("").Parse(q))
	if err := t.Execute(&b, ratings); err != nil {
		return app.renderErr(c, fiber.StatusInternalServerError, errors.Join(wrapErr, err))
	}

	return c.SendString(b.String())
}

// getOrCreateRoute - функция, возвращающая маршрут по запросу или создающая новый.
func (app *App) getOrCreateRoute(c *fiber.Ctx, wrapErr error, resultAddr string) (models.Route, error) {
	user, _ := app.getUser(c, wrapErr)

	route := models.Route{}
	if err := c.BodyParser(&route); err != nil {
		return models.Route{}, app.errToResult(c, errors.Join(wrapErr, err), resultAddr)
	}
	route.CreatorId = user.Id

	if r, err := app.db.GetRouteByCreds(route.Start, route.Finish); err != nil {
		if r := regexp.MustCompile(`.*wikipedia\.org\/wiki\/[^\s"]+`); !(r.Match([]byte(route.Start)) && r.Match([]byte(route.Finish))) {
			return models.Route{}, app.errToResult(c, errors.Join(wrapErr, errors.New("input must be a wikipedia article link")))
		}

		if id, err := app.db.AddRoute(route); err != nil {
			route.Id = id
		} else {
			return models.Route{}, app.errToResult(c, errors.Join(wrapErr, err), resultAddr)
		}
	} else {
		route = r
	}

	return route, nil
}
