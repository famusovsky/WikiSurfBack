package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// overrideDB - функция, перезаписывающая таблицы WikiSurf в БД.
func overrideDB(db *sql.DB) error {
	err := dropTables(db)
	if err != nil {
		return err
	}
	err = createTables(db)
	return err
}

// dropTables - функция, удаляющая таблицы WikiSurf в БД.
func dropTables(db *sql.DB) error {
	q := strings.Join([]string{
		dropUsers,
		dropRoutes,
		dropSprints,
		dropTours,
		dropTURelations,
		dropTCRelations,
		dropTRRelations,
	}, " ")

	_, err := db.Exec(q)
	if err != nil {
		return errors.Join(fmt.Errorf("error while dropping tables: %s", err))
	}

	return nil
}

// createTables - функция, добавляющая таблицы WikiSurf в БД.
func createTables(db *sql.DB) error {
	q := strings.Join([]string{
		createUsers,
		createRoutes,
		createSprints,
		createTours,
		createTURelations,
		createTCRelations,
		createTRRelations,
	}, " ")

	_, err := db.Exec(q)
	if err != nil {
		return errors.Join(fmt.Errorf("error while creating tables: %s", err))
	}

	return nil
}
