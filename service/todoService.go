package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"todoGoApi/common"
	"todoGoApi/db"
	"todoGoApi/types"
)

func GetTodoList(userId int) ([]types.Todo, int, error) {
	todoDb, err := db.DB_CONN.Query(context.Background(), "select * from get_todos($1)", userId)
	if err != nil {
		common.Logger(err.Error(), "GetTodoList")
		return []types.Todo{}, http.StatusInternalServerError, errors.New("something went wrong")
	}

	todoList := []types.Todo{}

	for todoDb.Next() {
		var todo types.Todo
		err := todoDb.Scan(&todo.Id, &todo.Name, &todo.Completed)
		if err != nil {
			common.Logger(err.Error(),"GetTodoList")
			return []types.Todo{}, http.StatusInternalServerError, errors.New("something went wrong")
		}
		todoList = append(todoList, todo)
	}
	return todoList, http.StatusOK, nil
}

func AddTodo(todo types.TodoSchema, userId int) (types.Todo, int, error) {

	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	var res int
	err := db.DB_CONN.QueryRow(
		context.Background(),
		"CALL insert_todo($1,$2,$3,$4,$5)", userId, todo.Name, todo.Completed, dbDate, res).Scan(&res)

	if err != nil {
		common.Logger(err.Error(),"GetTodoList")
		return types.Todo{}, http.StatusInternalServerError, errors.New("something went wrong")
	}

	return types.Todo{
		Id:        res,
		Name:      todo.Name,
		Completed: todo.Completed,
	}, http.StatusCreated, nil
}

func DeleteTodo(id int, userId int) (int, error) {
	_, err := db.DB_CONN.Exec(context.Background(), "CALL delete_todo($1,$2)", id, userId)

	if err != nil {
		common.Logger(err.Error(),"GetTodoList")
		if strings.Contains(err.Error(), "Invalid id") {
			return http.StatusBadRequest,errors.New("invalid id")
		}
		return http.StatusInternalServerError, errors.New("something went wrong")
	}
	return http.StatusNoContent, nil
}

func UpdateTodo(userId int, id int, name string, status bool) (types.Todo, int, error) {
	_, err := db.DB_CONN.Exec(context.Background(), "CALL update_todo($1,$2,$3,$4)", id, name, status, userId)

	if err != nil {
		common.Logger(err.Error(),"GetTodoList")
		if strings.Contains(err.Error(), "Invalid id") {
			return types.Todo{}, http.StatusBadRequest, errors.New("invalid id")
		}
		return types.Todo{},http.StatusInternalServerError, errors.New("something went wrong")
	}

	return types.Todo{
		Id:        id,
		Name:      name,
		Completed: status,
	}, http.StatusAccepted, err
}
