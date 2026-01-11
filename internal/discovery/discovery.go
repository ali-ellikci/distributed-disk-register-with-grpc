package discovery

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const StartPort = 5555

func FindAvailablePort(host string) int32 {
	for port := StartPort; ; port++ {
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
			// Port boş, bunu kullan
			return int32(port)
		}

		conn.Close()
	}
}
