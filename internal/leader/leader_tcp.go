package leader

import (
	"bufio"
	"context"
	pb "distributed-disk-register-with-grpc/proto/family"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func StartLeaderTCPListener(registry *NodeRegistry, self *pb.NodeInfo) {
	go func() {
		listener, err := net.Listen("tcp", ":6666")
		if err != nil {
			log.Fatalf("TCP listener error: %v", err)
		}
		defer listener.Close()
		log.Printf("Leader listening on TCP %s:6666", self.Host)

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("TCP accept error: %v", err)
				continue
			}
			go handleTCPClient(conn, registry, self)
		}
	}()
}

func handleTCPClient(conn net.Conn, registry *NodeRegistry, self *pb.NodeInfo) {
	defer conn.Close()
	reader := bufio.NewScanner(conn)

	for reader.Scan() {
		text := reader.Text()
		if text == "" {
			continue
		}
		log.Printf("Received from TCP: %s", text)

		msg := &pb.ChatMessage{
			Text:     text,
			FromHost: self.Host,
			FromPort: self.Port,
			// timestamp ekleyebilirsin
		}

		BroadcastToFamily(registry, self, msg)
	}

	if err := reader.Err(); err != nil {
		log.Printf("TCP client read error: %v", err)
	}
}

func BroadcastToFamily(registry *NodeRegistry, self *pb.NodeInfo, msg *pb.ChatMessage) {
	for _, n := range registry.Snapshot() {
		if n.Host == self.Host && n.Port == self.Port {
			continue
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", n.Host, n.Port), grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect %s:%d", n.Host, n.Port)
			continue
		}

		client := pb.NewFamilyServiceClient(conn)
		_, err = client.ReceiveChat(context.Background(), msg)
		if err != nil {
			log.Printf("Failed to send message to %s:%d", n.Host, n.Port)
		}

		conn.Close()
	}
}
