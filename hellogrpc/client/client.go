package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpcClient/hellogrpc"
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

	client := hellogrpc.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &hellogrpc.String{Value: " grpc"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)

	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			num := rand.Intn(100)
			if err := stream.Send(&hellogrpc.String{Value: "grpc" + strconv.Itoa(num)}); err != nil {
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
