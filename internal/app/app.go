package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/famusovsky/WikiSurfBack/internal/postgres"
	"github.com/gofiber/fiber/v2"
)

// App - структура, описывающая приложение.
type App struct {
	web     *fiber.App         // web - веб-приложение на основе фреймворка Fiber.
	db      postgres.DbHandler // db - обработчик БД.
	infoLog *log.Logger        // infoLog - логгер информации.
	errLog  *log.Logger        // errorLog - логгер ошибок.
}

// CreateApp - создание приложения.
//
// Принимает: логгеры, обработчик БД.
//
// Возвращает: приложение.
func CreateApp(db postgres.DbHandler, infoLog, errLog *log.Logger) *App {
	application := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			errLog.Printf("%v", err)
			return c.Status(http.StatusInternalServerError).JSON(err)
		},
	})

	result := &App{
		web: application,
		// db:         DbHandler,
		infoLog: infoLog,
		errLog:  errLog,
	}

	// TODO роуты

	return result
}

// Run - запуск приложения.
//
// Принимает: адрес.
func (app *App) Run(addr string) {
	gracefully := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := app.web.Shutdown(); err != nil {
			app.errLog.Printf("Error while shutting down the server: %v", err)
		} else {
			app.infoLog.Printf("App closed gracefully\n")
		}

		close(gracefully)
	}()

	app.infoLog.Printf("App started on adress\n")
	app.errLog.Fatalln(app.web.Listen(addr))
}