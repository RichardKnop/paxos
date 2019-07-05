## paxos

Golang implentation of [Paxos](https://pdos.csail.mit.edu/6.824/papers/paxos-simple.pdf) consensus algorithm.

[![Travis Status for RichardKnop/paxos](https://travis-ci.org/RichardKnop/merkle.svg?branch=master)](https://travis-ci.org/RichardKnop/paxos)

---

* [First Steps](#first-steps)

### First Steps

This project is still under development. I have tried to keep the algorithm implementation completely decoupled so you can just import from `github.com/RichardKnop/paxos/paxos` and extend `Acceptor`, `Proposer` and `Learner` structs. 

In order to provide method of communication best suited for you, implement the `AcceptorClientInterface` interface which is then used by proposers to send request to acceptors.

The communicaton / networking between agents is something which is not relevant for the algorithm. I have written a simple RPC agent system for testing purposes though.

Run multiple test agents in different tabs to test the algorithm:

```
go run cmd/main.go run --port 1234 --peers 127.0.0.1:2345
go run cmd/main.go run --port 2345 --peers 127.0.0.1:1234
```
