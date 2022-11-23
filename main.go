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
	receiverPort = flag.Int("receiver", 9999, "The port of peer who receives messages and tokens from this process")
	p            = &peer{
		id:        9999,
		isElected: false,
		receiver:  nil,
		ctx:       nil,
	}
	criticalOperation = false
)

func main() {

	// The node 'decides' that it wants to do the critical operation
	go func() {
		for {
			if !criticalOperation {
				num := rand.Intn(5-1) + 1
				time.Sleep(time.Duration(num) * time.Second)
				criticalOperation = true
				log.Println("I would like to do the critical operation")
			}
		}
	}()

	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	log.Printf("my port is %v, sending to %v", *port, *receiverPort)

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

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	// Dial peer node
	conn, err := grpc.Dial(fmt.Sprintf(":%v", *receiverPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()
	p.receiver = critical.NewCriticalClient(conn)

	p.receiver.Election(ctx, &critical.Candidate{Id: p.id})

	// To ensure the program keeps running
	for {

	}
}

type peer struct {
	critical.UnimplementedCriticalServer
	id        int32
	isElected bool
	receiver  critical.CriticalClient
	ctx       context.Context
}

func (c *peer) Election(ctx context.Context, in *critical.Candidate) (*critical.Empty, error) {

	// Wait for a connection to the receiving node
	for {
		if p.receiver != nil {
			break
		}
	}

	// Compare the id of this node with the id from the election
	eId := in.Id
	if eId < p.id {
		p.receiver.Election(p.ctx, &critical.Candidate{Id: p.id})
		log.Println("sending election, greater id")
	} else if eId > p.id {
		p.receiver.Election(p.ctx, &critical.Candidate{Id: eId})
		log.Println("sending election, lesser id")
	} else if eId == p.id {
		log.Println("i am elected")
		if !p.isElected {
			p.isElected = true
			// Token is created, passed to peer node
			p.receiver.PassToken(p.ctx, &critical.Token{})
		}
	} else {
		return &critical.Empty{}, errors.New("error receiving message")
	}
	return &critical.Empty{}, nil
}

func (c *peer) PassToken(ctx context.Context, in *critical.Token) (*critical.Empty, error) {

	// The critical operation, accessing of the 'critical section'
	if criticalOperation {
		log.Printf("CRITICAL OPERATION IN PROCESS %v", *port)
		time.Sleep(3 * time.Second)
		criticalOperation = false
	}

	p.receiver.PassToken(c.ctx, &critical.Token{})
	return &critical.Empty{}, nil
}
