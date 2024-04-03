package app

import (
	"html/template"
	"log"
	"net/http"

	"github.com/famusovsky/WikiSurfBack/internal/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// App - структура, представляющая собой приложение.
type App struct {
	web     *fiber.App         // web - веб-приложение на основе фреймворка Fiber.
	db      postgres.DbHandler // db - обработчик БД.
	ch      cookieHandler      // ch - обработчик cookie.
	infoLog *log.Logger        // infoLog - логгер информации.
	errLog  *log.Logger        // errorLog - логгер ошибок.
}

// CreateApp - создание приложения.
//
// Принимает: логгеры, обработчик БД.
//
// Возвращает: приложение.
func CreateApp(db postgres.DbHandler, infoLog, errLog *log.Logger) *App {
	engine := html.New("./ui/views", ".html")
	engine.AddFunc(
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)

	application := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			errLog.Printf("%s: %s %v", c.OriginalURL(), c.Method(), err)
			return c.Render("error", fiber.Map{
				"status":  http.StatusInternalServerError,
				"errText": err.Error(),
			}, "layouts/mini")
		},
		Views: engine,
	})

	result := &App{
		web:     application,
		db:      db,
		ch:      getCookieHandler("user-info", "email"),
		infoLog: infoLog,
		errLog:  errLog,
	}

	setRoutes(result)

	return result
}

// Run - запуск приложения.
//
// Принимает: адрес.
func (app *App) Run(addr string) {
	go app.errLog.Fatalln(app.web.Listen(addr))
}

// Shutdown - изящное отключение сервера.
func (app *App) Shutdown() error {
	return app.web.Shutdown()
}
