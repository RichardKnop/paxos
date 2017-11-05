[1]: http://patreon_public_assets.s3.amazonaws.com/sized/becomeAPatronBanner.png
[2]: http://richardknop.com/images/btcaddress.png

## paxos

Golang implentation of [Paxos](https://pdos.csail.mit.edu/6.824/papers/paxos-simple.pdf) consensus algorithm.

[![Travis Status for RichardKnop/paxos](https://travis-ci.org/RichardKnop/merkle.svg?branch=master)](https://travis-ci.org/RichardKnop/paxos)
[![godoc for RichardKnop/paxos](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/RichardKnop/paxos)
[![goreportcard for RichardKnop/paxos](https://goreportcard.com/badge/github.com/RichardKnop/paxos)](https://goreportcard.com/report/RichardKnop/paxos)
[![codecov for RichardKnop/paxos](https://codecov.io/gh/RichardKnop/paxos/branch/master/graph/badge.svg)](https://codecov.io/gh/RichardKnop/paxos)
[![Codeship Status for RichardKnop/paxos](https://app.codeship.com/projects/1a959950-27be-0135-6f5e-7693ff866668/status?branch=master)](https://app.codeship.com/projects/223055)

[![Sourcegraph for RichardKnop/paxos](https://sourcegraph.com/github.com/RichardKnop/paxos/-/badge.svg)](https://sourcegraph.com/github.com/RichardKnop/paxos?badge)
[![Donate Bitcoin](https://img.shields.io/badge/donate-bitcoin-orange.svg)](https://richardknop.github.io/donate/)

---

* [First Steps](#first-steps)
* [Development](#development)
  * [Dependencies](#dependencies)
* [Supporting the project](#supporting-the-project)

### First Steps

Run multiple agents in different tabs to test the algorithm:

```
go run cmd/main.go run --port 1234 --peers 127.0.0.1:2345
go run cmd/main.go run --port 2345 --peers 127.0.0.1:1234
```

### Development

#### Dependencies

According to [Go 1.5 Vendor experiment](https://docs.google.com/document/d/1Bz5-UB7g2uPBdOx-rw5t9MxJwkfpx90cqG9AFL0JAYo), all dependencies are stored in the vendor directory. This approach is called `vendoring` and is the best practice for Go projects to lock versions of dependencies in order to achieve reproducible builds.

This project uses [dep](https://github.com/golang/dep) for dependency management. To update dependencies during development:

```sh
dep ensure
```


### Supporting the project

Become a patreon:

[![http://patreon.com/richardknop][1]](http://patreon.com/richardknop)

Or donate BTC to my wallet: `12iFVjQ5n3Qdmiai4Mp9EG93NSvDipyRKV`

![Donate BTC][2]