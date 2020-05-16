package main

import (
	"context"
	"github.com/docker/docker/pkg/pubsub"
	"google.golang.org/grpc"
	pb "grpcClient/pubsubGrpc"
	"log"
	"net"
	"strings"
	"time"
)

type PubSubService struct {
	pub *pubsub.Publisher
}

func (p *PubSubService) Publish(ctx context.Context, arg *pb.String) (*pb.String, error) {
	p.pub.Publish(arg.GetValue())
	return &pb.String{}, nil
}

func (p *PubSubService) Subscribe(arg *pb.String, stream pb.PubSubService_SubscribeServer) error {
	ch := p.pub.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, arg.GetValue()) {
				return true
			}
		}
		return false
	})

	for v := range ch {
		if err := stream.Send(&pb.String{Value: v.(string)}); err != nil {
			return err
		}
	}

	return nil
}

func NewPubSubService() *PubSubService {
	return &PubSubService{pub: pubsub.NewPublisher(100*time.Millisecond, 10)}
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterPubSubServiceServer(grpcServer, NewPubSubService())

	listener, err := net.Listen("tcp", ":0215")
	if err != nil {
		log.Fatal(err)
	}

	_ = grpcServer.Serve(listener)
}
