# paxos

Golang implentation of [Paxos](https://pdos.csail.mit.edu/6.824/papers/paxos-simple.pdf) consensus algorithm.

Run multiple agents in different tabs to test the algorithm:

```
go run cmd/main.go run --port 1234 --peers 127.0.0.1:2345
go run cmd/main.go run --port 2345 --peers 127.0.0.1:1234
```
