package main

import (
	"fmt"
	"log"
	"net"

	"distributed-disk-register-with-grpc/internal/discovery"
	"distributed-disk-register-with-grpc/internal/leader"
	"distributed-disk-register-with-grpc/internal/node"
	pb "distributed-disk-register-with-grpc/proto/family"

	"google.golang.org/grpc"
)

func main() {
	host := "127.0.0.1"
	port := int32(5555)

	self := &pb.NodeInfo{
		Host: host,
		Port: port,
	}

	// Registry
	registry := node.NewRegistry()
	registry.Add(self)

	// 5555 boş mu diye bak
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// 5555 DOLU -> leader var
		log.Println("[INFO] Leader already running, discovering cluster")
		discovery.DiscoverExistingNodes(host, self, registry)
		log.Println("[INFO] Follower mode (şimdilik burada bitiyor)")
		return
	}

	// 5555 BOŞ ->  lider ol
	log.Println("[ROLE] LEADER")

	grpcServer := grpc.NewServer()

	familyService := node.NewFamilyService(registry, self)
	pb.RegisterFamilyServiceServer(grpcServer, familyService)

	leader.StartLeaderTCPListener(registry, self)
	node.StartHealthChecker(registry, self)
	leader.StartFamilyPrinter(registry, self)

	log.Printf("Leader started on %s:%d\n", host, port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc server error: %v", err)
	}
}
