package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpcClient/grpcDeadline"
	"io"
	"log"
	"net"
)

var (
	port = ":1107"
)

type HelloService struct {
}

func (p *HelloService) Hello(ctx context.Context, args *grpcDeadline.String) (*grpcDeadline.String, error) {
	reply := &grpcDeadline.String{Value: "pubsubGrpc" + args.GetValue()}

	return reply, nil
}

func (p *HelloService) Channel(stream grpcDeadline.HelloService_ChannelServer) error {
	for {
		if stream.Context().Err() == context.Canceled {
			return status.Errorf(codes.Canceled, "client context is canceled")
		}

		//select {
		//case <-time.After(1 * time.Second):
		//	fmt.Println("overslept")
		//case <-stream.Context().Done():
		//	fmt.Println(stream.Context().Err())
		//	return status.Errorf(codes.Canceled, "client context is canceled")
		//}

		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &grpcDeadline.String{Value: "hello " + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

func main() {
	grpcServer := grpc.NewServer()
	// RegisterHelloServiceServer把HelloService的方法反射给grpcServer
	grpcDeadline.RegisterHelloServiceServer(grpcServer, new(HelloService))

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Panicf("could not listen on %s: %s", port, err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Panicf("grpc serve error: %s", err)
	}

}
