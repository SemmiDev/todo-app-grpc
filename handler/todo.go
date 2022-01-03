package handler

import (
	"context"
	"github.com/SemmiDev/todo-app/common/context"
	"log"

	"github.com/SemmiDev/todo-app/model"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateTodo(
	c context.Context,
	req *model.CreateTodoRequest,
) (*model.CreateTodoResponse, error) {
	h.logger.Info().Interface("req", req).Msg("create todo")

	reqActivity := &model.GetActivityRequest{
		Id: req.GetActivityId(),
	}
	if _, err := h.GetActivity(c, reqActivity); err != nil {
		return nil, status.Errorf(codes.NotFound, "activity with ID %s not found", req.GetActivityId())
	}

	todo := &model.Todo{
		Id:          uuid.NewString(),
		Title:       req.Title,
		ActivityId:  req.ActivityId,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	}

	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	err := h.todoStore.Save(todo)
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "todo with ID %s already exists", todo.GetId())
	}

	res := &model.CreateTodoResponse{
		Todo: todo,
	}
	return res, nil
}

func (h *Handler) GetTodo(c context.Context, req *model.GetTodoRequest) (*model.GetTodoResponse, error) {
	h.logger.Info().Interface("req", req).Msg("get todo")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	todo, err := h.todoStore.Get(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "todo with ID %s not found", req.GetId())
	}
	res := &model.GetTodoResponse{
		Todo: todo,
	}
	return res, nil
}

func (h *Handler) ListTodo(c context.Context, req *model.EmptyRequest) (*model.ListTodoResponse, error) {
	h.logger.Info().Interface("req", req).Msg("list todo")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	todos, err := h.todoStore.List()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "todos is empty")
	}

	res := &model.ListTodoResponse{
		List: todos,
	}
	return res, nil
}

func (h *Handler) SearchTodo(filter *model.SearchTodoFilter, stream model.TodoService_SearchTodoServer) error {
	h.logger.Info().Interface("filter", filter).Msg("search todo")
	err := h.todoStore.Search(
		stream.Context(),
		filter,
		func(todo *model.Todo) error {
			res := &model.SearchTodoResponse{
				Todo: todo,
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent todo with id: %s", todo.GetId())
			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}

func (h *Handler) DeleteTodo(c context.Context, req *model.DeleteTodoRequest) (*model.EmptyResponse, error) {
	h.logger.Info().Interface("req", req).Msg("delete todo")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	err := h.todoStore.Delete(req.Id)
	if err != nil {
		return &model.EmptyResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "todo with ID %s not found", req.GetId())
	}
	return &model.EmptyResponse{
		Success: true,
	}, nil
}

func (h *Handler) UpdateTodo(c context.Context, req *model.UpdateTodoRequest) (*model.UpdateTodoResponse, error) {
	h.logger.Info().Interface("req", req).Msg("update todo")
	todo := &model.Todo{
		Id:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		UpdatedAt:   timestamppb.Now(),
	}
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	err := h.todoStore.Update(todo)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "todo with ID %s not found", req.GetId())
	}

	res := &model.UpdateTodoResponse{
		Todo: todo,
	}
	return res, nil
}

func (h *Handler) ListTodoByActivityId(
	c context.Context,
	req *model.ListTodoByActivityIdRequest,
) (*model.ListTodoByActivityIdResponse, error) {
	h.logger.Info().Interface("req", req).Msg("list todo by activity id")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	todos, err := h.todoStore.ListByActivityId(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "todos is empty")
	}

	res := &model.ListTodoByActivityIdResponse{
		Todos: todos,
	}
	return res, nil
}

func (h *Handler) ListTodoByActivityDate(
	c context.Context,
	req *model.ListTodoByActivityDateRequest,
) (*model.ListTodoByActivityDateResponse, error) {
	h.logger.Info().Interface("req", req).Msg("list todo by activity date")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	ids, err := h.activityStore.GetIdByDate(req.GetDate())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "activity on %s not found", req.GetDate())
	}

	todos, err := h.todoStore.ListByActivityIds(ids)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "todo with activities id %s is empty", ids)
	}

	res := &model.ListTodoByActivityDateResponse{
		Todos: todos,
	}
	return res, nil
}
