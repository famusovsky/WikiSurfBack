package postgres

import (
	"database/sql"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/jmoiron/sqlx"
)

// DbHandler - интерфейс, описывающий взаимодействие с БД WikiSurf.
type DbHandler interface {
	AddUser(user models.User) (int, error)                            // AddUser - добавление нового пользователя в БД.
	AddRoute(route models.Route) (int, error)                         // AddRoute - добавление нового маршрута в БД.
	AddSprint(sprint models.Sprint) (int, error)                      // AddSprint - добавление нового спринта в БД.
	AddTournament(tour models.Tournament, userId int) (int, error)    // AddTournament - добавление нового соревнования в БД.
	AddRouteToTour(tr models.TRRelation, userId int) error            // AddRouteToTour - добавление маршрута в соревнование.
	RemoveRouteFromTour(tr models.TRRelation, userId int) error       // AddRouteToTour - удаление маршрута из соревнования.
	AddUserToTour(tourId, userId int) error                           // AddUserToTour - добавление участника в соревнование.
	RemoveUserFromTour(tourId, userId int) error                      // RemoveUserFromTour - удаление участника из соревнования.
	AddCreatorToTour(tu models.TURelation, userId int) error          // AddCreatorToTour - добавление создателя в соревнование.
	RemoveCreatorFromTour(tu models.TURelation, userId int) error     // RemoveCreatorFromTour - удаление создателя из соревнования.
	GetUser(email string) (models.User, error)                        // GetUser - получение пользователя по email-у
	GetUserById(id int) (models.User, error)                          // GetUserById - получение пользователя по id
	GetRoute(id int) (models.Route, error)                            // GetRoute - получение маршрута по id.
	GetPopularRoutes() ([]models.Route, error)                        // GetRoutes - получение популярных маршрутов.
	GetRouteByCreds(start, finish string) (models.Route, error)       // GetRouteByCreds - получение маршрута по start, finish.
	GetSprint(id int) (models.Sprint, error)                          // GetSprint - получение спринта по id.
	GetTournament(id int) (models.Tournament, error)                  // GetTournament - получение соревнования по id.
	GetTournamentRoutes(id int) ([]models.Route, error)               // GetTournamentRoutes - получение маршрутов соревнования.
	GetTournamentCreators(id int) ([]models.User, error)              // GetTournamentRoutes - получение маршрутов соревнования.
	GetUserHistory(id int) ([]models.Sprint, error)                   // GetUserHistory - получение истории спринтов пользователя.
	GetUserRouteHistory(userId, routeId int) ([]models.Sprint, error) // GetUserRouteHistory - получение истории спринтов пользователя по маршруту.
	GetRouteRatings(routeId int) ([]models.RouteRating, error)        // GetRouteRatings - получение рейтинга по маршруту.
	GetOpenTournaments() ([]models.Tournament, error)                 // GetOpenTournaments - получение списка соревнований, открытых для вступления.
	GetUserTournaments(user int) ([]models.Tournament, error)         // GetUserTournaments - получение списка соревнований, в которых пользователь участвует.
	GetCreatorTournaments(user int) ([]models.Tournament, error)      // GetCreatorTournaments - получение списка соревнований, в которых пользователь выступает создателем.
	GetTournamentRatings(tour int) ([]models.TourRating, error)       // GetTournamentRatings - получение рейтинга по соревнованию .
	GetRatings() ([]models.TourRating, error)                         // GetRatings - получение общего рейтинга.
	CheckTournamentPassword(pswd string) (int, error)                 // CheckTournamentPassword - проверка на соответствие кода-пароля соревнования.
	CheckTournamentCreator(tourId, userId int) (bool, error)          // CheckTournamentCreator - проверка на соответствие Id пользователя с Id создателей соревнования.
	CheckTournamentParticipator(tourId, userId int) (bool, error)     // CheckTournamentParticipator - проверка на соответствие Id пользователя с Id участников соревнования.
	UpdateTournament(tour models.Tournament, user int) error          // UpdateTournament - обновление основных данных о соревновании.
	DeleteTournament(tourId, userId int) error                        // DeleteTournament - удаление данных о соревновании.
	UpdateUser(user models.User) error                                // UpdateUser - обновление основных данных о пользователе.
}

// Get - функция, возвращающая объект, реализующий интерфейс DbHandler.
func Get(db *sql.DB, override bool) (DbHandler, error) {
	if override {
		if err := overrideDB(db); err != nil {
			return nil, err
		}
	} else {
		if err := createTables(db); err != nil {
			return nil, err
		}
	}

	// err := checkDB(db)
	// if err != nil {
	// 	return nil, err
	// }

	return &dbProcessor{sqlx.NewDb(db, "postgres")}, nil
}
