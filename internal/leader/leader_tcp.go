package leader

import (
	"bufio"

	"log"
	"net"

	"distributed-disk-register-with-grpc/internal/common"
	"distributed-disk-register-with-grpc/internal/node"
	pb "distributed-disk-register-with-grpc/proto/family"
)

func StartLeaderTCPListener(registry *node.Registry, self *pb.NodeInfo) {
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

func handleTCPClient(conn net.Conn, registry *node.Registry, self *pb.NodeInfo) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Command parser
		cmd, err := common.ParseCommand(line)
		if err != nil {
			log.Printf("Failed to parse command: %v", err)
			continue
		}

		log.Printf("Received command: %+v", cmd)

		//Execute burada çağrılacak
		resp, _ := cmd.Execute()
		_, err = conn.Write([]byte(resp + "\n"))
		if err != nil {
			log.Printf("TCP client write error: %v", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("TCP client read error: %v", err)
	}
}
