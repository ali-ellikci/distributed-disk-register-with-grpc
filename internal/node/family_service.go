package node

import (
	"context"
	"log"

	pb "distributed-disk-register-with-grpc/proto/family"
)

type FamilyService struct {
	pb.UnimplementedFamilyServiceServer
	registry *Registry
	self     *pb.NodeInfo
}

func NewFamilyService(registry *Registry, self *pb.NodeInfo) *FamilyService {
	return &FamilyService{
		registry: registry,
		self:     self,
	}
}

func (s *FamilyService) Join(ctx context.Context, req *pb.NodeInfo) (*pb.FamilyView, error) {
	log.Printf("Node joined: %s:%d\n", req.Host, req.Port)
	s.registry.Add(req)

	return &pb.FamilyView{
		Members: s.registry.Snapshot(),
	}, nil
}

func (s *FamilyService) GetFamily(ctx context.Context, _ *pb.Empty) (*pb.FamilyView, error) {
	return &pb.FamilyView{
		Members: s.registry.Snapshot(),
	}, nil
}

func (s *FamilyService) ReceiveChat(ctx context.Context, msg *pb.ChatMessage) (*pb.Empty, error) {
	log.Println("💬 Incoming message")
	log.Printf("From: %s:%d\n", msg.FromHost, msg.FromPort)
	log.Printf("Text: %s\n", msg.Text)
	log.Printf("Timestamp: %d\n", msg.Timestamp)
	log.Println("--------------------------------")

	return &pb.Empty{}, nil
}
