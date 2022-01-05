package service_test

// import (
// 	"context"
// 	"net"
// 	"testing"

// 	"github.com/SemmiDev/todo-app/common/serializer"
// 	"github.com/SemmiDev/todo-app/common/token"
// 	"github.com/SemmiDev/todo-app/proto"
// 	"github.com/SemmiDev/todo-app/service"
// 	"github.com/SemmiDev/todo-app/store/user"
// 	"github.com/stretchr/testify/require"
// 	"google.golang.org/genproto/googleapis/api/httpbody"
// 	"google.golang.org/grpc"
// )

// func TestClientRegisterUser(t *testing.T) {
// 	t.Parallel()

// 	userStore := user.NewMapStore()

// 	laptopStore := service.NewInMemoryLaptopStore()
// 	serverAddress := startTestLaptopServer(t, laptopStore, nil, nil)
// 	laptopClient := newTestLaptopClient(t, serverAddress)

// 	laptop := sample.NewLaptop()
// 	expectedID := laptop.Id
// 	req := &pb.CreateLaptopRequest{
// 		Laptop: laptop,
// 	}

// 	res, err := laptopClient.CreateLaptop(context.Background(), req)
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// 	require.Equal(t, expectedID, res.Id)

// 	// check that the laptop is saved to the store
// 	other, err := laptopStore.Find(res.Id)
// 	require.NoError(t, err)
// 	require.NotNil(t, other)

// 	// check that the saved laptop is the same as the one we send
// 	requireSameLaptop(t, laptop, other)
// }

// func startTestServer(t *testing.T, userStore user.UserStore, token *token.JWTManager) string {
// 	authServer := service.NewAuthServer(userStore, token)
// 	grpcServer := grpc.NewServer()
// 	proto.RegisterAuthServiceServer(grpcServer, authServer)

// 	listener, err := net.Listen("tcp", ":0") // random available port
// 	require.NoError(t, err)
// 	go grpcServer.Serve(listener)
// 	return listener.Addr().String()
// }

// func requireSameLaptop(t *testing.T, user1 *httpbody.HttpBody, user2 *httpbody.HttpBody) {
// 	json1, err := serializer.ProtobufToJSON()
// 	require.NoError(t, err)

// 	json2, err := serializer.ProtobufToJSON(laptop2)
// 	require.NoError(t, err)

// 	require.Equal(t, json1, json2)
// }
