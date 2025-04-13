package handlers

import (
	"gorilla-mux-router/models"
	"gorilla-mux-router/utils"
	"net/http"
	"strconv"
)

var (
	todos []models.Todo = []models.Todo{
		{Id: 1, Title: "Pray", Body: "Pray 5 times everyday", Completed: true},
	}
)

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	var content []byte // the content to send as response

	// get the value of query completed eg: completed?=true
	query := r.URL.Query()["completed"]

	// if the num of queries is not 0, then send the filtered response
	if len(query) != 0 {
		isCompleted, _ := strconv.ParseBool(query[0])
		completedTodos := []models.Todo{}
		for _, todo := range todos {
			if todo.Completed == isCompleted {
				completedTodos = append(completedTodos, todo)
			}
			content, _ = utils.Json(completedTodos)
		}
	} else {
		content, _ = utils.Json(todos)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(content)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {

}
