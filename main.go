package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	critical "github.com/adamjhr/ds2-handin-4/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port         = flag.Int("port", 9999, "")
	recieverPort = flag.Int("reciever", 9999, "The port of peer who recieves messages and tokens from this process")
	p            = &peer{
		id:        9999,
		isElected: false,
		reciever:  nil,
		ctx:       nil,
	}
	criticalOperation = false
)

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	log.Printf("my port is %v, sending to %v", *port, *recieverPort)

	ctx, cancel := context.WithCancel(context.Background())
	p.id = int32(*port)
	p.ctx = ctx
	defer cancel()

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	grpcServer := grpc.NewServer()
	critical.RegisterCriticalServer(grpcServer, p)

	fmt.Printf("Trying to dial: %v\n", recieverPort)
	conn, err := grpc.Dial(fmt.Sprintf(":%v", *recieverPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()
	p.reciever = critical.NewCriticalClient(conn)

	log.Println("About to call election")

	p.reciever.Election(ctx, &critical.Candidate{Id: p.id})

	log.Println("election called")

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	for {
		num := rand.Intn(5-1) + 1
		time.Sleep(time.Duration(num) * time.Second)
		criticalOperation = true
	}
}

type peer struct {
	critical.UnimplementedCriticalServer
	id        int32
	isElected bool
	reciever  critical.CriticalClient
	ctx       context.Context
}

func (c *peer) Election(ctx context.Context, in *critical.Candidate) (*critical.Empty, error) {

	log.Println("received election")

	log.Printf("Id of %v", in.Id)

	eId := in.Id
	if eId < p.id {
		p.reciever.Election(p.ctx, &critical.Candidate{Id: p.id})
		log.Println("sending election, greater id")
	} else if eId > p.id {
		p.reciever.Election(p.ctx, &critical.Candidate{Id: eId})
		log.Println("sending election, lesser id")
	} else if eId == p.id {
		log.Println("i am elected")
		if !p.isElected {
			p.isElected = true
			p.reciever.PassToken(p.ctx, &critical.Token{})
		}
	} else {
		return &critical.Empty{}, errors.New("error recieving message")
	}
	return &critical.Empty{}, nil
}

func (c *peer) PassToken(ctx context.Context, in *critical.Token) (*critical.Empty, error) {

	if criticalOperation {
		log.Printf("CRITICAL OPERATION IN PROCESS %v", port)
		criticalOperation = false
	}

	p.reciever.PassToken(c.ctx, &critical.Token{})
	return &critical.Empty{}, nil
}

// func (c *peer) Elected(ctx context.Context, in *critical.Queen) error {
// 	if !p.isElected {

// 	} else {
// 		// Peer may begin passing token
// 	}
// 	return nil
// }
