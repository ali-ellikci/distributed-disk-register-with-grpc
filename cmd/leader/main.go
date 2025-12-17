package main

import (
	"fmt"
	"log"
	"net"

	"distributed-disk-register-with-grpc/internal/leader"
	pb "distributed-disk-register-with-grpc/proto/family"

	"google.golang.org/grpc"
)

func main() {
	host := "127.0.0.1"
	port := 5555

	self := &pb.NodeInfo{
		Host: host,
		Port: int32(port),
	}

	registry := leader.NewNodeRegistry()
	service := leader.NewFamilyService(registry, self)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	pb.RegisterFamilyServiceServer(server, service)

	leader.DiscoverExistingNodes(host, self, registry)
	leader.StartHealthChecker(registry, self)

	log.Printf("Node started on %s:%d\n", host, port)
	server.Serve(lis)
}
