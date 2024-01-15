package main

import (
	"fmt"
	"time"
)

// Message in TokenRing
type Message struct {
	Data      string
	Recipient int
	TTL       int
}

// Nodering is a node in the ring token
type Nodering struct {
	id     int
	input  chan Message
	output chan Message
}

// TokenRing network
type networkTokenring struct {
	Nodes []*Nodering
}

// Initializes n nodes in the token ring
func initialize(N int) *networkTokenring {
	network := &networkTokenring{}
	network.Nodes = make([]*Nodering, N)

	// Initialize nodes, run goroutines
	for i := 0; i < N; i++ {
		network.Nodes[i] = &Nodering{
			id:     i,
			input:  make(chan Message),
			output: make(chan Message),
		}
		go network.runNode(network.Nodes[i])
	}
	for i := 0; i < N; i++ {
		network.Nodes[i].connect(network.Nodes[(i+1)%N], network.Nodes[(i-1+N)%N])
	}

	return network
}

// Connects nodes (input output)
func (n *Nodering) connect(next, prev *Nodering) {
	go func() {
		for message := range n.input {
			// Send message to the next node
			next.output <- message
		}
	}()

	go func() {
		for message := range prev.output {
			// Receive message from the previous node
			n.input <- message
		}
	}()
}

// Node runs in the token network
func (network *networkTokenring) runNode(node *Nodering) {
	for message := range node.input {
		fmt.Printf("Remaining TTL for node %d: %d\n", node.id, message.TTL)

		// Decrease TTL
		message.TTL--

		// Check if TTL is still positive
		if message.TTL > 0 {
			// Resend the message
			node.output <- message
		}
	}
}

// Sends a message to all nodes in the network token ring
func (network *networkTokenring) SendMessageToAllNodes(data string, ttl int) {
	message := Message{
		Data: data,
		TTL:  ttl,
	}
	// Send message to the first node
	network.Nodes[0].input <- message
}

func main() {
	var N, ttl int
	var messageData string

	// Input number of nodes N
	fmt.Print("Enter the number of nodes: ")
	fmt.Scan(&N)

	// Initialize N nodes
	network := initialize(N)

	// Input TTL
	fmt.Print("Enter TTL (time to live): ")
	fmt.Scan(&ttl)

	// Input message
	fmt.Print("Enter the message: ")
	fmt.Scan(&messageData)

	// Send message to all nodes
	network.SendMessageToAllNodes(messageData, ttl)

	// Time delay for message send
	time.Sleep(5 * time.Second)
}
