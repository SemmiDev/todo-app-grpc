package random

import (
	"github.com/SemmiDev/todo-app/model"
)

func RandomTodo() *model.CreateTodoRequest {
	req := model.CreateTodoRequest{
		Title:       RandomTitle(),
		Description: RandomDescription(),
		Priority:    RandomPriority(),
		Status:      RandomStatus(),
	}
	return &req
}

func RandomActivity() *model.CreateActivityRequest {
	req := model.CreateActivityRequest{
		Email:       RandomEmail(),
		Title:       RandomTitle(),
		Description: RandomDescription(),
		Day:         RandomDay(),
	}
	return &req
}
