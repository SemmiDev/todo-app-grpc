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

func (h *Handler) CreateActivity(
	c context.Context,
	req *model.CreateActivityRequest,
) (*model.CreateActivityResponse, error) {
	h.logger.Info().Interface("req", req).Msg("create activity")

	Activity := &model.Activity{
		Id:          uuid.NewString(),
		Email:       req.Email,
		Title:       req.Title,
		Description: req.Description,
		Day:         req.Day,
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
	}

	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	err := h.activityStore.Save(Activity)
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "activity with ID %s already exists", Activity.GetId())
	}

	res := &model.CreateActivityResponse{
		Activity: Activity,
	}
	return res, nil
}

func (h *Handler) GetActivity(c context.Context, req *model.GetActivityRequest) (*model.GetActivityResponse, error) {
	h.logger.Info().Interface("req", req).Msg("get activity")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	Activity, err := h.activityStore.Get(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "activity with ID %s not found", req.GetId())
	}
	res := &model.GetActivityResponse{
		Activity: Activity,
	}
	return res, nil
}

func (h *Handler) ListActivity(c context.Context, req *model.EmptyRequest) (*model.ListActivityResponse, error) {
	h.logger.Info().Interface("req", req).Msg("list activity")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	Activitys, err := h.activityStore.List()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "activities is empty")
	}

	res := &model.ListActivityResponse{
		Activities: Activitys,
	}
	return res, nil
}

func (h *Handler) SearchActivity(filter *model.SearchActivityFilter, stream model.ActivityService_SearchActivityServer) error {
	h.logger.Info().Interface("filter", filter).Msg("search activity")
	err := h.activityStore.Search(
		stream.Context(),
		filter,
		func(Activity *model.Activity) error {
			res := &model.SearchActivityResponse{
				Activity: Activity,
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent activity with id: %s", Activity.GetId())
			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}

func (h *Handler) DeleteActivity(c context.Context, req *model.DeleteActivityRequest) (*model.EmptyResponse, error) {
	h.logger.Info().Interface("req", req).Msg("delete activity")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	err := h.activityStore.Delete(req.Id)
	if err != nil {
		return &model.EmptyResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "activity with ID %s not found", req.GetId())
	}
	return &model.EmptyResponse{
		Success: true,
	}, nil
}

func (h *Handler) UpdateActivity(c context.Context, req *model.UpdateActivityRequest) (*model.UpdateActivityResponse, error) {
	h.logger.Info().Interface("req", req).Msg("update activity")
	Activity := &model.Activity{
		Id:          req.Id,
		Email:       req.Email,
		Title:       req.Title,
		Description: req.Description,
		Day:         req.Day,
		UpdatedAt:   timestamppb.Now(),
	}
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	err := h.activityStore.Update(Activity)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "activity with ID %s not found", req.GetId())
	}

	res := &model.UpdateActivityResponse{
		Activity: Activity,
	}
	return res, nil
}
