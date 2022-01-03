package main

import (
	"fmt"
	"github.com/SemmiDev/todo-app/common/config"
	"net"
	"os"

	"github.com/SemmiDev/todo-app/handler"
	"github.com/SemmiDev/todo-app/model"
	activityStore "github.com/SemmiDev/todo-app/store/activity"
	todoStore "github.com/SemmiDev/todo-app/store/todo"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func main() {
	w := zerolog.ConsoleWriter{Out: os.Stderr}
	l := zerolog.New(w).With().Timestamp().Caller().Logger()

	ts := todoStore.NewMapStore()
	as := activityStore.NewMapStore()
	h := handler.New(&l, ts, as)

	lis, err := net.Listen("tcp", config.SERVER_PORT)
	if err != nil {
		l.Panic().Err(fmt.Errorf("failed to listen: %w", err))
	}

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
		),
	)

	model.RegisterTodoServiceServer(s, h)
	model.RegisterActivityServiceServer(s, h)

	l.Info().Str("port", config.SERVER_PORT).Msg("starting server")
	if err := s.Serve(lis); err != nil {
		l.Panic().Err(fmt.Errorf("failed to serve: %w", err))
	}
}
