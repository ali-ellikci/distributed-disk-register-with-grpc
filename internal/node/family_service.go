package node

import (
	"context"
	"distributed-disk-register-with-grpc/internal/storage"
	"fmt"
	"log"

	pb "distributed-disk-register-with-grpc/proto/family"
)

type FamilyService struct {
	pb.UnimplementedFamilyServiceServer
	registry *Registry
	self     *pb.NodeInfo
	role     string // leader | follower
}

func NewFamilyService(registry *Registry, self *pb.NodeInfo, role string) *FamilyService {
	return &FamilyService{
		registry: registry,
		self:     self,
		role:     role,
	}
}

// ================= CLUSTER =================

func (s *FamilyService) Join(ctx context.Context, req *pb.NodeInfo) (*pb.FamilyView, error) {
	log.Printf("[JOIN] %s:%d\n", req.Host, req.Port)
	s.registry.Add(req)

	// Dead nodes listesinde varsa sil
	deadNodes := s.registry.GetDeadNodes()
	key := fmt.Sprintf("%s:%d", req.Host, req.Port)
	if deadNodes[key] {
		s.registry.RemoveDeadNode(req.Host, req.Port)
		log.Printf("[REVIVED] %s:%d (from dead nodes)\n", req.Host, req.Port)
	}

	return &pb.FamilyView{
		Members: s.registry.Snapshot(),
	}, nil
}

func (s *FamilyService) GetFamily(ctx context.Context, _ *pb.Empty) (*pb.FamilyView, error) {
	return &pb.FamilyView{
		Members: s.registry.Snapshot(),
	}, nil
}

// ================= STORAGE =================

// follower
func (s *FamilyService) Store(ctx context.Context, msg *pb.StoredMessage) (*pb.StoreResult, error) {
	if s.role != "follower" {
		return &pb.StoreResult{
			Success: false,
			Error:   "not follower",
		}, nil
	}

	fmt.Println("[INFO] Gelen yazılıyor ", msg)

	err := storage.WriteMessage(int(msg.Id), msg.Text)
	if err != nil {
		return &pb.StoreResult{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	log.Printf("[STORE] id=%d\n", msg.Id)
	return &pb.StoreResult{Success: true}, nil
}

func (s *FamilyService) Retrieve(ctx context.Context, id *pb.MessageId) (*pb.StoredMessage, error) {
	if s.role != "follower" {
		return nil, fmt.Errorf("not follower")
	}

	text, err := storage.ReadMessage(int(id.Id))
	if err != nil {
		return nil, err
	}

	return &pb.StoredMessage{
		Id:   id.Id,
		Text: text,
	}, nil
}
