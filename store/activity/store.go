package activity

import (
	"context"
	"errors"
	"github.com/SemmiDev/todo-app/proto"
)

var (
	ErrActivityAlreadyExists = errors.New("activity already exists")
	ErrActivityNotFound      = errors.New("activity not found")
	ErrActivityIsEmpty       = errors.New("activity is empty")
)

type ActivityStore interface {
	Save(activity *proto.Activity) error
	Get(id string) (*proto.Activity, error)
	GetIdByDate(date string) ([]string, error)
	List() ([]*proto.Activity, error)
	Delete(id string) error
	Update(activity *proto.Activity) error
	Search(c context.Context, filter *proto.SearchActivityFilter, found func(activity *proto.Activity) error) error
}

func Copy(activity *proto.Activity) *proto.Activity {
	return &proto.Activity{
		Id:          activity.Id,
		Email:       activity.Email,
		Title:       activity.Title,
		Description: activity.Description,
		Day:         activity.Day,
		CreatedAt:   activity.CreatedAt,
		UpdatedAt:   activity.CreatedAt,
	}
}

func CopyAll(activities map[string]*proto.Activity) (data []*proto.Activity) {
	for _, v := range activities {
		data = append(data, Copy(v))
	}
	return
}
