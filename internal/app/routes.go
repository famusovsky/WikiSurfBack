package app

// TODO получать роуты по названию и ссылку показывать тоже его

// setRoutes - устанавливает маршрутизацию.
func setRoutes(app *App) {
	app.web.Static("/static", "./ui/static")

	auth := app.web.Group("/auth")
	auth.Get("/", app.auth)
	auth.Post("/", app.signUp)
	auth.Delete("/", app.signOut)
	auth.Get("/signin", app.renderSignin)
	auth.Get("/signup", app.renderSignup)

	service := app.web.Group("/service", app.checkReg)
	service.Get("/rating/route/:route", app.getRouteRating)
	service.Get("/rating/tour/:tour", app.getTourRating)
	service.Get("/rating/", app.getRating)
	service.Get("/tours", app.renderOpenedTournaments)
	service.Get("/tours/my", app.renderUserTournaments)
	service.Get("/tours/created", app.renderCreatorTournaments)
	service.Post("/tour/participate/:id", app.participateViaId, app.renderTournament)
	service.Delete("/tour/participate/:id", app.quitViaId, app.renderTournament)
	service.Post("/tour/participate/", app.participateViaPassword)
	service.Get("/tour/create", app.createTour)
	service.Delete("/tour/:id", app.deleteTour)
	service.Put("/tour/:id/route", app.addRouteToTour)
	service.Delete("/tour/:id/route", app.removeRouteFromTour)
	service.Put("/tour/:id/creator", app.addCreatorToTour)
	service.Delete("/tour/:id/creator", app.removeCreatorFromTour)
	service.Put("/tour/:id", app.updateTour)
	service.Post("/route/create", app.createRoute)

	app.web.Get("/ext/auth", app.authExt)
	ext := app.web.Group("/ext", app.checkRegExt)
	ext.Get("/", app.renderMainExt)
	ext.Get("/start", app.renderStartExt)
	ext.Post("/start", app.startRouteExt)
	ext.Get("/routes", app.renderRoutesExt)
	ext.Get("/tours", app.renderToursExt)
	ext.Post("/sprint", app.addSprintExt)

	base := app.web.Group("/", app.checkReg)
	base.Get("/", app.renderMain)
	base.Get("/history", app.renderHistory)
	base.Get("/settings", app.renderSettings)
	base.Put("/service/user", app.updateUser)
	base.Get("/sprint/:id", app.renderSprint)
	base.Get("/route/:id", app.renderRoute)
	base.Get("/tournaments", app.renderTournaments)
	base.Get("/tournament/:id", app.renderTournament)
	base.All("/tournament/edit/:id", app.renderEditTour)
}
