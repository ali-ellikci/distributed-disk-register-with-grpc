package leader

import (
	"distributed-disk-register-with-grpc/internal/common"
	"distributed-disk-register-with-grpc/internal/config"
	"distributed-disk-register-with-grpc/internal/node"
	"distributed-disk-register-with-grpc/internal/storage"
	pb "distributed-disk-register-with-grpc/proto/family"
	"fmt"
	"sort"
	"sync"
)

type Coordinator struct {
	registry         *node.Registry
	self             *pb.NodeInfo
	messageFollowers map[int][]int 
	followersMutex   sync.Mutex
	nextIndex        int        
	rrMutex          sync.Mutex 
}

func addFollowerForMessage(c *Coordinator, messageID int, followerPort int) {
	c.followersMutex.Lock()
	defer c.followersMutex.Unlock()
	if c.messageFollowers == nil {
		c.messageFollowers = make(map[int][]int)
	}
	c.messageFollowers[messageID] = append(c.messageFollowers[messageID], followerPort)
}

func NewCoordinator(registry *node.Registry, self *pb.NodeInfo) *Coordinator {
	return &Coordinator{
		registry:         registry,
		self:             self,
		messageFollowers: make(map[int][]int),
	}
}

func (c *Coordinator) Handle(cmd common.Command) string {
	switch v := cmd.(type) {

	case *common.SetCommand:

		return c.handleSet(v.ID, v.Text)

	case *common.GetCommand:
		return c.handleGet(v.ID)

	default:
		return "ERROR"
	}
}

func (c *Coordinator) handleSet(id int, text string) string {
	err := storage.WriteMessage(id, text)
	if err != nil {
		return ""
	}

	tolerance, err := config.LoadConfig()
	if err != nil {
		return "ERROR LOADING TOLERANCE"
	}
	members := c.pickMember(tolerance)

	fmt.Println("Members selected for replication:", members)

	for _, member := range members {
		if member.Port == c.self.Port && member.Host == c.self.Host {
			continue
		}
		fmt.Println("Sending SET to member:", member.Host, member.Port)
		addFollowerForMessage(c, id, int(member.Port))
		go SendSetToMember(member, id, text)
	}
	return "SET HANDLED BY LEADER"
}

func (c *Coordinator) handleGet(id int) string {
	msg, err := storage.ReadMessage(id)
	if err == nil {
		return msg
	}

	for i := 0; i < len(c.registry.Snapshot()); i++ {
		member := c.registry.Snapshot()[i]
		if member.Port == c.self.Port || member.Host == c.self.Host {
			continue
		}
		fmt.Println("Sending GET to member:", member.Host, member.Port)
		msg, err := SendGetToMember(member, id)
		if err == nil {
			return msg
		}

	}

	return "GET HANDLED BY LEADER"
}

func (c *Coordinator) pickMember(tolerance int) []*pb.NodeInfo {
	allNodes := c.registry.Snapshot()

	var candidates []*pb.NodeInfo
	for _, node := range allNodes {
		if node.Port != c.self.Port || node.Host != c.self.Host {
			candidates = append(candidates, node)
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Port < candidates[j].Port
	})

	totalCandidates := len(candidates)
	if totalCandidates == 0 {
		return nil
	}

	if totalCandidates < tolerance {
		tolerance = totalCandidates
	}

	c.rrMutex.Lock()
	startIndex := c.nextIndex
	c.nextIndex = (c.nextIndex + tolerance) % totalCandidates
	c.rrMutex.Unlock()

	var selected []*pb.NodeInfo
	for i := 0; i < tolerance; i++ {
		idx := (startIndex + i) % totalCandidates
		selected = append(selected, candidates[idx])
	}

	return selected
}
