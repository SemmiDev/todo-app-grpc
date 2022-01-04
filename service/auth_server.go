package service

import (
	"context"
	"github.com/SemmiDev/todo-app/model"
	"github.com/SemmiDev/todo-app/proto"

	"github.com/SemmiDev/todo-app/common/token"
	"github.com/SemmiDev/todo-app/store/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	userStore  user.UserStore
	jwtManager *token.JWTManager
}

func NewAuthServer(userStore user.UserStore, jwtManager *token.JWTManager) *AuthServer {
	authServer := &AuthServer{
		userStore:  userStore,
		jwtManager: jwtManager,
	}

	user1, _ := model.NewUser("sammi1", "sammi1", "admin")
	user2, _ := model.NewUser("sammi2", "sammi2", "user")

	authServer.userStore.Save(user1)
	authServer.userStore.Save(user2)

	return authServer
}

func (server *AuthServer) Login(c context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	u, err := server.userStore.Get(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if u == nil || !u.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	token, err := server.jwtManager.Generate(u)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &proto.LoginResponse{AccessToken: token}
	return res, nil
}
