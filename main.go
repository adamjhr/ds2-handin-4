package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	critical "github.com/adamjhr/ds2-handin-4/grpc"
	"google.golang.org/grpc"
)

var (
	port         = flag.Int("port", 9999, "")
	recieverPort = flag.Int("reciever", 9999, "The port of peer who recieves messages and tokens from this process")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:        int32(*port),
		isElected: false,
		reciever:  nil,
		ctx:       ctx,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	critical.RegisterCriticalServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	fmt.Printf("Trying to dial: %v\n", recieverPort)
	conn, err := grpc.Dial(fmt.Sprintf(":%v", recieverPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()
	c := critical.NewCriticalClient(conn)
	p.reciever = c
}

type peer struct {
	critical.UnimplementedCriticalServer
	id        int32
	isElected bool
	reciever  critical.CriticalClient
	ctx       context.Context
}
