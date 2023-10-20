package router

import (
	"github.com/gorilla/mux"
	"github.com/jatin-02/todo/middleware"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/tasks", middleware.GetAllTasks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/task", middleware.InsertOneTask).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/task/{id}", middleware.TaskComplete).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/undoTask/{id}", middleware.UndoTask).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteTask/{id}", middleware.DeleteOneTask).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/deleteTasks", middleware.DeleteAllTasks).Methods("DELETE", "OPTIONS")

	return router
}