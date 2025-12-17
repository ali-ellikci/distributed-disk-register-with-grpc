package leader

import (
	"fmt"
	"time"

	pb "distributed-disk-register-with-grpc/proto/family"
)

func StartFamilyPrinter(registry *NodeRegistry, self *pb.NodeInfo) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			members := registry.Snapshot()
			fmt.Println("======================================")
			fmt.Printf("Family at %s:%d (me)\n", self.Host, self.Port)
			fmt.Println("Members:")
			for _, n := range members {
				meMark := ""
				if n.Host == self.Host && n.Port == self.Port {
					meMark = " (me)"
				}
				fmt.Printf(" - %s:%d%s\n", n.Host, n.Port, meMark)
			}
			fmt.Println("======================================")
		}
	}()
}
