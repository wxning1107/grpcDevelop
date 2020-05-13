package main

import (
	"context"
	"google.golang.org/grpc"
	"grpcClient/hellogrpc"
	"io"

	"log"
	"net"
)

type HelloService struct {
}

func (p *HelloService) Hello(ctx context.Context, args *hellogrpc.String) (*hellogrpc.String, error) {
	reply := &hellogrpc.String{Value: "hello" + args.GetValue()}

	return reply, nil
}

func (p *HelloService) Channel(stream hellogrpc.HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &hellogrpc.String{Value: "hello " + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

func main() {
	grpcServer := grpc.NewServer()
	hellogrpc.RegisterHelloServiceServer(grpcServer, new(HelloService))

	listener, err := net.Listen("tcp", ":1107")
	if err != nil {
		log.Fatal(err)
	}

	_ = grpcServer.Serve(listener)

}

