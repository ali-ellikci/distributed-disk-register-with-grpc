package discovery

import (
	"context"
	"fmt"
	"time"

	"distributed-disk-register-with-grpc/internal/node"
	pb "distributed-disk-register-with-grpc/proto/family"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const StartPort = 5555

func DiscoverExistingNodes(host string, self *pb.NodeInfo, registry *node.Registry) {
	for port := StartPort; port < int(self.Port); port++ {
		addr := fmt.Sprintf("%s:%d", host, port)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		conn, err := grpc.DialContext(
			ctx,
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()

		if err != nil {
			continue
		}

		client := pb.NewFamilyServiceClient(conn)

		view, err := client.Join(ctx, self)
		if err == nil {
			registry.AddAll(view.Members)
		}

		conn.Close()
	}
}
