package activity

import (
	"context"
	"errors"

	"github.com/SemmiDev/todo-app/model"
)

var (
	ErrActivityAlreadyExists = errors.New("activity already exists")
	ErrActivityNotFound      = errors.New("activity not found")
	ErrActivityIsEmpty       = errors.New("activity is empty")
)

type ActivityStore interface {
	Save(activity *model.Activity) error
	Get(id string) (*model.Activity, error)
	GetIdByDate(date string) ([]string, error)
	List() ([]*model.Activity, error)
	Delete(id string) error
	Update(activity *model.Activity) error
	Search(c context.Context, filter *model.SearchActivityFilter, found func(activity *model.Activity) error) error
}

func Copy(activity *model.Activity) *model.Activity {
	return &model.Activity{
		Id:          activity.Id,
		Email:       activity.Email,
		Title:       activity.Title,
		Description: activity.Description,
		Day:         activity.Day,
		CreatedAt:   activity.CreatedAt,
		UpdatedAt:   activity.CreatedAt,
	}
}

func CopyAll(activities map[string]*model.Activity) (data []*model.Activity) {
	for _, v := range activities {
		data = append(data, Copy(v))
	}
	return
}
