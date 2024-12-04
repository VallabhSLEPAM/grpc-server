package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/hello"
	"google.golang.org/grpc"
)

func (adapter *GRPCAdapter) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {

	greet := adapter.helloService.GenerateHello(req.Name)

	return &hello.HelloResponse{
		Greet: greet,
	}, nil

}

func (adapter *GRPCAdapter) HelloServerStream(req *hello.HelloRequest, ss grpc.ServerStreamingServer[hello.HelloResponse]) error {

	for i := 0; i < 10; i++ {
		greet := adapter.helloService.GenerateHello(req.Name)
		res := fmt.Sprintf("[%d] %s", i, greet)
		ss.Send(
			&hello.HelloResponse{
				Greet: res,
			},
		)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (adapter *GRPCAdapter) HelloClientStream(clientStreamServer grpc.ClientStreamingServer[hello.HelloRequest, hello.HelloResponse]) error {
	res := ""

	for {
		req, err := clientStreamServer.Recv()
		if err == io.EOF {
			return clientStreamServer.SendAndClose(
				&hello.HelloResponse{
					Greet: res,
				},
			)
		}
		if err != nil {
			log.Fatalln("Error receiving data: ", err)
		}
		greet := adapter.helloService.GenerateHello(req.Name)
		res += greet + " "
	}
}

func (adapter *GRPCAdapter) HelloContinuous(stream grpc.BidiStreamingServer[hello.HelloRequest, hello.HelloResponse]) error {

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalln("Error receiiving data: ", err)
		}

		greet := adapter.helloService.GenerateHello(req.Name)

		err = stream.Send(&hello.HelloResponse{
			Greet: greet,
		})
		if err != nil {
			log.Fatalln("Error sending data to client: ", err)
		}
	}

}
