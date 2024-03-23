package app

import (
	"bytes"
	"errors"
	"html/template"

	"github.com/gofiber/fiber/v2"
)

func (app *App) renderStartExt(c *fiber.Ctx) error {
	return c.Render("ext/start", fiber.Map{})
}

func (app *App) authExt(c *fiber.Ctx) error {
	return c.Render("ext/auth", fiber.Map{})
}

func (app *App) renderMainExt(c *fiber.Ctx) error {
	return c.Render("ext/main", fiber.Map{})
}

func (app *App) renderToursExt(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting user tours")
	user, _ := app.getUser(c, wrapErr)

	tours, err := app.db.GetUserTournaments(user.Id)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}

	var b bytes.Buffer
	q := `{{range .}}<tr><td><a href={{printf "http://127.0.0.1:8080/tournament/%d" .Id }} target="_blank">#{{.Id}}</a></td></tr>{{end}}`

	temp := template.Must(template.New("").Parse(q))
	if err := temp.Execute(&b, tours); err != nil {
		return app.renderErr(c, fiber.StatusNotFound, wrapErr)
	}

	return c.Render("partials/tourList", fiber.Map{
		"name":  "Tours in which I participate",
		"tbody": b.String(),
	})
}

func (app *App) renderRoutesExt(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting routes")

	routes, err := app.db.GetPopularRoutes()
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}

	var b bytes.Buffer
	q := `{{range .}}<tr>
	<td><a href={{printf "http://127.0.0.1:8080/route/%d" .Id }} target="_blank">#{{.Id}}</a></td>
	<td>{{.Start}}</td><td>{{.Finish}}</td>
	</tr>{{end}}`

	temp := template.Must(template.New("").Parse(q))
	if err := temp.Execute(&b, routes); err != nil {
		return app.renderErr(c, fiber.StatusNotFound, wrapErr)
	}

	return c.Render("partials/tourList", fiber.Map{
		"tbody": b.String(),
	})
}
