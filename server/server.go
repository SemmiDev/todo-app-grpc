package main

import (
	"fmt"
	"net"
	"os"

	"github.com/SemmiDev/todo-app/common/config"
	"github.com/SemmiDev/todo-app/common/token"

	"github.com/SemmiDev/todo-app/handler"
	"github.com/SemmiDev/todo-app/model"
	activityStore "github.com/SemmiDev/todo-app/store/activity"
	todoStore "github.com/SemmiDev/todo-app/store/todo"
	userStore "github.com/SemmiDev/todo-app/store/user"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func accessibleRoles() map[string][]string {
	const activityServicePath = "/model.ActivityService/"
	const todoServicePath = "/model.TodoService/"

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

	userStore := userStore.NewMapStore()
	todoStore := todoStore.NewMapStore()
	activityStore := activityStore.NewMapStore()

	handlers := handler.New(&l, todoStore, activityStore)

	listener, err := net.Listen("tcp", config.ServerPort)
	if err != nil {
		l.Panic().Err(fmt.Errorf("failed to listen: %w", err))
	}

	jwtManager := token.NewJWTManager(config.SecretKey, config.TokenDuration)
	authServer := handler.NewAuthServer(userStore, jwtManager)
	interceptor := handler.NewAuthInterceptor(jwtManager, accessibleRoles())

	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
		),
	}

	grpcServer := grpc.NewServer(serverOpts...)

	model.RegisterAuthServiceServer(grpcServer, authServer)
	model.RegisterTodoServiceServer(grpcServer, handlers)
	model.RegisterActivityServiceServer(grpcServer, handlers)

	l.Info().Str("port", config.ServerPort).Msg("Start GRPC Server")
	if err := grpcServer.Serve(listener); err != nil {
		l.Panic().Err(fmt.Errorf("failed to serve: %w", err))
	}
}
