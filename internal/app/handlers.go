package app

import (
	"bytes"
	"errors"
	"html/template"
	"strconv"
	"sync"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (app *App) auth(c *fiber.Ctx) error {
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
		"ratingType": "/service/rating",
	}, "layouts/base")
}

func (app *App) renderHistory(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting users history")
	user, _ := app.getUser(c, wrapErr)

	history, err := app.db.GetUserHistory(user.Id)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
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

	q := `{{range .}}<tr hx-get={{printf "/sprint/%d" .id }} hx-target="body">
	<td>{{.Start}}</td>
	<td>{{.Finish}}</td>
	<td>{{.StartTime}}</td>
	<td>{{.LengthTime}}</td>
	<td>{{.Steps}}</td>
	</tr>{{end}}`
	t := template.Must(template.New("").Parse(q))

	var body bytes.Buffer
	if err := t.Execute(&body, res); err != nil {
		return app.renderErr(c, fiber.StatusInternalServerError, errors.Join(wrapErr, err))
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
	wrapErr := errors.New("error while getting sprint data")
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.renderErr(c, fiber.StatusBadRequest, errors.Join(wrapErr, err))
	}
	sprint, err := app.db.GetSprint(id)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}

	var infoTbody, stepsTbody bytes.Buffer
	wg := sync.WaitGroup{}
	errs := make([]error, 3)
	wg.Add(3)

	go func(b *bytes.Buffer, wg *sync.WaitGroup, e []error) {
		q := `{{range .}}<tr><td>{{.Start}}</td><td>{{.Finish}}</td><td>{{.StartTime}}</td><td>{{.LengthTime}}</td><td>{{.Steps}}</td></tr>{{end}}`

		data := app.getFullSprintDate(sprint)
		t := template.Must(template.New("").Parse(q))

		if err := t.Execute(b, data); err != nil {
			e[0] = err
		}
		wg.Done()
	}(&infoTbody, &wg, errs)

	go func(b *bytes.Buffer, wg *sync.WaitGroup, e []error) {
		q := `{{range .}}<tr><td>{{.}}</td></tr>{{end}}`
		t := template.Must(template.New("").Parse(q))

		if err := t.Execute(b, sprint.Path); err != nil {
			e[1] = err
		}
		wg.Done()
	}(&stepsTbody, &wg, errs)

	var place int
	go func(place *int, wg *sync.WaitGroup, e []error) {
		rating, err := app.db.GetRouteRatings(sprint.RouteId)
		if err == nil {
			for i := range rating {
				if rating[i].UserId == sprint.UserId {
					*place = i + 1
					break
				}
			}
		} else {
			e[2] = err
		}
		wg.Done()
	}(&place, &wg, errs)

	for _, err := range errs {
		if err != nil {
			return app.renderErr(c, fiber.StatusInternalServerError, err)
		}
	}

	return c.Render("sprint", fiber.Map{
		"ind":        strconv.Itoa(sprint.Id),
		"infoTbody":  infoTbody.String(),
		"place":      strconv.Itoa(place),
		"routeId":    strconv.Itoa(sprint.RouteId),
		"stepsTbody": stepsTbody.String(),
	}, "layouts/base")
}

func (app *App) renderRoute(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting a route")
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.renderErr(c, fiber.StatusBadRequest, errors.Join(wrapErr, err))
	}

	route, err := app.db.GetRoute(id)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}

	return c.Render("route", fiber.Map{
		"ind":        c.Params("id"),
		"start":      route.Start,
		"finish":     route.Finish,
		"link":       route.Start, // FIXME сделать конпкой с хедерами
		"ratingType": "/service/rating/route/" + c.Params("id"),
	}, "layouts/base")
}

func (app *App) renderTournaments(c *fiber.Ctx) error {
	return c.Render("tournaments", fiber.Map{}, "layouts/base")
}

func (app *App) renderOpenedTournaments(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting opened tours")
	tours, err := app.db.GetOpenTournaments()
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}

	res, err := getToursTable(tours)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, wrapErr)
	}
	return c.Render("partials/tourList", fiber.Map{
		"name":  "Opened tours",
		"tbody": res,
	})
}

func (app *App) renderUserTournaments(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting user tours")
	user, _ := app.getUser(c, wrapErr)

	tours, err := app.db.GetUserTournaments(user.Id)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}

	res, err := getToursTable(tours)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, wrapErr)
	}
	return c.Render("partials/tourList", fiber.Map{
		"name":  "Tours in which I participate",
		"tbody": res,
	})
}

func (app *App) renderCreatorTournaments(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting creator tours")
	user, _ := app.getUser(c, wrapErr)

	tours, err := app.db.GetCreatorTournaments(user.Id)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, wrapErr)
	}

	res, err := getToursTable(tours)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, wrapErr)
	}
	return c.Render("partials/tourList", fiber.Map{
		"name":  "Tours I have created",
		"tbody": res,
	})
}

