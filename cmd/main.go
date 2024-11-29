package main

import (
	grpc "github.com/VallabhSLEPAM/grpc-server/internal/adapters/grpc"
	app "github.com/VallabhSLEPAM/grpc-server/internal/application"
)

func main() {

	helloService := &app.HelloService{}

	grpcAdapter := grpc.NewGRPCAdapter(helloService, 9090)

	grpcAdapter.Run()

}
