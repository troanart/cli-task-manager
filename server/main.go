package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID      int
	Title   string
	Content string
}

func main() {
	//  открываем или создеём файл бд (task.db) если его нету
	db, err := sql.Open("sqlite3", "./task.db")
	if err != nil {
		log.Fatal("Не удалось подключиться к DB:", err)
	}
	defer db.Close()

	err = initDB(db)
	if err != nil {
		log.Fatal("Ошибка шаблона создания таблицы: ", err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Использование:")
		fmt.Println("  go run . list")
		fmt.Println("  go run . add \"заголовок\" \"описание\"")
		fmt.Println("  go run . delete 3")
		return
	}

	command := os.Args[1]

	switch command {
	case "list":
		allTasks, err := getAllTasks(db)
		if err != nil {
			log.Fatal("ошибка чтения:", err)
		}

		for _, task := range allTasks {
			fmt.Printf("%d: %s %s\n", task.ID, task.Title, task.Content)
		}
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("нужно: go run . add \"заголовок\" \"описание\"")
			return
		}

		title := os.Args[2]
		content := os.Args[3]

		err := createTask(db, title, content)
		if err != nil {
			log.Fatal("ошибка создания:", err)
		}
		fmt.Println("задача добавлена")

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("нужно: go run . delete \"id\"")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("id должен быть числом:", err)
		}

		err = deleteTask(db, id)
		if err != nil {
			log.Fatal("ошибка удаления:", err)
		}
		fmt.Println("задача удалена")

	default:
		fmt.Println("Неизвестная команда:", os.Args[1])
	}

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
