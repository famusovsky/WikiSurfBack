package app

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"strconv"
	"sync"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/gofiber/fiber/v2"
)

// TODO change log fatal to return error html

func (app *App) checkReg(c *fiber.Ctx) error {
	_, ok := app.getUser(c, errors.New("error while checking authorization"))
	if !ok {
		return c.Redirect("/auth")
	}

	return c.Next()
}

func (app *App) renderAuth(c *fiber.Ctx) error {
	if c.Query("signin") != "" {
		return app.signIn(c)
	}
	return c.Render("auth/auth", fiber.Map{}, "layouts/mini")
}

func (app *App) renderSignin(c *fiber.Ctx) error {
	return c.Render("auth/signin", fiber.Map{})
}

func (app *App) renderSignup(c *fiber.Ctx) error {
	return c.Render("auth/signup", fiber.Map{})
}

func (app *App) renderMain(c *fiber.Ctx) error {
	return c.Render("main", fiber.Map{
		"ratingType": "/rating",
	}, "layouts/base")
}

func (app *App) renderHistory(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting users history")
	user, _ := app.getUser(c, wrapErr)

	history, err := app.db.GetUserHistory(user.Id)
	if err != nil {
		app.errLog.Println(errors.Join(wrapErr, err))
		return fiber.ErrInternalServerError
	}

	res := make([]sprintData, len(history))
	getData := func(res []sprintData, i int) {
		res[i] = app.getFullSprintDate(history[i])
	}

	wg := sync.WaitGroup{}
	wg.Add(len(res))
	for i := range res {
		go getData(res, i)
		wg.Done()
	}
	wg.Wait()

	q := `{{range .}}<tr hx-get={{printf "/sprint/%s" .id }} hx-target="body">
	<td>{{.Start}}</td>
	<td>{{.Finish}}</td>
	<td>{{.StartTime}}</td>
	<td>{{.LengthTime}}</td>
	<td>{{.Steps}}</td>
	</tr>{{end}}`
	t := template.Must(template.New("").Parse(q))

	var body bytes.Buffer
	if err := t.Execute(&body, res); err != nil {
		log.Fatal(err)
	}

	return c.Render("history", fiber.Map{
		"tbody": body.String(),
	}, "layouts/base")
}

func (app *App) renderSettings(c *fiber.Ctx) error {
	usr, _ := app.getUser(c, errors.New(""))
	return c.Render("settings", fiber.Map{
		"email": usr.Email,
		"name":  usr.Name,
	}, "layouts/base")
}

func (app *App) renderSprint(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrNotFound
	}
	sprint, err := app.db.GetSprint(id)
	if err != nil {
		return fiber.ErrNotFound
	}

	var infoTbody, stepsTbody bytes.Buffer
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func(b *bytes.Buffer, wg *sync.WaitGroup) {
		q := `{{range .}}<tr>
	<td>{{.Start}}</td>
	<td>{{.Finish}}</td>
	<td>{{.StartTime}}</td>
	<td>{{.LengthTime}}</td>
	<td>{{.Steps}}</td>
	</tr>{{end}}`

		data := app.getFullSprintDate(sprint)
		t := template.Must(template.New("").Parse(q))

		if err := t.Execute(b, data); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}(&infoTbody, &wg)

	go func(b *bytes.Buffer, wg *sync.WaitGroup) {
		q := `{{range .}}<tr>
	<td>{{.}}</td>
	</tr>{{end}}`
		t := template.Must(template.New("").Parse(q))

		if err := t.Execute(b, sprint.Path); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}(&stepsTbody, &wg)

	var place int
	go func(place *int, wg *sync.WaitGroup) {
		rating, err := app.db.GetRouteRatings(sprint.RouteId)
		if err != nil {
			log.Fatal(err)
		}
		for i := range rating {
			if rating[i].UserId == sprint.UserId {
				*place = i + 1
				break
			}
		}
	}(&place, &wg)

	return c.Render("sprint", fiber.Map{
		"ind":        sprint.Id,
		"infoTbody":  infoTbody.String(),
		"place":      place,
		"routeId":    sprint.RouteId,
		"stepsTbody": stepsTbody.String(),
	}, "layouts/base")
}

