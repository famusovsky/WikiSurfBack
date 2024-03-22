package postgres

import (
	"errors"
	"sort"
	"time"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/jmoiron/sqlx"
)

type dbProcessor struct {
	db *sqlx.DB
}

var (
	errBeginTx    = errors.New("error while starting transaction")
	errCommitTx   = errors.New("error while committing transaction")
	errNotCreator = errors.New("user is not the tournament's creator")
)

// AddUser implements DbHandler.
func (d *dbProcessor) AddUser(user models.User) (int, error) {
	wrapErr := errors.New("error while inserting user to the database")

	tx, err := d.db.Begin()
	if err != nil {
		return 0, errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	var id int

	if err = tx.QueryRow(addUser, user.Name, user.Email, user.Password).Scan(&id); err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	if err = tx.Commit(); err != nil {
		return 0, errors.Join(wrapErr, errCommitTx, err)
	}

	return id, nil
}

// AddRoute implements DbHandler.
func (d *dbProcessor) AddRoute(route models.Route) (int, error) {
	wrapErr := errors.New("error while inserting route to the database")

	tx, err := d.db.Begin()
	if err != nil {
		return 0, errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	var id int

	if err = tx.QueryRow(addRoute, route.Start, route.Finish, route.CreatorId).Scan(&id); err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	if err = tx.Commit(); err != nil {
		return 0, errors.Join(wrapErr, errCommitTx, err)
	}

	return id, nil
}

// AddSprint implements DbHandler.
func (d *dbProcessor) AddSprint(sprint models.Sprint) (int, error) {
	wrapErr := errors.New("error while inserting sprint to the database")
	tx, err := d.db.Begin()
	if err != nil {
		return 0, errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	var id int

	if err = tx.QueryRow(addSprint, sprint.StartTime, sprint.LengthTime, sprint.Success,
		sprint.RouteId, sprint.UserId, sprint.TournamentId, sprint.Path).Scan(&id); err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	if err = tx.Commit(); err != nil {
		return 0, errors.Join(wrapErr, errCommitTx, err)
	}

	return id, nil
}

// AddTournament implements DbHandler.
func (d *dbProcessor) AddTournament(tour models.Tournament, userId int) (int, error) {
	wrapErr := errors.New("error while inserting tournament to the database")

	tx, err := d.db.Begin()
	if err != nil {
		return 0, errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	var id int

	if err := tx.QueryRow(addTour, tour.StartTime, tour.EndTime, tour.Pswd, tour.Private).Scan(&id); err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	if _, err := tx.Exec(addCreatorToTour, id, userId); err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Join(wrapErr, errCommitTx, err)
	}

	return id, nil
}

// AddRouteToTour implements DbHandler.
func (d *dbProcessor) AddRouteToTour(tr models.TRRelation, userId int) error {
	wrapErr := errors.New("error while adding route to the tournament in the database")

	ok, err := d.CheckTournamentCreator(tr.TournamentId, userId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Join(wrapErr, errNotCreator)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(addRouteToTour, tr.TournamentId, tr.RouteId); err != nil {
		return errors.Join(wrapErr, err)
	}

	if err := tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// AddCreatorToTour implements DbHandler.
func (d *dbProcessor) AddCreatorToTour(tu models.TURelation, userId int) error {
	wrapErr := errors.New("error while adding creator to the tournament in the database")

	ok, err := d.CheckTournamentCreator(tu.TournamentId, userId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Join(wrapErr, errNotCreator)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(errBeginTx, err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(addCreatorToTour, tu.TournamentId, tu.UserId); err != nil {
		return errors.Join(wrapErr, err)
	}

	if err := tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// AddUserToTour implements DbHandler.
func (d *dbProcessor) AddUserToTour(tourId, userId int) error {
	wrapErr := errors.New("error while adding user to the tournament in the database")

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(addUserToTour, tourId, userId); err != nil {
		return errors.Join(wrapErr, err)
	}

	if err := tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// CheckTournamentCreator implements DbHandler.
func (d *dbProcessor) CheckTournamentCreator(tourId, userId int) (bool, error) {
	var cnt int

	if err := d.db.Get(&cnt, checkTournamentCreator, tourId, userId); err != nil {
		return false, errors.Join(errors.New("error while checking tournament's creator in the database"), err)
	}

	return cnt > 0, nil
}

func (d *dbProcessor) CheckTournamentParticipator(tourId int, userId int) (bool, error) {
	var cnt int

	if err := d.db.Get(&cnt, checkTournamentParticipator, tourId, userId); err != nil {
		return false, errors.Join(errors.New("error while checking tournament's creator in the database"), err)
	}

	return cnt > 0, nil
}

// CheckTournamentPassword implements DbHandler.
func (d *dbProcessor) CheckTournamentPassword(pswd string) (int, error) {
	var id int

	if err := d.db.Get(&id, checkTournamentPassword, pswd); err != nil {
		return 0, errors.Join(errors.New("error while checking tournament's password in the database"), err)
	}

	return id, nil
}

// GetCreatorTournaments implements DbHandler.
func (d *dbProcessor) GetCreatorTournaments(user int) ([]models.Tournament, error) {
	var res []models.Tournament

	if err := d.db.Select(&res, getCreatorTournaments, user); err != nil {
		return []models.Tournament{}, errors.Join(errors.New("error while getting all user created tournaments from the database"), err)
	}

	return res, nil
}

// GetOpenTournaments implements DbHandler.
func (d *dbProcessor) GetOpenTournaments() ([]models.Tournament, error) {
	var res []models.Tournament

	if err := d.db.Select(&res, getOpenTournaments, time.Now()); err != nil {
		return []models.Tournament{}, errors.Join(errors.New("error while getting opened tournaments from the database"), err)
	}

	return res, nil
}

// GetTournament implements DbHandler.
func (d *dbProcessor) GetTournament(tour int) (models.Tournament, error) {
	wrapErr := errors.New("error while getting tournament from the database")
	var res models.Tournament

	if err := d.db.Get(&res, getTournament, tour); err != nil {
		return models.Tournament{}, errors.Join(wrapErr, err)
	}

	return res, nil
}

func (d *dbProcessor) GetRoute(routeId int) (models.Route, error) {
	wrapErr := errors.New("error while getting route from the database")
	var route models.Route

	if err := d.db.Get(&route, getRoute, routeId); err != nil {
		return models.Route{}, errors.Join(wrapErr, err)
	}

	return route, nil
}

func (d *dbProcessor) GetRouteByCreds(start, finish string) (models.Route, error) {
	wrapErr := errors.New("error while getting route from the database")
	var route models.Route

	if err := d.db.Get(&route, getRouteByCreds, start, finish); err != nil {
		return models.Route{}, errors.Join(wrapErr, err)
	}

	return route, nil
}

// GetTournamentRoutes implements DbHandler.
func (d *dbProcessor) GetTournamentRoutes(tour int) ([]models.Route, error) {
	wrapErr := errors.New("error while getting tournament routes from the database")
	var routes []models.Route

	if err := d.db.Select(&routes, getTournamentRoutes, tour); err != nil {
		return []models.Route{}, errors.Join(wrapErr, err)
	}

	return routes, nil
}

// GetTournamentCreators implements DbHandler.
func (d *dbProcessor) GetTournamentCreators(tour int) ([]models.User, error) {
	wrapErr := errors.New("error while getting tournament creators from the database")
	var creators []models.User

	if err := d.db.Select(&creators, getTournamentCreators, tour); err != nil {
		return []models.User{}, errors.Join(wrapErr, err)
	}

	return creators, nil
}

// GetRouteRatings implements DbHandler.
func (d *dbProcessor) GetRouteRatings(routeId int) ([]models.RouteRating, error) {
	wrapErr := errors.New("error while getting route ratings from the database")
	var ratings []models.RouteRating

	if err := d.db.Select(&ratings, getRouteBest, routeId); err != nil {
		return []models.RouteRating{}, errors.Join(wrapErr, err)
	}

	return ratings, nil
}

// GetTournamentRatings implements DbHandler.
func (d *dbProcessor) GetTournamentRatings(tour int) ([]models.TourRating, error) {
	wrapErr := errors.New("error while getting tournament ratings from the database")

	var routes []int

	if err := d.db.Select(&routes, getTournamentRoutes, tour); err != nil {
		return []models.TourRating{}, errors.Join(wrapErr, err)
	}

	users := map[int]int{}
	for i := 0; i < len(routes); i++ {
		var rr []models.RouteRating

		if err := d.db.Select(&rr, getRouteTourBest, routes[i], tour); err != nil {
			return []models.TourRating{}, errors.Join(wrapErr, err)
		}

		if len(rr) == 0 {
			continue
		}

		sort.Slice(rr, func(i, j int) bool {
			return rr[i].SprintLengthTime < rr[j].SprintLengthTime
		})

		min, id := rr[1].SprintLengthTime, 1
		for j := 1; j < len(rr); j++ {
			if rr[j].SprintLengthTime < min {
				min = rr[j].SprintLengthTime
				id = rr[j].UserId
			}
		}

		users[id]++
	}

	ratings := make([]models.TourRating, 0, len(users))
	for id, points := range users {
		var name string

		if err := d.db.Get(&name, "SELECT name FROM users WHERE id = $1", id); err != nil {
			name = "Unknown Name - (try reload the window)"
		}
		ratings = append(ratings, models.TourRating{
			UserName: name,
			Points:   points,
		})
	}

	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i].Points < ratings[j].Points
	})

	return ratings, nil
}

// GetTournamentRatings implements DbHandler.
func (d *dbProcessor) GetRatings() ([]models.TourRating, error) {
	wrapErr := errors.New("error while getting ratings from the database")

	var routes []int

	if err := d.db.Select(&routes, getRoutes); err != nil {
		return []models.TourRating{}, errors.Join(wrapErr, err)
	}

	users := map[int]int{}
	for i := 0; i < len(routes); i++ {
		var rr []models.RouteRating

		if err := d.db.Select(&rr, getRouteBest, routes[i]); err != nil {
			return []models.TourRating{}, errors.Join(wrapErr, err)
		}

		if len(rr) == 0 {
			continue
		}
		sort.Slice(rr, func(i, j int) bool {
			return rr[i].SprintLengthTime < rr[j].SprintLengthTime
		})

		min, id := rr[1].SprintLengthTime, 1
		for j := 1; j < len(rr); j++ {
			if rr[j].SprintLengthTime < min {
				min = rr[j].SprintLengthTime
				id = rr[j].UserId
			}
		}

		users[id]++
	}

	ratings := make([]models.TourRating, 0, len(users))
	for id, points := range users {
		var name string

		if err := d.db.Get(&name, "SELECT name FROM users WHERE id = $1", id); err != nil {
			name = "Unknown Name - (try reload the window)"
		}
		ratings = append(ratings, models.TourRating{
			UserName: name,
			Points:   points,
		})
	}

	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i].Points < ratings[j].Points
	})

	return ratings, nil
}

// GetUser implements DbHandler.
func (d *dbProcessor) GetUser(email string) (models.User, error) {
	var user models.User

	if err := d.db.Get(&user, getUser, email); err != nil {
		return models.User{}, errors.Join(errors.New("error while getting user from the database"), err)
	}

	return user, nil
}

// GetUser implements DbHandler.
func (d *dbProcessor) GetUserById(id int) (models.User, error) {
	var user models.User

	if err := d.db.Get(&user, getUserById, id); err != nil {
		return models.User{}, errors.Join(errors.New("error while getting user from the database"), err)
	}

	return user, nil
}

func (d *dbProcessor) GetSprint(id int) (models.Sprint, error) {
	var sprint models.Sprint

	if err := d.db.Get(&sprint, getSprint, id); err != nil {
		return models.Sprint{}, errors.Join(errors.New("error while getting sprint from the database"), err)
	}

	return sprint, nil
}

// GetUserHistory implements DbHandler.
func (d *dbProcessor) GetUserHistory(id int) ([]models.Sprint, error) {
	var user []models.Sprint

	if err := d.db.Select(&user, getUserHistory, id); err != nil {
		return []models.Sprint{}, errors.Join(errors.New("error while getting user's history from the database"), err)
	}

	return user, nil
}

// GetUserRouteHistory implements DbHandler.
func (d *dbProcessor) GetUserRouteHistory(userId int, routeId int) ([]models.Sprint, error) {
	var user []models.Sprint

	if err := d.db.Select(&user, getUserRouteHistory, userId, routeId); err != nil {
		return []models.Sprint{}, errors.Join(errors.New("error while getting user's route history from the database"), err)
	}

	return user, nil
}

// GetUserTournaments implements DbHandler.
func (d *dbProcessor) GetUserTournaments(user int) ([]models.Tournament, error) {
	var tournaments []models.Tournament

	if err := d.db.Select(&tournaments, getUserTournaments, user); err != nil {
		return []models.Tournament{}, errors.Join(errors.New("error while getting tournaments in which user participates from the database"), err)
	}

	return tournaments, nil
}

// RemoveCreatorFromTour implements DbHandler.
func (d *dbProcessor) RemoveCreatorFromTour(tu models.TURelation, userId int) error {
	wrapErr := errors.New("error while removing creator from the tournament in the database")

	ok, err := d.CheckTournamentCreator(tu.TournamentId, userId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Join(wrapErr, errNotCreator)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(removeCreatorsFromTour, tu.TournamentId, tu.UserId); err != nil {
		return errors.Join(wrapErr, err)
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// RemoveRouteFromTour implements DbHandler.
func (d *dbProcessor) RemoveRouteFromTour(tr models.TRRelation, userId int) error {
	wrapErr := errors.New("error while removing route from the tournament in the database")

	ok, err := d.CheckTournamentCreator(tr.TournamentId, userId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Join(wrapErr, errNotCreator)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(removeRouteFromTour, tr.TournamentId, tr.RouteId); err != nil {
		return errors.Join(wrapErr, err)
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// RemoveUserFromTour implements DbHandler.
func (d *dbProcessor) RemoveUserFromTour(tourId, userId int) error {
	wrapErr := errors.New("error while removing route from the tournament in the database")

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(removeUserFromTour, tourId, userId); err != nil {
		return errors.Join(wrapErr, err)
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// UpdateTournament implements DbHandler.
func (d *dbProcessor) UpdateTournament(tour models.Tournament, user int) error {
	wrapErr := errors.New("error while updating the tournament in the database")

	ok, err := d.CheckTournamentCreator(tour.Id, user)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Join(wrapErr, errNotCreator)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(updateTournament, tour.Id, tour.StartTime, tour.EndTime, tour.Pswd, tour.Private); err != nil {
		return errors.Join(wrapErr, err)
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// UpdateTournament implements DbHandler.
func (d *dbProcessor) DeleteTournament(tourId, userId int) error {
	wrapErr := errors.New("error while updating the tournament in the database")

	ok, err := d.CheckTournamentCreator(tourId, userId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Join(wrapErr, errNotCreator)
	}

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(deleteTourFromRoutes, tourId); err != nil {
		return errors.Join(wrapErr, err)
	}
	if _, err = tx.Exec(deleteTourFromCreators, tourId); err != nil {
		return errors.Join(wrapErr, err)
	}
	if _, err = tx.Exec(deleteTourFromUsers, tourId); err != nil {
		return errors.Join(wrapErr, err)
	}
	if _, err = tx.Exec(deleteTournament, tourId); err != nil {
		return errors.Join(wrapErr, err)
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// UpdateUser implements DbHandler.
func (d *dbProcessor) UpdateUser(user models.User) error {
	wrapErr := errors.New("error while updating the user in the database")

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(updateUser, user.Id, user.Name, user.Email, user.Password); err != nil {
		return errors.Join(wrapErr, err)
	}

	if err = tx.Commit(); err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}
