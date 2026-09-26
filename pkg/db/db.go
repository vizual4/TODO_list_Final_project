package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const scheme = `CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(32) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT '');
	CREATE INDEX IF NOT EXISTS date_id ON scheduler (date);`

var DB *sql.DB

func Init(dbFile string) error {
	var err error

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	_, err = DB.Exec(scheme)
	if err != nil {
		return err
	}

	return nil
}
