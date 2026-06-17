package main

import "database/sql"

type Task struct {
	ID      int
	Title   string
	Content string
}

func getAllTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, title, content FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Content)
		if err != nil {

			return nil, err
		}

		tasks = append(tasks, t)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil

}

func createTask(db *sql.DB, title, content string) error {

	// готовим SQL-шаблон
	statement, err := db.Prepare("INSERT INTO tasks (title, content) VALUES (?, ?)")
	if err != nil {
		return err
	}

	defer statement.Close()

	// выполняем запрос
	_, err = statement.Exec(title, content)
	if err != nil {
		return err
	}

	return nil

}

func deleteTask(db *sql.DB, id int) error {
	statement, err := db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func editTask(db *sql.DB, id int, title, content string) error {
	statement, err := db.Prepare("UPDATE tasks SET title = ? , content = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(title, content, id)
	if err != nil {
		return err
	}

	return nil
}