func (app *App) renderRoute(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrNotFound
	}

	route, err := app.db.GetRoute(id)
	if err != nil {
		return fiber.ErrNotFound
	}

	return c.Render("route", fiber.Map{
		"ind":        id,
		"start":      route.Start,
		"finish":     route.Finish,
		"link":       route.Start, // FIXME
		"ratingType": "/service/rating/route/" + c.Params("id"),
	}, "layouts/base")
}

func (app *App) renderTournaments(c *fiber.Ctx) error {
	return c.Render("tournaments", fiber.Map{}, "layouts/base")
}

func (app *App) renderOpenedTournaments(c *fiber.Ctx) error {
	tours, err := app.db.GetOpenTournaments()
	if err != nil {
		log.Fatal(err)
	}

	res := getToursTable(tours)
	return c.Render("partials/tourList", fiber.Map{
		"name":  "Opened tours",
		"tbody": res,
	})
}

func (app *App) renderUserTournaments(c *fiber.Ctx) error {
	user, ok := app.getUser(c, errors.New("error while getting user"))
	if !ok {
		return c.Redirect("/auth")
	}
	tours, err := app.db.GetUserTournaments(user.Id)
	if err != nil {
		log.Fatal(err)
	}

	res := getToursTable(tours)
	return c.Render("partials/tourList", fiber.Map{
		"name":  "Tours in which I participate",
		"tbody": res,
	})
}
func (app *App) renderCreatorTournaments(c *fiber.Ctx) error {
	user, ok := app.getUser(c, errors.New("error while getting user"))
	if !ok {
		return c.Redirect("/auth")
	}
	tours, err := app.db.GetCreatorTournaments(user.Id)
	if err != nil {
		log.Fatal(err)
	}

	res := getToursTable(tours)
	return c.Render("partials/tourList", fiber.Map{
		"name":  "Tours I have created",
		"tbody": res,
	})
}

func (app *App) renderTournament(c *fiber.Ctx) error {
	user, ok := app.getUser(c, errors.New("error while getting user"))
	if !ok {
		return c.Redirect("/auth")
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrNotFound
	}

	var (
		body         bytes.Buffer
		tour         models.Tournament
		participates bool
		isCreator    bool
	)

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func(t *models.Tournament, wg *sync.WaitGroup) {
		tmp, err := app.db.GetTournament(id)
		if err != nil {
			log.Fatal(err)
		}
		*t = tmp
		wg.Done()
	}(&tour, &wg)

	go func(b *bytes.Buffer, wg *sync.WaitGroup) {
		routes, err := app.db.GetTournamentRoutes(id)
		if err != nil {
			log.Fatal(err)
		}

		q := `{{range .}}<tr hx-get={{printf "/route/%s" .Id }} hx-target="body">
	<td>{{.Id}}</td>
	<td>{{.Start}}</td>
	<td>{{.Finish}}</td>
	</tr>{{end}}`

		t := template.Must(template.New("").Parse(q))

		if err := t.Execute(b, routes); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}(&body, &wg)

	go func(b *bool, wg *sync.WaitGroup) {
		participates, err := app.db.CheckTournamentParticipator(id, user.Id)
		if err != nil {
			log.Fatal(err)
		}
		*b = participates
		wg.Done()
	}(&participates, &wg)

	go func(b *bool, wg *sync.WaitGroup) {
		isCreator, err := app.db.CheckTournamentCreator(id, user.Id)
		if err != nil {
			log.Fatal(err)
		}
		*b = isCreator
		wg.Done()
	}(&isCreator, &wg)

	wg.Wait()

	return c.Render("tournament", fiber.Map{
		"password":     tour.Pswd,
		"routesTbody":  body.String(),
		"participates": participates,
		"isCreator":    isCreator,
		"ratingType":   "/service/rating/tour/" + c.Params("id"),
		"start":        tour.StartTime.Format("2006 Jan 2 15:04"),
		"end":          tour.EndTime.Format("2006 Jan 2 15:04"),
	}, "layouts/base")
}