func (app *App) renderTournament(c *fiber.Ctx) error {
	wrapErr := errors.New("error while getting tour")
	user, _ := app.getUser(c, errors.New("error while getting user"))

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}

	var (
		body         bytes.Buffer
		tour         models.Tournament
		participates bool
		isCreator    bool
	)

	wg := sync.WaitGroup{}
	errs := make([]error, 4)
	wg.Add(4)

	go func(t *models.Tournament, wg *sync.WaitGroup, e []error) {
		if tmp, err := app.db.GetTournament(id); err == nil {
			*t = tmp
		} else {
			e[0] = err
		}

		wg.Done()
	}(&tour, &wg, errs)

	go func(b *bytes.Buffer, wg *sync.WaitGroup, e []error) {
		routes, err := app.db.GetTournamentRoutes(id)
		if err != nil {
			e[1] = err
			wg.Done()
			return
		}
		q := `{{range .}}<tr hx-get={{printf "/route/%d" .Id }} hx-target="body">
	<td>{{.Id}}</td>
	<td>{{.Start}}</td>
	<td>{{.Finish}}</td>
	</tr>{{end}}`
		t := template.Must(template.New("").Parse(q))
		if err := t.Execute(b, routes); err != nil {
			e[1] = err
		}
		wg.Done()
	}(&body, &wg, errs)

	go func(b *bool, wg *sync.WaitGroup, e []error) {
		if participates, err := app.db.CheckTournamentParticipator(id, user.Id); err == nil {
			*b = participates
		} else {
			e[2] = err
		}
		wg.Done()
	}(&participates, &wg, errs)

	go func(b *bool, wg *sync.WaitGroup, e []error) {
		if isCreator, err := app.db.CheckTournamentCreator(id, user.Id); err == nil {
			*b = isCreator
		} else {
			e[3] = err
		}
		wg.Done()
	}(&isCreator, &wg, errs)

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return app.renderErr(c, fiber.StatusInternalServerError, err)
		}
	}

	return c.Render("tournament", fiber.Map{
		"ind":          c.Params("id"),
		"password":     tour.Pswd,
		"routesTbody":  body.String(),
		"participates": participates,
		"isCreator":    isCreator,
		"ratingType":   "/service/rating/tour/" + c.Params("id"),
		"start":        tour.StartTime.Format("2006 Jan 2 15:04"),
		"end":          tour.EndTime.Format("2006 Jan 2 15:04"),
	}, "layouts/base")
}

func (app *App) renderEditTour(c *fiber.Ctx) error {
	wrapErr := errors.New("error while rendering tour editor")

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return app.renderErr(c, fiber.StatusBadRequest, errors.Join(wrapErr, err))
	}

	user, _ := app.getUser(c, wrapErr)
	ok, err := app.db.CheckTournamentCreator(id, user.Id)
	if !ok || err != nil {
		return app.renderErr(c, fiber.StatusForbidden, errors.Join(wrapErr, err))
	}

	tour, err := app.db.GetTournament(id)
	if err != nil {
		return app.renderErr(c, fiber.StatusNotFound, errors.Join(wrapErr, err))
	}

	var creatorsTbody, routesTbody bytes.Buffer
	wg := sync.WaitGroup{}
	errs := make([]error, 2)
	wg.Add(2)

	go func(b *bytes.Buffer, wg *sync.WaitGroup, e []error) {
		q := `{{range .}}<tr><td>{{.Name}}</td><td>{{.Email}}</td></tr>{{end}}`

		if users, err := app.db.GetTournamentCreators(id); err == nil {
			t := template.Must(template.New("").Parse(q))
			if err := t.Execute(b, users); err != nil {
				e[0] = err
			}
		} else {
			e[0] = err
		}
		wg.Done()
	}(&creatorsTbody, &wg, errs)

	go func(b *bytes.Buffer, wg *sync.WaitGroup, e []error) {
		q := `{{range .}}<tr hx-get={{printf "/route/%d" .Id }} hx-target="body"><td>{{.Id}}</td><td>{{.Start}}</td><td>{{.Finish}}</td></tr>{{end}}`

		if routes, err := app.db.GetTournamentRoutes(id); err == nil {
			t := template.Must(template.New("").Parse(q))
			if err := t.Execute(b, routes); err != nil {
				e[1] = err
			}
		} else {
			e[1] = err
		}
		wg.Done()
	}(&routesTbody, &wg, errs)

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return app.renderErr(c, fiber.StatusInternalServerError, err)
		}
	}

	return c.Render("editTour", fiber.Map{
		"start":         tour.StartTime.Format("2006 Jan 2 15:04"),
		"end":           tour.EndTime.Format("2006 Jan 2 15:04"),
		"ind":           c.Params("id"),
		"routesTbody":   routesTbody.String(),
		"creatorsTbody": creatorsTbody.String(),
		"password":      tour.Pswd,
	}, "layouts/base")
}
