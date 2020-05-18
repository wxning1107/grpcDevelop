package main

import (
	"context"
	"google.golang.org/grpc"
	"grpcClient/certificattedRpc"
	"io"
	"log"
	"net"
)

var (
	port = ":1107"
)

type HelloService struct {
}

func (p *HelloService) Hello(ctx context.Context, args *certificattedRpc.String) (*certificattedRpc.String, error) {
	reply := &certificattedRpc.String{Value: "pubsubGrpc" + args.GetValue()}

	return reply, nil
}

func (p *HelloService) Channel(stream certificattedRpc.HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &certificattedRpc.String{Value: "hello " + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

func main() {
	grpcServer := grpc.NewServer()
	// RegisterHelloServiceServer把HelloService的方法反射给grpcServer
	certificattedRpc.RegisterHelloServiceServer(grpcServer, new(HelloService))

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("could not listen on %s: %s", port, err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Panicf("grpc serve error: %s", err)
	}

}
