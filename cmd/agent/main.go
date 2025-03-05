package main

import (
	"log"

	"github.com/popvictor123/distributed-calc/internal/agent"
)

func main() {
	log.Println("Starting agent...")
	agent := agent.NewAgent()
	agent.Start()
}
