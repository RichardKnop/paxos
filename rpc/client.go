package rpc

import (
	"log"
	"net/rpc"
	"time"

	"github.com/RichardKnop/paxos/paxos"
)

// Client ...
type Client struct {
	url string
}

// NewClient returns new instance of Client
func NewClient(url string) *Client {
	return &Client{url: url}
}

// GetName ...
func (c *Client) GetName() string {
	return c.url
}

// SendPrepare sends a prepare request to acceptor
func (c *Client) SendPrepare(proposal *paxos.Proposal) (*paxos.Proposal, error) {
	// fib := Fibonacci()

	reply, err := makeRequest(c.url, AcceptorReceivePrepare, proposal)

	return reply, err
}

// SendPropose sends a propose request to acceptor
func (c *Client) SendPropose(proposal *paxos.Proposal) (*paxos.Proposal, error) {
	// fib := Fibonacci()

	reply, err := makeRequest(c.url, AcceptorReceivePropose, proposal)

	return reply, err
}

// makeRequest is a generic function to make RPC calls
func makeRequest(to, serviceMethod string, proposal *paxos.Proposal) (*paxos.Proposal, error) {
	fib := Fibonacci()

	var (
		client *rpc.Client
		err    error
	)

	for {
		client, err = rpc.DialHTTP("tcp", to)
		if err != nil {
			// Use fibonacci sequence to space out retry attempts
			waitSec := fib()
			log.Printf("Failed to dial %s. Retrying in %ds", to, waitSec)
			<-time.After(time.Duration(waitSec) * time.Second)

			continue
		}

		break
	}

	var reply *paxos.Proposal
	if err = client.Call(serviceMethod, proposal, &reply); err != nil {
		return nil, err
	}

	return reply, nil
}
