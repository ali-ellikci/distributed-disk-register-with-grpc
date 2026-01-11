package leader

import (
	"bufio"
	"log"
	"net"

	"distributed-disk-register-with-grpc/internal/common"
)

func StartLeaderTCPListener(coordinator *Coordinator) {
	go func() {
		listener, err := net.Listen("tcp", ":6666")
		if err != nil {
			log.Fatalf("TCP listener error: %v", err)
		}
		defer listener.Close()
		log.Printf("Leader listening on TCP :6666")

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("TCP accept error: %v", err)
				continue
			}
			go handleTCPClient(conn, coordinator)
		}
	}()
}

func handleTCPClient(conn net.Conn, coordinator *Coordinator) {
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
			conn.Write([]byte("ERROR\n"))
			continue
		}

		log.Printf("Received command: %+v", cmd)

		resp := coordinator.Handle(cmd)

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
