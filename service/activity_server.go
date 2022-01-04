package service

import (
	"context"
	"github.com/SemmiDev/todo-app/common/seeder"
	"github.com/SemmiDev/todo-app/proto"
	"github.com/SemmiDev/todo-app/store/activity"
	"github.com/rs/zerolog"
	"log"

	ctx "github.com/SemmiDev/todo-app/common/context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ActivityServer struct {
	proto.UnimplementedActivityServiceServer
	activityStore activity.ActivityStore
	logger        *zerolog.Logger
}

func NewActivityServer(activityStore activity.ActivityStore, logger *zerolog.Logger) *ActivityServer {
	activityServer := &ActivityServer{
		activityStore: activityStore,
		logger:        logger,
	}

	// seed activities
	activities, _ := seeder.Seed(100)
	for i := 0; i < len(activities); i++ {
		activityServer.activityStore.Save(activities[i])
	}

	activityServer.logger.Info().Interface("SEED ACTIVITY", len(activities)).Msg("[COMPLETED]")
	return activityServer
}

func (h *ActivityServer) CreateActivity(c context.Context, req *proto.CreateActivityRequest) (*proto.CreateActivityResponse, error) {
	h.logger.Info().Interface("REQ", req).Msg("CREATE ACTIVITY")

	Activity := &proto.Activity{
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

	res := &proto.CreateActivityResponse{
		Activity: Activity,
	}
	return res, nil
}

func (h *ActivityServer) GetActivity(c context.Context, req *proto.GetActivityRequest) (*proto.GetActivityResponse, error) {
	h.logger.Info().Interface("REQ", req).Msg("GET ACTIVITY")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	Activity, err := h.activityStore.Get(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "activity with ID %s not found", req.GetId())
	}
	res := &proto.GetActivityResponse{
		Activity: Activity,
	}
	return res, nil
}

func (h *ActivityServer) ListActivity(c context.Context, req *proto.EmptyRequest) (*proto.ListActivityResponse, error) {
	h.logger.Info().Interface("REQ", req).Msg("LIST ACTIVITY")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	activities, err := h.activityStore.List()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "activities is empty")
	}

	res := &proto.ListActivityResponse{
		Activities: activities,
	}
	return res, nil
}

func (h *ActivityServer) SearchActivity(filter *proto.SearchActivityFilter, stream proto.ActivityService_SearchActivityServer) error {
	h.logger.Info().Interface("FILTER", filter).Msg("SEARCH ACTIVITY")
	err := h.activityStore.Search(
		stream.Context(),
		filter,
		func(activity *proto.Activity) error {
			res := &proto.SearchActivityResponse{
				Activity: activity,
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent activity with id: %s", activity.GetId())
			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}

func (h *ActivityServer) UpdateActivity(c context.Context, req *proto.UpdateActivityRequest) (*proto.UpdateActivityResponse, error) {
	h.logger.Info().Interface("REQ", req).Msg("UPDATE ACTIVITY")
	Activity := &proto.Activity{
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

	res := &proto.UpdateActivityResponse{
		Activity: Activity,
	}
	return res, nil
}

func (h *ActivityServer) DeleteActivity(c context.Context, req *proto.DeleteActivityRequest) (*proto.EmptyResponse, error) {
	h.logger.Info().Interface("REQ", req).Msg("DELETE ACTIVITY")
	if err := ctx.ContextError(c); err != nil {
		return nil, err
	}

	err := h.activityStore.Delete(req.Id)
	if err != nil {
		return &proto.EmptyResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "activity with ID %s not found", req.GetId())
	}
	return &proto.EmptyResponse{
		Success: true,
	}, nil
}
