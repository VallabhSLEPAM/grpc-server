package main

import (
	grpc "github.com/VallabhSLEPAM/grpc-server/internal/adapters/grpc"
	app "github.com/VallabhSLEPAM/grpc-server/internal/application"
)

func main() {

	helloService := &app.HelloService{}
	bankService := &app.BankService{}

	grpcAdapter := grpc.NewGRPCAdapter(helloService, bankService, 9090)

	grpcAdapter.Run()

}
