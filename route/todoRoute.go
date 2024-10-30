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

	todoList, err := service.GetTodoList(userData.Id)

	if err != nil {
		errResponse := types.ErrorResponse{
			Message: err.Error(),
			Code:    500,
		}
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todoList)

}

func (h *TodoApi) AddTodo(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")
	var data types.TodoSchema
	bodyErr := json.NewDecoder(r.Body).Decode(&data)

	if r.Body == nil {
		common.ErrorResponse(w, "Empty request body", http.StatusBadRequest)
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	if data.Name == "" {
		common.ErrorResponse(w, "Name field can't be empty", http.StatusBadRequest)
		return
	}

	newTodo, err := service.AddTodo(data, userData.Id)

	if err != nil {
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func (h *TodoApi) DeleteTodo(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")

	id, parseErr := strconv.Atoi(r.PathValue("id"))
	if parseErr != nil {
		common.ErrorResponse(w, "Invalid id", http.StatusInternalServerError)
		return
	}

	err := service.DeleteTodo(id, userData.Id)

	if err != nil {
		if err.Error() == "Invalid id" {
			common.ErrorResponse(w, "Id Not Found", http.StatusBadRequest)
			return
		}
		common.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TodoApi) UpdateTodo(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")

	id, parseErr := strconv.Atoi(r.PathValue("id"))
	if parseErr != nil {
		common.ErrorResponse(w, "Invalid id", http.StatusInternalServerError)
		return
	}

	var data types.TodoSchema
	bodyErr := json.NewDecoder(r.Body).Decode(&data)

	if r.Body == nil {
		common.ErrorResponse(w, "Empty request body", http.StatusBadRequest)
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	if data.Name == "" {
		common.ErrorResponse(w, "Name field can't be empty", http.StatusBadRequest)
		return
	}

	updatedTodo, err := service.UpdateTodo(userData.Id, id, data.Name, data.Completed)

	if err != nil {
		if err.Error() == "invalid id" {
			common.ErrorResponse(w, "Invalid Id", http.StatusBadRequest)
			return
		}

		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(updatedTodo)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(updatedTodo)
}
