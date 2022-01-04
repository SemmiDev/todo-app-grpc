package handler

import (
	"context"

	"github.com/SemmiDev/todo-app/common/token"
	"github.com/SemmiDev/todo-app/model"
	"github.com/SemmiDev/todo-app/store/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer is the server for authentication
type AuthServer struct {
	model.UnimplementedAuthServiceServer
	userStore  user.UserStore
	jwtManager *token.JWTManager
}

func NewAuthServer(userStore user.UserStore, jwtManager *token.JWTManager) *AuthServer {
	return &AuthServer{userStore: userStore, jwtManager: jwtManager}
}

// Login is a unary RPC to login user
func (server *AuthServer) Login(c context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
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

	res := &model.LoginResponse{AccessToken: token}
	return res, nil
}
