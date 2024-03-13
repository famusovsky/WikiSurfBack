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
	AddTournament(tour models.Tournament) (int, error)                // AddTournament - добавление нового соревнования в БД.
	AddRouteToTour(tr models.TRRelation, userId int) error            // AddRouteToTour - добавление маршрута в соревнование.
	RemoveRouteFromTour(tr models.TRRelation, userId int) error       // AddRouteToTour - удаление маршрута из соревнования.
	AddUserToTour(tu models.TURelation, userId int) error             // AddUserToTour - добавление участника в соревнование.
	RemoveUserFromTour(tu models.TURelation, userId int) error        // AddUserToTour - удаление участника из соревнования.
	AddCreatorToTour(tu models.TURelation, userId int) error          // AddCreatorToTour - добавление создателя в соревнование.
	RemoveCreatorFromTour(tu models.TURelation, userId int) error     // RemoveCreatorFromTour - удаление создателя из соревнования.
	GetUser(email, pswd string) (models.User, error)                  // GetUser - получение Id пользователя по email-у и паролю.
	GetUserHistory(id int) ([]models.Sprint, error)                   // GetUserHistory - получение истории спринтов пользователя.
	GetUserRouteHistory(userId, routeId int) ([]models.Sprint, error) // GetUserRouteHistory - получение истории спринтов пользователя по маршруту.
	GetRouteRatings(routeId int) ([]byte, error)                      // GetRouteRatings - получение рейтинга по маршруту в формате JSON.
	GetOpenTournaments() ([]models.Tournament, error)                 // GetOpenTournaments - получение списка соревнований, открытых для вступления.
	GetUserTournaments(user int) ([]models.Tournament, error)         // GetUserTournaments - получение списка соревнований, в которых пользователь участвует.
	GetCreatorTournaments(user int) ([]models.Tournament, error)      // GetCreatorTournaments - получение списка соревнований, в которых пользователь выступает создателем.
	GetTournamentRatings(tour int) ([]byte, error)                    // GetTournamentRatings - получение рейтинга по соревнованию в формате JSON.
	CheckTournamentPassword(tourId int, pswd string) (bool, error)    // CheckTournamentPassword - проверка на соответствие кода-пароля соревнования.
	CheckTournamentCreator(tourId, userId int) (bool, error)          // CheckTournamentCreator - проверка на соответствие Id пользователя с Id создателей соревнования.
	UpdateTournament(tour models.Tournament, user int) error          // UpdateTournament - обновление основных данных о соревновании.
	UpdateUser(user models.User) error                                // UpdateUser -  - обновление основных данных о пользователе.
}

// Get - функция, возвращающая объект, реализующий интерфейс DbHandler.
func Get(db *sql.DB, createTables bool) (DbHandler, error) {
	if createTables {
		err := overrideDB(db)
		if err != nil {
			return nil, err
		}
	}

	// err := checkDB(db)
	// if err != nil {
	// 	return nil, err
	// }

	return &dbProcessor{sqlx.NewDb(db, "postgres")}, nil
}
