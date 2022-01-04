gen:
	cd proto && \
	protoc --go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--grpc-gateway_out=. \
	--openapiv2_out=:swagger \
	--go-grpc_opt=paths=source_relative \
	model.proto && \
	cd proto && \
	mv model.pb.gw.go .. && \
	cd .. && \
	rmdir proto

run-GRPC-server:
	go run cmd/server/grpc/main.go
run-REST-server:
	go run cmd/server/rest/main.go