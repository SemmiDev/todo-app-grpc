package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/SemmiDev/todo-app/common/config"
	"log"
	"net/http"
	"strconv"

	"github.com/SemmiDev/todo-app/model"
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
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := model.RegisterTodoServiceHandlerFromEndpoint(ctx, mux, config.SERVER_PORT, opts)
	if err != nil {
		return err
	}

	err = model.RegisterActivityServiceHandlerFromEndpoint(ctx, mux, config.SERVER_PORT, opts)
	if err != nil {
		return err
	}

	log.Println("starting gateway server on port 8081")
	return http.ListenAndServe(config.GATEWAY_PORT, mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
