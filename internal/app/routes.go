package app

import (
	"github.com/gofiber/fiber/v2"
)

// TODO rating

func setRoutes(app *App) {
	app.web.Static("/static", "./ui/static")

	service := app.web.Group("/service")

	service.Get("/rating/route/:route", app.getRouteRating)
	service.Get("/rating/tour/:tour", app.getTourRating)
	service.Get("/rating/", app.getRating)
	service.Get("/tours", app.renderOpenedTournaments)
	service.Get("/tours/my", app.renderUserTournaments)
	service.Get("/tours/created", app.renderCreatorTournaments)
	service.Get("/signin", app.renderSignin)
	service.Get("/signup", app.renderSignup) // FIXME

	app.web.Get("/auth", app.renderAuth)
	app.web.Post("/auth", app.signUp)
	app.web.Get("/", app.checkReg, app.renderMain)
	app.web.Get("/history", app.checkReg, app.renderHistory)
	app.web.Get("/settings", app.checkReg, app.renderSettings)
	app.web.Put("/service/user", app.checkReg, app.updateUser)
	app.web.Get("/sprint/:id", app.checkReg, app.renderSprint)
	app.web.Get("/route/:id", app.checkReg, app.renderRoute)
	app.web.Get("/tournaments", app.checkReg, app.renderTournaments)
	app.web.Get("/tournament/:id", app.checkReg, app.renderTournament)

	app.web.Post("/create-route", func(c *fiber.Ctx) error {
		return c.Format(fiber.Map{
			"error": "already existes",
		})
	}) // TODO
}
