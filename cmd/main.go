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
	leaderPort := int32(5555)

	// 5555 boş mu diye bak
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", leaderPort))
	if err != nil {
		// 5555 DOLU -> leader var

		log.Println("[INFO] Leader already running, discovering cluster")
		port := discovery.FindAvailablePort(host)

		log.Printf("[ROLE] FOLLOWER on port %d\n", port)
		role := "follower"
		self := &pb.NodeInfo{
			Host: host,
			Port: port,
		}

		node.JoinCluster(host, leaderPort, self)

		registry := node.NewRegistry()
		registry.Add(self)

		grpcServer := grpc.NewServer()
		familyService := node.NewFamilyService(registry, self, role)
		pb.RegisterFamilyServiceServer(grpcServer, familyService)

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		grpcServer.Serve(lis)

		return
	}
	role := "leader"

	self := &pb.NodeInfo{
		Host: host,
		Port: leaderPort,
	}

	// Registry
	registry := node.NewRegistry()
	registry.Add(self)

	// 5555 BOŞ ->  leader ol
	log.Println("[ROLE] LEADER")

	grpcServer := grpc.NewServer()

	familyService := node.NewFamilyService(registry, self, role)
	pb.RegisterFamilyServiceServer(grpcServer, familyService)

	coordinator := leader.NewCoordinator(registry, self)

	leader.StartLeaderTCPListener(coordinator)
	node.StartHealthChecker(registry, self)
	leader.StartFamilyPrinter(registry, self)
	leader.StartMessagePrinter(coordinator)

	log.Printf("Leader started on %s:%d\n", host, leaderPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc server error: %v", err)
	}
}
