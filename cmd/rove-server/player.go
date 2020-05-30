package main

import "github.com/google/uuid"

type Player struct {
	id uuid.UUID
}

func NewPlayer() Player {
	return Player{
		id: uuid.New(),
	}
}
