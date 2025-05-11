package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Conn *sql.DB

// Init открывает соединение с SQLite и создаёт все нужные таблицы
func Init(path string) error {
	var err error
	Conn, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	schema := `
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      login TEXT UNIQUE NOT NULL,
      password TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS expressions (
      id TEXT PRIMARY KEY,
      user_id INTEGER NOT NULL,
      expression TEXT NOT NULL,
      status TEXT NOT NULL,
      result REAL,
      FOREIGN KEY(user_id) REFERENCES users(id)
    );
    `
	_, err = Conn.Exec(schema)
	return err
}
