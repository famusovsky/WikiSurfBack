package app

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/gofiber/fiber/v2"
)

// getUser - функция, возвращающая пользователя по fiber.Ctx.
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

func getToursTable(tours []models.Tournament) string {
	var b *bytes.Buffer
	q := `{{range .}}<tr><td hx-get={{printf "/tour/%s" .Id }}>{{.Id}}</td></tr>{{end}}`

	temp := template.Must(template.New("").Parse(q))

	if err := temp.Execute(b, tours); err != nil {
		log.Fatal(err)
	}

	return b.String()
}
