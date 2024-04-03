package app

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/gofiber/fiber/v2"
)

// renderStartExt - функция, производящая рендер страницы запуска спринта в расширении.
func (app *App) renderStartExt(c *fiber.Ctx) error {
	return c.Render("ext/start", fiber.Map{
		"baseUrl": c.BaseURL(),
	})
}

// startRouteExt - функция, запускающая прохождение спринта в расширении.
func (app *App) startRouteExt(c *fiber.Ctx) error {
	wrapErr := errors.New("error while starting a route")

	route, err := app.getOrCreateRoute(c, wrapErr, "#result")
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	return c.Render("ext/startRoute", fiber.Map{
		"rid":     route.Id,
		"start":   route.Start,
		"finish":  route.Finish,
		"baseUrl": c.BaseURL(),
	})
}

// authExt - функция, проводящая авторизацию в расширении.
func (app *App) authExt(c *fiber.Ctx) error {
	return c.Render("ext/auth", fiber.Map{
		"baseUrl": c.BaseURL(),
	})
}

// renderMainExt - функция, производящая рендер главной страницы расширения.
func (app *App) renderMainExt(c *fiber.Ctx) error {
	return c.Render("ext/main", fiber.Map{
		"baseUrl": c.BaseURL(),
	})
}

// renderToursExt - функция, производящая рендер страницы соревнований в расширении.
func (app *App) renderToursExt(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting user tours")
	user, _ := app.getUser(c, wrapErr)

	tours, err := app.db.GetUserTournaments(user.Id)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err), "")
	}

	var b bytes.Buffer
	q := fmt.Sprintf(
		`{{range .}}<tr><td><a href={{printf "%s/tournament/%%d" .Id }} target="_blank">#{{.Id}}</a></td></tr>{{end}}`,
		c.BaseURL())

	temp := template.Must(template.New("").Parse(q))
	if err := temp.Execute(&b, tours); err != nil {
		return app.renderErr(c, fiber.StatusNotFound, wrapErr, "")
	}

	return c.Render("partials/tourList", fiber.Map{
		"name":    "Tours in which I participate",
		"tbody":   b.String(),
		"baseUrl": c.BaseURL(),
	})
}

// renderRoutesExt - функция, производящая рендер страницы маршрутов в расширении.
func (app *App) renderRoutesExt(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting routes")

	routes, err := app.db.GetPopularRoutes()
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err), "")
	}

	var b bytes.Buffer
	q := fmt.Sprintf(`{{range .}}<tr>
	<td><a href={{printf "%s/route/%%d" .Id }} target="_blank">#{{.Id}}</a></td><td>{{.Start}}</td><td>{{.Finish}}</td>
	</tr>{{end}}`, c.BaseURL())

	temp := template.Must(template.New("").Parse(q))
	if err := temp.Execute(&b, routes); err != nil {
		return app.renderErr(c, fiber.StatusNotFound, wrapErr, "")
	}

	return c.Render("partials/tourList", fiber.Map{
		"tbody":   b.String(),
		"baseUrl": c.BaseURL(),
	})
}

// addSprintExt - функция, сохраняющая пройденный спринт.
func (app *App) addSprintExt(c *fiber.Ctx) error {
	sprint := models.Sprint{}
	wrapErr := errors.New("error while adding a sprint")

	user, _ := app.getUser(c, wrapErr)
	if err := c.BodyParser(&sprint); err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}
	sprint.UserId = user.Id
	id, err := app.db.AddSprint(sprint)
	if err != nil {
		return app.errToResult(c, errors.Join(wrapErr, err))
	}

	return c.SendString(fmt.Sprintf("%s/sprint/%d", c.BaseURL(), id))
}
