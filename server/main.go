package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

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
		fmt.Println("  go run . edit 3 \"заголовок\" \"описание\"")
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

	case "edit":
		if len(os.Args) < 5 {
			fmt.Println("нужно: go run . edit \"id\" \"title\" \"content\"")
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("id должен быть числом ", err)
		}

		title := os.Args[3]
		content := os.Args[4]

		err = editTask(db, id, title, content)
		if err != nil {
			log.Fatal("ошибка изменения", err)
		}
		fmt.Println("задача успешно изменена")

	default:
		fmt.Println("Неизвестная команда:", os.Args[1])
	}

}
