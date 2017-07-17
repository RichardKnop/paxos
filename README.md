[![Travis Status for RichardKnop/paxos](https://travis-ci.org/RichardKnop/merkle.svg?branch=master)](https://travis-ci.org/RichardKnop/paxos)
[![godoc for RichardKnop/paxos](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/RichardKnop/paxos)
[![goreportcard for RichardKnop/paxos](https://goreportcard.com/badge/github.com/RichardKnop/paxos)](https://goreportcard.com/report/RichardKnop/paxos)
[![codecov for RichardKnop/paxos](https://codecov.io/gh/RichardKnop/paxos/branch/master/graph/badge.svg)](https://codecov.io/gh/RichardKnop/paxos)
[![Codeship Status for RichardKnop/paxos](https://app.codeship.com/projects/1a959950-27be-0135-6f5e-7693ff866668/status?branch=master)](https://app.codeship.com/projects/223055)

[![Sourcegraph for RichardKnop/paxos](https://sourcegraph.com/github.com/RichardKnop/paxos/-/badge.svg)](https://sourcegraph.com/github.com/RichardKnop/paxos?badge)
[![Donate Bitcoin](https://img.shields.io/badge/donate-bitcoin-orange.svg)](https://richardknop.github.io/donate/)

# paxos

Golang implentation of [Paxos](https://pdos.csail.mit.edu/6.824/papers/paxos-simple.pdf) consensus algorithm.

Run multiple agents in different tabs to test the algorithm:

```
go run cmd/main.go run --port 1234 --peers 127.0.0.1:2345
go run cmd/main.go run --port 2345 --peers 127.0.0.1:1234
```
