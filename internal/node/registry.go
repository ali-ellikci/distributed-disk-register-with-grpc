package node

import (
	"fmt"
	"sync"

	pb "distributed-disk-register-with-grpc/proto/family"
)

type Registry struct {
	mu        sync.Mutex
	nodes     map[string]*pb.NodeInfo
	deadNodes map[string]bool
}

func NewRegistry() *Registry {
	return &Registry{
		nodes:     make(map[string]*pb.NodeInfo),
		deadNodes: make(map[string]bool),
	}
}

func nodeKey(n *pb.NodeInfo) string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

func (r *Registry) Add(node *pb.NodeInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[nodeKey(node)] = node
}

func (r *Registry) AddAll(nodes []*pb.NodeInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, n := range nodes {
		r.nodes[nodeKey(n)] = n
	}
}

func (r *Registry) Remove(node *pb.NodeInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := nodeKey(node)
	delete(r.nodes, key)
	r.deadNodes[key] = true
}

func (r *Registry) RemoveByAddress(host string, port int32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s:%d", host, port)
	delete(r.nodes, key)
	r.deadNodes[key] = true
}

func (r *Registry) Snapshot() []*pb.NodeInfo {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]*pb.NodeInfo, 0, len(r.nodes))
	for _, n := range r.nodes {
		out = append(out, n)
	}
	return out
}

func (r *Registry) GetDeadNodes() map[string]bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.deadNodes
}
