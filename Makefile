gen:
	cd proto && \
	protoc --go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--grpc-gateway_out=. \
	--openapiv2_out=:swagger \
	--go-grpc_opt=paths=source_relative \
	model.proto && \
	mv model.pb.go model_grpc.pb.go model/model.pb.gw.go ../model && \
	rm -r model && cd ..

run-server:
	go run server/server.go
run-gateway:
	go run gateway/gateway.go