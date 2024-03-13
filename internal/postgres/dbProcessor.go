package postgres

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/famusovsky/WikiSurfBack/internal/models"
	"github.com/jmoiron/sqlx"
)

// TODO wrap all errors

type dbProcessor struct {
	db *sqlx.DB
}

var (
	errBeginTx     = errors.New("error while starting transaction")
	errCommitTx    = errors.New("error while committing transaction")
	errNotCreator  = errors.New("user is not the tournament's creator")
	errMarshalling = errors.New("error while marshalling route ratings to json")
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
	err = tx.QueryRow(addUser, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
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
	err = tx.QueryRow(addRoute, route.Start, route.Finish, route.CreatorId).Scan(&id)
	if err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
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
	err = tx.QueryRow(addSprint, sprint.StartTime, sprint.LengthTime, sprint.Success,
		sprint.RouteId, sprint.UserId, sprint.TournamentId, sprint.Path).Scan(&id)
	if err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, errors.Join(wrapErr, errCommitTx, err)
	}

	return id, nil
}

// AddTournament implements DbHandler.
func (d *dbProcessor) AddTournament(tour models.Tournament) (int, error) {
	wrapErr := errors.New("error while inserting tournament to the database")

	tx, err := d.db.Begin()
	if err != nil {
		return 0, errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRow(addTour, tour.StartTime, tour.EndTime, tour.Pswd, tour.Private).Scan(&id)
	if err != nil {
		return 0, errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
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

	_, err = tx.Exec(addRouteToTour, tr.TournamentId, tr.RouteId)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
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

	_, err = tx.Exec(addCreatorToTour, tu.TournamentId, tu.UserId)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// AddUserToTour implements DbHandler.
func (d *dbProcessor) AddUserToTour(tu models.TURelation, userId int) error {
	wrapErr := errors.New("error while adding user to the tournament in the database")

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(addUserToTour, tu.TournamentId, tu.UserId)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// CheckTournamentCreator implements DbHandler.
func (d *dbProcessor) CheckTournamentCreator(tourId int, userId int) (bool, error) {
	var cnt int
	err := d.db.Get(&cnt, checkTournamentCreator, tourId, userId)
	if err != nil {
		return false, errors.Join(errors.New("error while checking tournament's creator in the database"), err)
	}

	return cnt > 0, nil
}

// CheckTournamentPassword implements DbHandler.
func (d *dbProcessor) CheckTournamentPassword(tourId int, pswd string) (bool, error) {
	var cnt int
	err := d.db.Get(&cnt, checkTournamentPassword, tourId, pswd)
	if err != nil {
		return false, errors.Join(errors.New("error while checking tournament's password in the database"), err)
	}

	return cnt > 0, nil
}

// GetCreatorTournaments implements DbHandler.
func (d *dbProcessor) GetCreatorTournaments(user int) ([]models.Tournament, error) {
	var res []models.Tournament
	err := d.db.Select(&res, getCreatorTournaments, user)
	if err != nil {
		return nil, errors.Join(errors.New("error while getting all user created tournaments from the database"), err)
	}

	return res, nil
}

// GetOpenTournaments implements DbHandler.
func (d *dbProcessor) GetOpenTournaments() ([]models.Tournament, error) {
	var res []models.Tournament
	err := d.db.Select(&res, getOpenTournaments)
	if err != nil {
		return nil, errors.Join(errors.New("error while getting opened tournaments from the database"), err)
	}

	return res, nil
}

// GetRouteRatings implements DbHandler.
func (d *dbProcessor) GetRouteRatings(routeId int) ([]byte, error) {
	wrapErr := errors.New("error while getting route ratings from the database")
	var ratings []models.RouteRating
	err := d.db.Select(&ratings, getRouteBest, routeId)
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	res, err := json.Marshal(ratings)
	if err != nil {
		return nil, errors.Join(wrapErr, errMarshalling, err)
	}

	return res, nil
}

// GetTournamentRatings implements DbHandler.
func (d *dbProcessor) GetTournamentRatings(tour int) ([]byte, error) {
	wrapErr := errors.New("error while getting tournament ratings from the database")

	var routes []int
	err := d.db.Select(&routes, getTournamentRoutes, tour)
	if err != nil {
		return nil, errors.Join(wrapErr, err)
	}

	users := map[int]int{}
	for i := 0; i < len(routes); i++ {
		var rr []models.RouteRating
		err = d.db.Select(&rr, getRouteTourBest, routes[i], tour)
		if err != nil {
			return nil, errors.Join(wrapErr, err)
		}

		if len(rr) == 0 {
			continue
		}

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
		err := d.db.Get(&name, "SELECT name FROM users WHERE id = $1", id)
		if err != nil {
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

	res, err := json.Marshal(ratings)
	if err != nil {
		return nil, errors.Join(wrapErr, errMarshalling, err)
	}

	return res, nil
}

// GetUser implements DbHandler.
func (d *dbProcessor) GetUser(email string, pswd string) (models.User, error) {
	var user models.User
	err := d.db.Get(&user, getUser, email, pswd)
	if err != nil {
		return models.User{}, errors.Join(errors.New("error while getting user from the database"), err)
	}

	return user, nil
}

// GetUserHistory implements DbHandler.
func (d *dbProcessor) GetUserHistory(id int) ([]models.Sprint, error) {
	var user []models.Sprint
	err := d.db.Select(&user, getUserHistory, id)
	if err != nil {
		return nil, errors.Join(errors.New("error while getting user's history from the database"), err)
	}

	return user, nil
}

// GetUserRouteHistory implements DbHandler.
func (d *dbProcessor) GetUserRouteHistory(userId int, routeId int) ([]models.Sprint, error) {
	var user []models.Sprint
	err := d.db.Select(&user, getUserRouteHistory, userId, routeId)
	if err != nil {
		return nil, errors.Join(errors.New("error while getting user's route history from the database"), err)
	}

	return user, nil
}

// GetUserTournaments implements DbHandler.
func (d *dbProcessor) GetUserTournaments(user int) ([]models.Tournament, error) {
	var tournaments []models.Tournament
	err := d.db.Select(&tournaments, getUserTournaments, user)
	if err != nil {
		return nil, errors.Join(errors.New("error while getting tournaments in which user participates from the database"), err)
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

	_, err = tx.Exec(removeCreatorsFromTour, tu.TournamentId, tu.UserId)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
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

	_, err = tx.Exec(removeRouteFromTour, tr.TournamentId, tr.RouteId)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}

// RemoveUserFromTour implements DbHandler.
func (d *dbProcessor) RemoveUserFromTour(tu models.TURelation, userId int) error {
	wrapErr := errors.New("error while removing route from the tournament in the database")

	tx, err := d.db.Begin()
	if err != nil {
		return errors.Join(wrapErr, errBeginTx, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(removeUserFromTour, tu.TournamentId, tu.UserId)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
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
	_, err = tx.Exec(updateTournament, tour.Id, tour.StartTime, tour.EndTime, tour.Pswd, tour.Private)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
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
	_, err = tx.Exec(updateUser, user.Id, user.Name, user.Email, user.Password)
	if err != nil {
		return errors.Join(wrapErr, err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(wrapErr, errCommitTx, err)
	}

	return nil
}
