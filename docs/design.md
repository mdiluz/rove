Rove
=====

An asynchronous nomadic game about exploring a planet as part of a loose community.

-------------------------------------------

## The Basics

### Core

Control a rover on the surface of the planet using a remote control interface.

Commands are sent and happen asynchronously, and the rover feeds back information about position and surroundings, as well as photos.

### General

Movement is slow and sometimes dangerous.

Resources can be collected, and rovers recharge power during the day.

Hazards damage the rover. Resources can be spent to repair.

Spend resources to create and spawn a new improved rover a significant distance away, leaving the current one dormant.

"Dying" leaves the current rover dormant and assigns the users a new rover.

Players can repair dormant rovers to gain control of them, taking on their improvements and inventory.

### Multiplayer

Players can see each other and use very rudimentary signals.

Dormant rovers store full history of travel and owners, as well as their improvements and resources.

-------------------------------------------

### Implementation

* A server that receives the commands, sends out data, and handles interactions between players.

* An app, or apps, that interface with the server to let you control and view rover information

-------------------------------------------

### To Solve

* What kinds of progression/upgrades exist?

* How does the game encourage cooperation?

* How would the gameplay prevent griefing?

* What drives the exploration?

-------------------------------------------

### Key ideas left to integrate

Feeling “known for” something -  the person who did X thing. Limit number of  X things that can be done, possibly over time.

A significant aspect of failure - failing must be a focus of the game. Winning the game might actually be failing in some specific way.

A clear and well  defined investment vs. payoff curve.

Not an infinite game, let the game have a point where you’re done and can move on.

