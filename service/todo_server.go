package service

import (
	"context"
	"github.com/SemmiDev/todo-app/common/context"
	"github.com/SemmiDev/todo-app/common/seeder"
	"github.com/SemmiDev/todo-app/proto"
	"github.com/SemmiDev/todo-app/store/activity"
	"github.com/SemmiDev/todo-app/store/todo"
	"github.com/rs/zerolog"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TodoServer struct {
	proto.UnimplementedTodoServiceServer
	todoStore     todo.TodoStore
	activityStore activity.ActivityStore
	logger        *zerolog.Logger
}

func NewTodoServer(todoStore todo.TodoStore, activityStore activity.ActivityStore, logger *zerolog.Logger) *TodoServer {
	todoServer := &TodoServer{
		todoStore:     todoStore,
		activityStore: activityStore,
		logger:        logger,
	}

	// seed todos
	_, todos := seeder.Seed(100)
	for i := 0; i < len(todos); i++ {
		todoServer.todoStore.Save(todos[i])
	}

	todoServer.logger.Info().Interface("SEED TODO", len(todos)).Msg("[COMPLETED]")
	return todoServer
}

func (h *TodoServer) CreateTodo(
	c context.Context,
	req *proto.CreateTodoRequest,
) (*proto.CreateTodoResponse, error) {
	h.logger.Info().Interface("req", req).Msg("create todo")

	if _, err := h.activityStore.Get(req.ActivityId); err != nil {
		return nil, status.Errorf(codes.NotFound, "activity with ID %s not found", req.GetActivityId())
	}

	todo := &proto.Todo{
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

	res := &proto.CreateTodoResponse{
		Todo: todo,
	}
	return res, nil
}

func (h *TodoServer) GetTodo(c context.Context, req *proto.GetTodoRequest) (*proto.GetTodoResponse, error) {
	h.logger.Info().Interface("req", req).Msg("get todo")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	todo, err := h.todoStore.Get(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "todo with ID %s not found", req.GetId())
	}
	res := &proto.GetTodoResponse{
		Todo: todo,
	}
	return res, nil
}

func (h *TodoServer) ListTodo(c context.Context, req *proto.EmptyRequest) (*proto.ListTodoResponse, error) {
	h.logger.Info().Interface("req", req).Msg("list todo")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	todos, err := h.todoStore.List()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "todos is empty")
	}

	res := &proto.ListTodoResponse{
		List: todos,
	}
	return res, nil
}

func (h *TodoServer) SearchTodo(filter *proto.SearchTodoFilter, stream proto.TodoService_SearchTodoServer) error {
	h.logger.Info().Interface("filter", filter).Msg("search todo")
	err := h.todoStore.Search(
		stream.Context(),
		filter,
		func(todo *proto.Todo) error {
			res := &proto.SearchTodoResponse{
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

func (h *TodoServer) DeleteTodo(c context.Context, req *proto.DeleteTodoRequest) (*proto.EmptyResponse, error) {
	h.logger.Info().Interface("req", req).Msg("delete todo")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	err := h.todoStore.Delete(req.Id)
	if err != nil {
		return &proto.EmptyResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "todo with ID %s not found", req.GetId())
	}
	return &proto.EmptyResponse{
		Success: true,
	}, nil
}

func (h *TodoServer) UpdateTodo(c context.Context, req *proto.UpdateTodoRequest) (*proto.UpdateTodoResponse, error) {
	h.logger.Info().Interface("req", req).Msg("update todo")
	todo := &proto.Todo{
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

	res := &proto.UpdateTodoResponse{
		Todo: todo,
	}
	return res, nil
}

func (h *TodoServer) ListTodoByActivityId(
	c context.Context,
	req *proto.ListTodoByActivityIdRequest,
) (*proto.ListTodoByActivityIdResponse, error) {
	h.logger.Info().Interface("req", req).Msg("list todo by activity id")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	todos, err := h.todoStore.ListByActivityId(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "todos is empty")
	}

	res := &proto.ListTodoByActivityIdResponse{
		Todos: todos,
	}
	return res, nil
}

func (h *TodoServer) ListTodoByActivityDate(
	c context.Context,
	req *proto.ListTodoByActivityDateRequest,
) (*proto.ListTodoByActivityDateResponse, error) {
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

	res := &proto.ListTodoByActivityDateResponse{
		Todos: todos,
	}
	return res, nil
}
