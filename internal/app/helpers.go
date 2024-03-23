package app

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (app *App) checkReg(c *fiber.Ctx) error {
	_, ok := app.getUser(c, errors.New("error while checking authorization"))
	if !ok {
		return c.Redirect("/auth")
	}

	return c.Next()
}

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

type sprintData struct {
	Start      string
	Finish     string
	StartTime  string
	LengthTime string
	Steps      int
}

func (app *App) getFullSprintDate(sprint models.Sprint) sprintData {
	res := sprintData{}
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

func getToursTable(tours []models.Tournament) (string, error) {
	var b bytes.Buffer
	q := `{{range .}}<tr><td hx-get={{printf "/tournament/%d" .Id }} hx-target="body">{{.Id}}</td></tr>{{end}}`

	temp := template.Must(template.New("").Parse(q))
	if err := temp.Execute(&b, tours); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (app *App) errToResult(c *fiber.Ctx, err error, name ...string) error {
	result := "#result"
	if len(name) != 0 {
		result = name[0]
	}
	app.errLog.Print(err)
	c.Set("HX-Retarget", result)
	return c.SendString(err.Error())
}

func (app *App) renderErr(c *fiber.Ctx, status int, err error) error {
	app.errLog.Println(err)
	return c.Render("error", fiber.Map{
		"status":  status,
		"errText": err.Error(),
	}, "layouts/base")
}

func (app *App) renderSimpleRating(c *fiber.Ctx, ratings []models.TourRating, wrapErr error) error {
	var b bytes.Buffer
	q := `{{range .}}<tr><td>{{.UserName}}</td><td>{{.Points}}</td></tr>{{end}}`
	t := template.Must(template.New("").Parse(q))
	if err := t.Execute(&b, ratings); err != nil {
		return app.renderErr(c, fiber.StatusInternalServerError, errors.Join(wrapErr, err))
	}

	return c.SendString(b.String())
}
