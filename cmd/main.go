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
	port := int32(5555) // İlk çalışan node lider kabul edilecek

	self := &pb.NodeInfo{
		Host: host,
		Port: port,
	}

	// Registry oluştur
	registry := leader.NewNodeRegistry()
	service := leader.NewFamilyService(registry, self)

	// gRPC server başlat
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("gRPC listen error: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterFamilyServiceServer(server, service)

	// Eğer lider node ise TCP 6666 listener başlat
	if port == 5555 {
		leader.StartLeaderTCPListener(registry, self)
	}

	// Var olan node'ları keşfet
	leader.DiscoverExistingNodes(host, self, registry)

	// Health check başlat
	leader.StartHealthChecker(registry, self)

	// Family durumunu periyodik yazdır
	leader.StartFamilyPrinter(registry, self)

	log.Printf("Node started on %s:%d\n", host, port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}
