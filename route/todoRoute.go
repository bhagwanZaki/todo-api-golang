package route

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todoGoApi/common"
	"todoGoApi/service"
	"todoGoApi/types"
)

type TodoApi struct{}

func (h *TodoApi) HealthCheckAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	hashedPassword, hashErr := common.HashPassword("adminadmin")

	if hashErr != nil {
		log.Println(hashErr.Error())
	}
	log.Println(hashedPassword)

	res := map[string]string{
		"message": "API is up",
	}
	json.NewEncoder(w).Encode(res)
}

func (h *TodoApi) GetTodos(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")

	todoList, statusCode, err := service.GetTodoList(userData.Id)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "GetTodos")
		return
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(todoList)

}

func (h *TodoApi) AddTodo(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")
	var data types.TodoSchema
	bodyErr := json.NewDecoder(r.Body).Decode(&data)

	if r.Body == nil {
		common.ErrorResponse(w, "Empty request body", http.StatusBadRequest, "AddTodo")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError, "AddTodo")
		return
	}

	if data.Name == "" {
		common.ErrorResponse(w, "Name field can't be empty", http.StatusBadRequest, "AddTodo")
		return
	}

	newTodo, statusCode,err := service.AddTodo(data, userData.Id)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "AddTodo")
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(newTodo)
}

func (h *TodoApi) DeleteTodo(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")

	id, parseErr := strconv.Atoi(r.PathValue("id"))
	if parseErr != nil {
		common.ErrorResponse(w, "Invalid id", http.StatusInternalServerError, "DeleteTodo")
		return
	}

	statusCode, err := service.DeleteTodo(id, userData.Id)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "DeleteTodo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TodoApi) UpdateTodo(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")

	id, parseErr := strconv.Atoi(r.PathValue("id"))
	if parseErr != nil {
		common.ErrorResponse(w, "Invalid id", http.StatusInternalServerError, "UpdateTodo")
		return
	}

	var data types.TodoSchema
	bodyErr := json.NewDecoder(r.Body).Decode(&data)

	if r.Body == nil {
		common.ErrorResponse(w, "Empty request body", http.StatusBadRequest, "UpdateTodo")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError, "UpdateTodo")
		return
	}

	if data.Name == "" {
		common.ErrorResponse(w, "Name field can't be empty", http.StatusBadRequest, "UpdateTodo")
		return
	}

	updatedTodo, statusCode,err := service.UpdateTodo(userData.Id, id, data.Name, data.Completed)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "UpdateTodo")
		return
	}

	log.Println(updatedTodo)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(updatedTodo)
}
