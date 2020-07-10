Rove
=====

Rove is an asynchronous nomadic game about exploring a planet as part of a loose community.

-------------------------------------------

## Core gameplay

Remotely explore the surface of a planet with an upgradable and customisable rover. Send commands to be executed asynchronously, view the rover's radar, and communicate and coordinate with other nearby rovers.

### Key Components

* Navigate an expansive world
* Collect resources to repair and upgrade
* Keep the rover batteries charged as you explore
* Help other players on their journey
* Explore north to discover more

-------------------------------------------

## Installing

On Ubuntu:
```
$ snap install rove
```

Elsewhere (with [go](https://golang.org/doc/install) installed)
```
go get github.com/mdiluz/rove
cd $GOPATH/src/github.com/mdiluz/rove/
make install
```

-------------------------------------------

### Implementation Details

`rove-server` hosts the game world and a gRPC server to allow users to interact from any client.

`rove` is a basic example command-line client that allows for simple play, to explore it's usage, see the output of `rove help`

-------------------------------------------

### "Find the fun" issues to solve

* What kinds of progression/upgrades exist?
* How does the game encourage cooperation?
* How would the gameplay prevent griefing?
* What drives the exploration?

-------------------------------------------

### Key ideas left to integrate

* Feeling “known for” something -  the person who did X thing. Limit number of  X things that can be done, possibly over time.
* A significant aspect of failure - failing must be a focus of the game. Winning the game might actually be failing in some specific way.
* A clear and well  defined investment vs. payoff curve.
* Not an infinite game, let the game have a point where you’re done and can move on.

