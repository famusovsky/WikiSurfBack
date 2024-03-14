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
	// SQL запрос для создания таблицы маршрутов.
	createRoutes = `CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    start TEXT NOT NULL,
    finish TEXT NOT NULL,
    creator_id INTEGER NOT NULL,
    CONSTRAINT start_finish UNIQUE (start, finish),
    -- ratings JSON,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);`
	// SQL запрос для создания таблицы спринтов.
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
	// SQL запрос для создания таблицы соревнований.
	createTours = `CREATE TABLE IF NOT EXISTS tournaments (
    id SERIAL PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    pswd TEXT NOT NULL,
    private BOOLEAN NOT NULL
);`
	// SQL запрос для создания таблицы отношений соревнований и пользователей.
	createTURelations = `CREATE TABLE IF NOT EXISTS tournament_users (
    tour_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (tour_id) REFERENCES tournaments(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (tour_id, user_id)
);`
	// SQL запрос для создания таблицы отношений соревнований и создателей.
	createTCRelations = `CREATE TABLE IF NOT EXISTS tournament_creators (
    tour_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (tour_id) REFERENCES tournaments(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (tour_id, user_id)
);`
	// SQL запрос для создания таблицы отношений соревнований и маршрутов.
	createTRRelations = `CREATE TABLE IF NOT EXISTS tournament_routes (
    tour_id INTEGER NOT NULL,
    route_id INTEGER NOT NULL,
    FOREIGN KEY (tour_id) REFERENCES tournaments(id),
    FOREIGN KEY (route_id) REFERENCES routes(id),
    PRIMARY KEY (tour_id, route_id)
);`
)

// SQL запросы для удаления таблиц.
const (
	// SQL запрос для удаления таблицы пользователей.
	dropUsers = `DROP TABLE IF EXISTS users;`
	// SQL запрос для удаления таблицы маршрутов.
	dropRoutes = `DROP TABLE IF EXISTS routes;`
	// SQL запрос для удаления таблицы спринтов.
	dropSprints = `DROP TABLE IF EXISTS sprints;`
	// SQL запрос для удаления таблицы соревнований.
	dropTours = `DROP TABLE IF EXISTS tournaments;`
	// SQL запрос для удаления таблицы отношений соревнований и пользователей.
	dropTURelations = `DROP TABLE IF EXISTS tournament_users;`
	// SQL запрос для удаления таблицы отношений соревнований и создателей.
	dropTCRelations = `DROP TABLE IF EXISTS tournament_creators;`
	// SQL запрос для удаления таблицы отношений соревнований и маршрутов.
	dropTRRelations = `DROP TABLE IF EXISTS tournament_routes;`
)

// SQL запросы для получения данных.
const (
	// SQL запрос для получения пользователя по user.Email.
	getUser = `SELECT * FROM users WHERE email = $1;`
	// SQL запрос для получения истории спринтов пользователя по user.Email.
	getUserHistory = `SELECT * FROM sprints WHERE user_id = $1;`
	// SQL запрос для получения истории спринтов пользователя по user.Email, route.Id.
	getUserRouteHistory = `SELECT * FROM sprints WHERE user_id = $1 AND route_id = $2;`
	// SQL запрос для получения открытых соревнований.
	getOpenTournaments = `SELECT * FROM tournaments WHERE private = false AND end_time < $1;`
	// SQL запрос для получения соревнований по user.Id.
	getUserTournaments = `SELECT * FROM tournaments WHERE id IN (
        SELECT tour_id FROM tournament_users WHERE user_id = $1
    );`
	// SQL запрос для получения соревнований по creator.Id.
	getCreatorTournaments = `SELECT * FROM tournaments WHERE id IN (
        SELECT tour_id FROM tournament_creators WHERE user_id = $1
    );`
	// SQL запрос для получения соревнований по tournament.Id.
	getTournamentRoutes = `SELECT * FROM routes WHERE id IN (
        SELECT route_id FROM tournament_routes WHERE tour_id = $1
    );`
	// SQL запрос для получения данных о лучших результатах спринтов (length_time, path, user_id, spring_id) в маршруте по route.Id.
	getRouteBest = `SELECT DISTINCT ON (s.user_id) 
    s.length_time AS length_time, s.path AS path, s.user_id AS user_id, s.id AS sprint_id
    FROM sprints s INNER JOIN (
      SELECT user_id, MIN(length_time) AS min_length_time
      FROM sprints WHERE success = true AND route_id = $1
      GROUP BY user_id
    ) AS min_times ON s.user_id = min_times.user_id AND s.length_time = min_times.min_length_time
    ORDER BY s.length_time;`
	// SQL запрос для получения данных о лучших результатах спринтов (length_time, path, user_id, spring_id) в соревовании и маршруте по route.Id, tour.Id.
	getRouteTourBest = `SELECT DISTINCT ON (s.user_id) s.length_time, s.path, s.user_id, s.id
    FROM sprints s INNER JOIN (
      SELECT user_id, MIN(length_time) AS min_length_time
      FROM sprints WHERE success = true AND route_id = $1 AND tour_id = $2
      GROUP BY user_id
    ) AS min_times ON s.user_id = min_times.user_id AND s.length_time = min_times.min_length_time
    ORDER BY s.length_time;`
)

