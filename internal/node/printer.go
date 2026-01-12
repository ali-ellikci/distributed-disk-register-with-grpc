package node

import (
	"fmt"
	"time"
)

func StartMsgNumberPrinter(r *Registry) {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {

			fmt.Println("======================================")
			fmt.Println("Message Replication Status:")

			r.mu.Lock()
			fmt.Printf("  Message Number: %d\n", r.msgNumber)
			r.mu.Unlock()

			fmt.Println("======================================")
		}

	}()
}
