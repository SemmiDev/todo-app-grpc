package service

import (
	"context"
	"fmt"
	"log"

	"github.com/SemmiDev/todo-app/model"
	"github.com/SemmiDev/todo-app/proto"

	"github.com/SemmiDev/todo-app/common/response"
	"github.com/SemmiDev/todo-app/common/token"
	"github.com/SemmiDev/todo-app/store/user"
	"google.golang.org/genproto/googleapis/api/httpbody"
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

	user1, _ := model.NewUser("sammi1", "sammi@gmail.com", "sammi1", "admin")
	user2, _ := model.NewUser("sammi2", "sammi@gmail.com", "sammi2", "user")

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

func (server *AuthServer) Register(c context.Context, req *proto.RegisterRequest) (*httpbody.HttpBody, error) {
	exists := server.userStore.ExistsByUsername(req.GetUsername()) || server.userStore.ExistsByEmail(req.GetEmail())
	if !exists {
		user, _ := model.NewUser(req.Username, req.Email, req.Password, "user")
		server.userStore.Save(user)
		resp := response.OKResponse(user)
		log.Println(resp)
		return resp, nil
	}
	err := fmt.Errorf("cannot find user: username/email already exists")
	return response.ErrAlreadyExists(err), nil
}
