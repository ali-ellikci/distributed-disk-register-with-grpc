package leader

import (
	pb "distributed-disk-register-with-grpc/proto/family"
	"sync"
)

type NodeRegistry struct {
	mu    sync.Mutex
	nodes map[string]*pb.NodeInfo
}

func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		nodes: make(map[string]*pb.NodeInfo),
	}
}

func nodeKey(n *pb.NodeInfo) string {
	return n.Host + ":" + string(rune(n.Port))
}

func (r *NodeRegistry) Add(node *pb.NodeInfo) {
	r.mu.Lock()
	defer r.mu.Lock()
	r.nodes[nodeKey(node)] = node
}

func (r *NodeRegistry) AddAll(nodes []*pb.NodeInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, n := range nodes {
		r.nodes[nodeKey(n)] = n
	}
}

func (r *NodeRegistry) Remove(node *pb.NodeInfo) {
	r.mu.Lock()
	delete(r.nodes, nodeKey(node))
}

func (r *NodeRegistry) Snapshot() []*pb.NodeInfo {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]*pb.NodeInfo, 0, len(r.nodes))
	for _, n := range r.nodes {
		out = append(out, n)
	}
	return out
}
