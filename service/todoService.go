package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"todoGoApi/db"
	"todoGoApi/types"
)

func GetTodoList(userId int) ([]types.Todo, error) {
	todoDb, err := db.DB_CONN.Query(context.Background(), "select * from get_todos($1)", userId)
	if err != nil {
		fmt.Println("DB ERROR")
		return []types.Todo{}, err
	}

	todoList := []types.Todo{}

	for todoDb.Next() {
		var todo types.Todo
		err := todoDb.Scan(&todo.Id, &todo.Name, &todo.Completed)
		if err != nil {
			fmt.Println("ROW ERRO BITCH")
			return []types.Todo{}, err
		}
		todoList = append(todoList, todo)
	}
	return todoList, nil
}

func AddTodo(todo types.TodoSchema, userId int) (types.Todo, error) {

	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	var res int
	err := db.DB_CONN.QueryRow(
		context.Background(),
		"CALL insert_todo($1,$2,$3,$4,$5)", userId, todo.Name, todo.Completed, dbDate, res).Scan(&res)

	if err != nil {
		return types.Todo{}, err
	}

	log.Println(res, todo.Name, todo.Completed)
	return types.Todo{
		Id:        res,
		Name:      todo.Name,
		Completed: todo.Completed,
	}, nil
}

func DeleteTodo(id int, userId int) error {
	_, err := db.DB_CONN.Exec(context.Background(), "CALL delete_todo($1,$2)", id, userId)

	if err != nil {
		if strings.Contains(err.Error(), "Invalid id") {
			return errors.New("invalid id")
		}
		return err
	}
	return nil
}

func UpdateTodo(userId int, id int, name string, status bool) (types.Todo, error) {
	_, err := db.DB_CONN.Exec(context.Background(), "CALL update_todo($1,$2,$3,$4)", id, name, status, userId)

	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "Invalid id") {
			return types.Todo{}, errors.New("invalid id")
		}
		return types.Todo{}, err
	}

	return types.Todo{
		Id:        id,
		Name:      name,
		Completed: status,
	}, err
}
