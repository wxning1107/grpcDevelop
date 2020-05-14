package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpcClient/certificattedRpc"
	"grpcClient/certificattedRpc/tokenCertication/token"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	auth := token.Authentication{
		User:     "wxning",
		Password: "gopher",
	}

	conn, err := grpc.Dial("localhost:1107", grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := certificattedRpc.NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &certificattedRpc.String{Value: " grpc"})
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
			if err := stream.Send(&certificattedRpc.String{Value: "grpc" + strconv.Itoa(num)}); err != nil {
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
