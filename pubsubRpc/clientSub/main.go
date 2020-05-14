package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pubsub "grpcClient/pubsubRpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:0215", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pubsub.NewPubSubServiceClient(conn)

	stream, err := client.Subscribe(context.Background(), &pubsub.String{Value: "Golang:"})
	if err != nil {
		log.Fatal(err)
	}

	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		fmt.Println(reply)
	}

}
