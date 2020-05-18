package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpcClient/grpcDeadline"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:1107", grpc.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	client := grpcDeadline.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &grpcDeadline.String{Value: " grpc"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Println(reply)

	stream, err := client.Channel(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			num := rand.Intn(100)
			if err := stream.Send(&grpcDeadline.String{Value: "grpc" + strconv.Itoa(num)}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}()

	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		fmt.Println(reply.GetValue())

	}
}
