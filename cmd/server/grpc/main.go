package main

import (
	"fmt"
	"github.com/SemmiDev/todo-app/proto"
	"google.golang.org/grpc/reflection"
	"net"
	"os"

	"github.com/SemmiDev/todo-app/common/config"
	"github.com/SemmiDev/todo-app/common/token"

	"github.com/SemmiDev/todo-app/service"
	activityStore "github.com/SemmiDev/todo-app/store/activity"
	todoStore "github.com/SemmiDev/todo-app/store/todo"
	userStore "github.com/SemmiDev/todo-app/store/user"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func accessibleRoles() map[string][]string {
	const activityServicePath = "/proto.ActivityService/"
	const todoServicePath = "/proto.TodoService/"

	return map[string][]string{
		activityServicePath + "CreateActivity": {"admin", "user"},
		activityServicePath + "GetActivity":    {"admin", "user"},
		activityServicePath + "ListActivity":   {"admin", "user"},
		activityServicePath + "SearchActivity": {"admin", "user"},
		activityServicePath + "UpdateActivity": {"admin", "user"},
		activityServicePath + "DeleteActivity": {"admin", "user"},

		todoServicePath + "CreateTodo":             {"admin", "user"},
		todoServicePath + "GetTodo":                {"admin", "user"},
		todoServicePath + "ListTodo":               {"admin", "user"},
		todoServicePath + "ListTodoByActivityId":   {"admin", "user"},
		todoServicePath + "ListTodoByActivityDate": {"admin", "user"},
		todoServicePath + "SearchTodo":             {"admin", "user"},
		todoServicePath + "UpdateTodo":             {"admin", "user"},
		todoServicePath + "DeleteTodo":             {"admin", "user"},
	}
}

func main() {
	w := zerolog.ConsoleWriter{Out: os.Stderr}
	l := zerolog.New(w).With().Timestamp().Caller().Logger()

	jwtManager := token.NewJWTManager(config.SecretKey, config.TokenDuration)
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}

	userStore := userStore.NewMapStore()
	todoStore := todoStore.NewMapStore()
	activityStore := activityStore.NewMapStore()

	authServer := service.NewAuthServer(userStore, jwtManager)
	todoServer := service.NewTodoServer(todoStore, activityStore, &l)
	activityServer := service.NewActivityServer(activityStore, &l)

	grpcServer := grpc.NewServer(serverOpts...)

	listener, err := net.Listen("tcp", config.GRPCServerPort)
	if err != nil {
		l.Panic().Err(fmt.Errorf("failed to listen: %w", err))
	}

	proto.RegisterAuthServiceServer(grpcServer, authServer)
	proto.RegisterTodoServiceServer(grpcServer, todoServer)
	proto.RegisterActivityServiceServer(grpcServer, activityServer)
	reflection.Register(grpcServer)

	l.Info().Str("ADDR", config.GRPCServerPort).Msg("Start GRPC Server")
	if err := grpcServer.Serve(listener); err != nil {
		l.Panic().Err(fmt.Errorf("failed to serve: %w", err))
	}
}
