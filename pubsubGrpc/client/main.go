package main

import (
	"context"
	"google.golang.org/grpc"
	pubsub "grpcClient/pubsubGrpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:0215", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pubsub.NewPubSubServiceClient(conn)
	_, err = client.Publish(context.Background(), &pubsub.String{Value: "Golang: hello Go"})
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Publish(context.Background(), &pubsub.String{Value: "docekr: hello docker"})
	if err != nil {
		log.Fatal(err)
	}
}
