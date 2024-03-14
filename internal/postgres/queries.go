package postgres

// SQL запросы для создания таблиц.
const (
	// SQL запрос для создания таблицы пользователей.
	createUsers string = `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);`

	createRoutes = `CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    start TEXT NOT NULL,
    finish TEXT NOT NULL,
    creator_id INTEGER NOT NULL,
    CONSTRAINT start_finish UNIQUE (start, finish),
    -- ratings JSON,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);`

	createSprints = `CREATE TABLE IF NOT EXISTS sprints (
    id SERIAL PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    length_time INTEGER NOT NULL,
    success BOOLEAN NOT NULL,
    route_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    tour_id INTEGER,
    path JSON NOT NULL,
    FOREIGN KEY (route_id) REFERENCES routes(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (tour_id) REFERENCES tournaments(id)
);`

	createTours = `CREATE TABLE IF NOT EXISTS tournaments (
    id SERIAL PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    pswd TEXT NOT NULL,
    private BOOLEAN NOT NULL
);`

	createTURelations = `CREATE TABLE IF NOT EXISTS tournament_users (
    tour_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (tour_id) REFERENCES tournaments(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (tour_id, user_id)
);`

	createTCRelations = `CREATE TABLE IF NOT EXISTS tournament_creators (
    tour_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (tour_id) REFERENCES tournaments(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (tour_id, user_id)
);`

	createTRRelations = `CREATE TABLE IF NOT EXISTS tournament_routes (
    tour_id INTEGER NOT NULL,
    route_id INTEGER NOT NULL,
    FOREIGN KEY (tour_id) REFERENCES tournaments(id),
    FOREIGN KEY (route_id) REFERENCES routes(id),
    PRIMARY KEY (tour_id, route_id)
);`
)

const (
	dropUsers       = `DROP TABLE IF EXISTS users;`
	dropRoutes      = `DROP TABLE IF EXISTS routes;`
	dropSprints     = `DROP TABLE IF EXISTS sprints;`
	dropTours       = `DROP TABLE IF EXISTS tournaments;`
	dropTURelations = `DROP TABLE IF EXISTS tournament_users;`
	dropTCRelations = `DROP TABLE IF EXISTS tournament_creators;`
	dropTRRelations = `DROP TABLE IF EXISTS tournament_routes;`
)

const (
	getUser             = `SELECT * FROM users WHERE email = $1;`
	getUserHistory      = `SELECT * FROM sprints WHERE user_id = $1;`
	getUserRouteHistory = `SELECT * FROM sprints WHERE user_id = $1 AND route_id = $2;`
	getOpenTournaments  = `SELECT * FROM tournaments WHERE private = false AND end_time < $1;`
	getUserTournaments  = `SELECT * FROM tournaments WHERE id IN (
        SELECT tour_id FROM tournament_users WHERE user_id = $1
    );`
	getCreatorTournaments = `SELECT * FROM tournaments WHERE id IN (
        SELECT tour_id FROM tournament_creators WHERE user_id = $1
    );`
	getTournamentRoutes = `SELECT * FROM routes WHERE id IN (
        SELECT route_id FROM tournament_routes WHERE route_id = $1
    );`
	getRouteBest = `SELECT DISTINCT ON (s.user_id) 
    s.length_time AS length_time, s.path AS path, s.user_id AS user_id, s.id AS sprint_id
    FROM sprints s INNER JOIN (
      SELECT user_id, MIN(length_time) AS min_length_time
      FROM sprints WHERE success = true AND route_id = $1
      GROUP BY user_id
    ) AS min_times ON s.user_id = min_times.user_id AND s.length_time = min_times.min_length_time
    ORDER BY s.length_time;`
	getRouteTourBest = `SELECT DISTINCT ON (s.user_id) s.length_time, s.path, s.user_id, s.id
    FROM sprints s INNER JOIN (
      SELECT user_id, MIN(length_time) AS min_length_time
      FROM sprints WHERE success = true AND route_id = $1 AND tour_id = $2
      GROUP BY user_id
    ) AS min_times ON s.user_id = min_times.user_id AND s.length_time = min_times.min_length_time
    ORDER BY s.length_time;`
)

const (
	addUser   = `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id;`
	addRoute  = `INSERT INTO routes (start, finish, creator_id) VALUES ($1, $2, $3) RETURNING id;`
	addSprint = `INSERT INTO sprints (start_time, length_time, success, route_id, user_id, tour_id, path)
    VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	addTour          = `INSERT INTO tournaments (start_time, end_time, pswd, private) VALUES ($1, $2, $3, $4) RETURNING id;`
	addRouteToTour   = `INSERT INTO tournament_routes (tour_id, route_id) VALUES ($1, $2);`
	addUserToTour    = `INSERT INTO tournament_users (tour_id, user_id) VALUES ($1, $2);`
	addCreatorToTour = `INSERT INTO tournament_creators (tour_id, user_id) VALUES ($1, $2);`
)

const (
	removeRouteFromTour    = `DELETE FROM tournament_routes WHERE tour_id = $1 AND route_id = $2;`
	removeUserFromTour     = `DELETE FROM tournament_users WHERE tour_id = $1 AND user_id = $2;`
	removeCreatorsFromTour = `DELETE FROM tournament_creators WHERE tour_id = $1 AND user_id = $2;`
)

const (
	checkTournamentPassword = `SELECT COUNT(*) FROM tournaments WHERE id = $1 AND pswd = $2;`
	checkTournamentCreator  = `SELECT COUNT(*) FROM tournaments t JOIN tournament_creators tc ON t.id = tc.tour_id WHERE tc.user_id = $1;`
)

const (
	updateTournament = `UPDATE tournaments SET start_time = $2, end_time = $3, pswd = $4, private = $5 WHERE id = $1;`
	updateUser       = `UPDATE users SET name = $2, email = $3, password = $4 WHERE id = $1;`
)