// SQL запросы для добавления данных.
const (
	// SQL запрос для добавления пользователя по name, email, password.
	addUser = `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id;`
	// SQL запрос для добавления маршрута по start, finish, creator_id.
	addRoute = `INSERT INTO routes (start, finish, creator_id) VALUES ($1, $2, $3) RETURNING id;`
	// SQL запрос для добавления спринта по start_time, length_time, success, route_id, user_id, tour_id, path.
	addSprint = `INSERT INTO sprints (start_time, length_time, success, route_id, user_id, tour_id, path)
    VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	// SQL запрос для добавления соревнования по start_time, end_time, pswd, private.
	addTour = `INSERT INTO tournaments (start_time, end_time, pswd, private) VALUES ($1, $2, $3, $4) RETURNING id;`
	// SQL запрос для добавления маршрута в соревнование по tour_id, route_id.
	addRouteToTour = `INSERT INTO tournament_routes (tour_id, route_id) VALUES ($1, $2);`
	// SQL запрос для добавления пользователя в соревнование по tour_id, user_id.
	addUserToTour = `INSERT INTO tournament_users (tour_id, user_id) VALUES ($1, $2);`
	// SQL запрос для добавления создателя в соревнование по tour_id, user_id.
	addCreatorToTour = `INSERT INTO tournament_creators (tour_id, user_id) VALUES ($1, $2);`
)

// SQL запросы для удаления данных.
const (
	// SQL запрос для добавления маршрута из соревнования по tour_id, route_id.
	removeRouteFromTour = `DELETE FROM tournament_routes WHERE tour_id = $1 AND route_id = $2;`
	// SQL запрос для добавления пользователя из соревнования по tour_id, user_id.
	removeUserFromTour = `DELETE FROM tournament_users WHERE tour_id = $1 AND user_id = $2;`
	// SQL запрос для добавления создателя из соревнования по tour_id, user_id.
	removeCreatorsFromTour = `DELETE FROM tournament_creators WHERE tour_id = $1 AND user_id = $2;`
)

// SQL запросы для проверки данных.
const (
	// SQL запрос для проверки пароля соревнования по id, pswd.
	checkTournamentPassword = `SELECT COUNT(*) FROM tournaments WHERE id = $1 AND pswd = $2;`
	// SQL запрос для проверки создателя соревнования по tour_id, user_id.
	checkTournamentCreator = `SELECT COUNT(*) FROM tournaments t JOIN tournament_creators tc ON t.id = tc.tour_id WHERE tc.user_id = $2 AND tc.tour_id = $1;`
)

// SQL запросы для обновления данных.
const (
	// SQL запрос для обновления соревнования по id, start_time, end_time, pswd, private.
	updateTournament = `UPDATE tournaments SET start_time = $2, end_time = $3, pswd = $4, private = $5 WHERE id = $1;`
	// SQL запрос для обновления пользователя по id, name, email, password.
	updateUser = `UPDATE users SET name = $2, email = $3, password = $4 WHERE id = $1;`
)
