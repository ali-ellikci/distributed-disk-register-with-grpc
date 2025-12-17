package leader

import (
	"context"
	"fmt"

	pb "distributed-disk-register-with-grpc/proto/family"

	"google.golang.org/grpc"
)

const StartPort = 5555

func DiscoverExistingNodes(host string, self *pb.NodeInfo, registry *NodeRegistry) {
	for port := StartPort; port < int(self.Port); port++ {
		addr := fmt.Sprintf("%s:%d", host, port)

		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			continue
		}

		client := pb.NewFamilyServiceClient(conn)

		view, err := client.Join(context.Background(), self)
		if err == nil {
			registry.AddAll(view.Members)
		}

		conn.Close()
	}
}
