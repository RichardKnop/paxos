package paxos

import (
	"fmt"

	"github.com/satori/go.uuid"
)

// Node represents a hostname and a port with a unique ID
type Node struct {
	ID   string
	Host string
	Port int
}

// NewNode returns a new node instance
func NewNode(id string, host string, port int) Node {
	if id == "" {
		id = uuid.NewV4().String()
	}
	if host == "" {
		host = "127.0.0.1"
	}
	return Node{
		ID:   id,
		Host: host,
		Port: port,
	}
}

// GetAddress concatenates host and port
func (n *Node) GetAddress() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

// String returns a human readable representation
func (n *Node) String() string {
	return n.GetAddress()
}
