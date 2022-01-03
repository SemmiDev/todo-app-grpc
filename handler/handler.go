package handler

import (
	"context"

	"github.com/SemmiDev/todo-app/common/seeder"
	"github.com/SemmiDev/todo-app/model"
	"github.com/SemmiDev/todo-app/store/activity"
	"github.com/SemmiDev/todo-app/store/todo"
	"github.com/rs/zerolog"
)

type Handler struct {
	logger        *zerolog.Logger
	todoStore     todo.TodoStore
	activityStore activity.ActivityStore
	model.UnimplementedTodoServiceServer
	model.UnimplementedActivityServiceServer
}

func New(l *zerolog.Logger, ts todo.TodoStore, as activity.ActivityStore) *Handler {
	h := &Handler{
		logger:        l,
		todoStore:     ts,
		activityStore: as,
	}

	// seed todos
	todosReq := seeder.SeedTodos()
	for i := 0; i < len(todosReq); i++ {
		_, _ = h.CreateTodo(context.Background(), todosReq[i])
	}
	// seed activities
	activitiesReq := seeder.SeedActivities()
	for i := 0; i < len(activitiesReq); i++ {
		_, _ = h.CreateActivity(context.Background(), activitiesReq[i])
	}

	h.logger.Info().Interface("seed todos", len(todosReq)).Msg("[COMPLETED]")
	h.logger.Info().Interface("seed todos", len(activitiesReq)).Msg("[COMPLETED]")

	return h
}
