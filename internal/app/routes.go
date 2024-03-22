package app

// TODO получать роуты по названию и ссылку показывать тоже его

func setRoutes(app *App) {
	app.web.Static("/static", "./ui/static")

	auth := app.web.Group("/auth")
	auth.Get("/", app.auth)
	auth.Post("/", app.signUp)
	auth.Delete("/", app.signOut)

	service := app.web.Group("/service")

	service.Get("/rating/route/:route", app.checkReg, app.getRouteRating)
	service.Get("/rating/tour/:tour", app.checkReg, app.getTourRating)
	service.Get("/rating/", app.checkReg, app.getRating)
	service.Get("/tours", app.checkReg, app.renderOpenedTournaments)
	service.Get("/tours/my", app.checkReg, app.renderUserTournaments)
	service.Get("/tours/created", app.checkReg, app.renderCreatorTournaments)
	service.Post("/tour/participate/:id", app.checkReg, app.participateViaId, app.renderTournament)
	service.Delete("/tour/participate/:id", app.checkReg, app.quitParticipateViaId, app.renderTournament)
	service.Post("/tour/participate/", app.checkReg, app.participateViaPassword)
	service.Get("/tour/create", app.checkReg, app.createTour)
	service.Delete("/tour/:id", app.checkReg, app.deleteTour)
	service.Put("/tour/:id/route", app.checkReg, app.addRouteToTour)
	service.Delete("/tour/:id/route", app.checkReg, app.removeRouteFromTour)
	service.Put("/tour/:id/creator", app.checkReg, app.addCreatorToTour)
	service.Delete("/tour/:id/creator", app.checkReg, app.removeCreatorFromTour)
	service.Put("/tour/:id", app.checkReg, app.updateTour)
	service.Get("/signin", app.renderSignin)
	service.Get("/signup", app.renderSignup)
	service.Post("/route/create", app.checkReg, app.createRoute)

	app.web.Get("/", app.checkReg, app.renderMain)
	app.web.Get("/history", app.checkReg, app.renderHistory)
	app.web.Get("/settings", app.checkReg, app.renderSettings)
	app.web.Put("/service/user", app.checkReg, app.updateUser)
	app.web.Get("/sprint/:id", app.checkReg, app.renderSprint)
	app.web.Get("/route/:id", app.checkReg, app.renderRoute)
	app.web.Get("/tournaments", app.checkReg, app.renderTournaments)
	app.web.Get("/tournament/:id", app.checkReg, app.renderTournament)
	app.web.All("/tournament/edit/:id", app.checkReg, app.renderEditTour)
}
