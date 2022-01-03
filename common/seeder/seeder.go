package seeder

import (
	"github.com/SemmiDev/todo-app/common/random"
	"github.com/SemmiDev/todo-app/model"
)

func SeedTodos() (todosReq []*model.CreateTodoRequest) {
	for i := 0; i < 100; i++ {
		todosReq = append(todosReq, random.RandomTodo())
	}
	return
}

func SeedActivities() (todosReq []*model.CreateActivityRequest) {
	for i := 0; i < 100; i++ {
		todosReq = append(todosReq, random.RandomActivity())
	}
	return
}
