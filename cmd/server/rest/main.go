package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/SemmiDev/todo-app/proto"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"strconv"

	"github.com/SemmiDev/todo-app/common/config"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type ErrBody struct {
	Message string `json:"message,omitempty"`
	Details string `json:"detail"`
}

func CustomHTTPError(
	_ context.Context,
	_ *runtime.ServeMux,
	_ runtime.Marshaler,
	w http.ResponseWriter,
	_ *http.Request,
	err error,
) {
	const fallback = `{"error": "failed to marshal error message"}`
	w.WriteHeader(runtime.HTTPStatusFromCode(status.Code(err)))
	jErr := json.NewEncoder(w).Encode(ErrBody{
		Message: status.Convert(err).Message(),
		Details: strconv.Itoa(int(status.Code(err))),
	})
	if jErr != nil {
		w.Write([]byte(fallback))
	}
}

func run() error {
	w := zerolog.ConsoleWriter{Out: os.Stderr}
	l := zerolog.New(w).With().Timestamp().Caller().Logger()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := proto.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, config.GRPCServerPort, opts)
	if err != nil {
		return err
	}

	err = proto.RegisterTodoServiceHandlerFromEndpoint(ctx, mux, config.GRPCServerPort, opts)
	if err != nil {
		return err
	}

	err = proto.RegisterActivityServiceHandlerFromEndpoint(ctx, mux, config.GRPCServerPort, opts)
	if err != nil {
		return err
	}

	l.Info().Str("ADDR", config.RestServerPort).Msg("Start RESET Server")
	return http.ListenAndServe(config.RestServerPort, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
