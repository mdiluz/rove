Rove
=====

An asynchronous nomadic game about exploring a planet as part of a loose 
community.

-------------------------------------------

## The Basics

### Core

Control a rover on the surface of the planet using a remote control interface.

Commands are sent and happen asynchronously, and the rover feeds back information about position and surroundings, as well as photos.

### Goal

To reach the pole.

### General

Movement is slow and sometimes dangerous. Hazards damage the rover.

Resources can be collected to fix and upgrade the rover.

Rovers recharge power during the day.

Enough collected resources allow you to create and fire a new rover a significant distance in any direction, losing control of the current one and leaving it dormant.

Finding a dormant rover gives you a choice - scrap it to gain minor resources, or fire it a distance just like a new rover, taking control of it.

“Dying” triggers a self destruct and fires a new basic robot in a random direction towards the equator

## Multiplayer

The planet itself and things that happen on it are persistent. Players can view each other, and use very rudimentary signals.

Dormant rovers store full history of travel, owners, and keep damage, improvements and resources.

Players have no other forms of direct communication.

Players can view progress of all rovers attached to their name.

Limit too many players in one location with a simple interference mechanic - only a certain density can exist at once to operate properly, additional players can’t move within range.

-------------------------------------------

### Implementation

Two functional parts

A server that receives the commands, sends out data, and handles interactions between players.

An app, or apps, that interface with the server to let you control and view rover information

-------------------------------------------

### To Solve

#### What kinds of progression/upgrades exist?
Needs a very simple set of rover internals defined, each of which can be upgraded.

#### How does the game encourage lateral movement?
Could simply be the terrain is constructed in very lateral ways, blocking progress frequently

#### How does the game encourage cooperation?
How exactly would a time delay mechanic enhance the experience?
Currently it’s just to make the multiplayer easier to use, and to make interactions a little more complicated. The game could limit the number of bytes (commands) you can send over time. 

#### How would the gameplay prevent griefing?

-------------------------------------------

### Key ideas left to integrate

Feeling “known for” something -  the person who did X thing. Limit number of  X things that can be done, possibly over time.

Build up a certain level of knowledge and ownership of a place, but then destroy it or give it up. Or build up a character and then leave it behind.

A significant aspect of failure - failing must be a focus of the game. Winning the game might actually be failing in some specific way.

A clear and well  defined investment vs. payoff curve.

Not an infinite game, let the game have a point where you’re done and can move on.

