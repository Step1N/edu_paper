gen:
	protoc --proto_path=eduproto eduproto/*.proto --go_out=plugins=grpc:edupb --grpc-gateway_out=:edupb --swagger_out=:swagger

clean:
	rm edupb/*.go 

server:
	go run cmd/eduserv/main.go -port 8080

client:
	go run cmd/educlient/main.go -address 0.0.0.0:8080

test:
	go test -cover -race ./...

cert:
	cd cert; ./gen.sh; cd ..

.PHONY: gen clean server client test cert 