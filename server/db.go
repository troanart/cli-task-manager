package main

import "database/sql"

func initDB(db *sql.DB) error {
	// готовим SQL-шаблон
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, title TEXT, content TEXT )")
	if err != nil {

		return err
	}

	defer statement.Close()

	// выполняем запрос
	_, err = statement.Exec()
	if err != nil {

		return err
	}

	return nil
}
