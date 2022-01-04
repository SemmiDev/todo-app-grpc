package seeder

import (
	"github.com/SemmiDev/todo-app/common/random"
	"github.com/SemmiDev/todo-app/proto"
)

func Seed(n int) (activities []*proto.Activity, todos []*proto.Todo) {
	for i := 0; i < n; i++ {
		activity := &proto.Activity{
			Id:          random.RandomID(),
			Email:       random.RandomEmail(),
			Title:       random.RandomTitle(),
			Description: random.RandomDescription(),
			Day:         random.RandomDay(),
			CreatedAt:   random.RandomCreatedAt(),
			UpdatedAt:   random.RandomUpdatedAt(),
		}

		todo := &proto.Todo{
			Id:          random.RandomID(),
			ActivityId:  activity.GetId(),
			Title:       random.RandomTitle(),
			Description: random.RandomDescription(),
			Priority:    random.RandomPriority(),
			Status:      random.RandomStatus(),
			CreatedAt:   random.RandomCreatedAt(),
			UpdatedAt:   random.RandomUpdatedAt(),
		}

		activities = append(activities, activity)
		todos = append(todos, todo)
	}

	return
}
