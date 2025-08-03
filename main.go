package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// todo represents a single todo item.
// ID is a unique identifier for the todo.
// Item is the description of the task.
// Completed indicates whether the task is done.
type todo struct {
	ID        string `json:"id"`
	Item      string `json:"title"`
	Completed bool   `json:"completed"`
}

// todos is an in-memory slice that stores all todo items.
var todos = []todo{
	{ID: "1", Item: "Learn Go", Completed: false},
	{ID: "2", Item: "Build a web app", Completed: false},
	{ID: "3", Item: "Deploy to production", Completed: false},
}

// getTodos handles GET /todos
// Responds with the list of all todos in JSON format.
func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

// addTodo handles POST /todos
// Adds a new todo item with an auto-incremented ID.
// Expects a JSON body with "title" and "completed" fields.
func addTodo(context *gin.Context) {
	var newTodo todo
	if err := context.BindJSON(&newTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	// Generate new ID by incrementing the last ID
	if len(todos) == 0 {
		newTodo.ID = "1"
	} else {
		lastID := todos[len(todos)-1].ID
		// Convert lastID to int, increment, then back to string
		var nextID int
		fmt.Sscanf(lastID, "%d", &nextID)
		nextID++
		newTodo.ID = fmt.Sprintf("%d", nextID)
	}
	todos = append(todos, newTodo)
	context.IndentedJSON(http.StatusCreated, newTodo)
}

// deleteTodo handles DELETE /todos/:id
// Deletes a todo item by its ID.
func deleteTodo(context *gin.Context) {
	id := context.Param("id")
	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			context.Status(http.StatusNoContent)
			return
		}
	}
	context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

// editTodo handles PUT /todos/:id
// Updates the title and completed status of a todo item by its ID.
// Expects a JSON body with "title" and "completed" fields.
func editTodo(context *gin.Context) {
	id := context.Param("id")
	var updatedTodo todo
	if err := context.BindJSON(&updatedTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	for i, t := range todos {
		if t.ID == id {
			todos[i].Item = updatedTodo.Item
			todos[i].Completed = updatedTodo.Completed
			context.IndentedJSON(http.StatusOK, todos[i])
			return
		}
	}
	context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

// getTodoByID handles GET /todos/:id
// Responds with a single todo item by its ID.
func getTodoByID(context *gin.Context) {
	id := context.Param("id")
	for _, t := range todos {
		if t.ID == id {
			context.IndentedJSON(http.StatusOK, t)
			return
		}
	}
	context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

// main sets up the Gin router and routes, then starts the server.
func main() {
	router := gin.Default()
	router.GET("/todos", func(context *gin.Context) {
		getTodos(context)
	})
	router.POST("/todos", func(context *gin.Context) {
		addTodo(context)
	})
	router.DELETE("/todos/:id", func(context *gin.Context) {
		deleteTodo(context)
	})
	router.PUT("/todos/:id", func(context *gin.Context) {
		editTodo(context)
	})
	router.GET("/todos/:id", func(context *gin.Context) {
		getTodoByID(context)
	})

	router.Run(":8080")

}
