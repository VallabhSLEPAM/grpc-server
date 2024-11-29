package grpc

import (
	"context"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/hello"
)

func (adapter *GRPCAdapter) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {

	greet := adapter.helloService.GenerateHello(req.Name)

	return &hello.HelloResponse{
		Greet: greet,
	}, nil

}
