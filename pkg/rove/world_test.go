package rove

import (
	"testing"

	"github.com/mdiluz/rove/pkg/maths"
	"github.com/mdiluz/rove/proto/roveapi"
	"github.com/stretchr/testify/assert"
)

func TestNewWorld(t *testing.T) {
	// Very basic for now, nothing to verify
	world := NewWorld(4)
	if world == nil {
		t.Error("Failed to create world")
	}
}

func TestWorld_CreateRover(t *testing.T) {
	world := NewWorld(8)
	a, err := world.SpawnRover("")
	assert.NoError(t, err)
	b, err := world.SpawnRover("")
	assert.NoError(t, err)

	// Basic duplicate check
	if a == b {
		t.Errorf("Created identical rovers")
	} else if len(world.Rovers) != 2 {
		t.Errorf("Incorrect number of rovers created")
	}
}

func TestWorld_GetRover(t *testing.T) {
	world := NewWorld(4)
	a, err := world.SpawnRover("")
	assert.NoError(t, err)

	rover, err := world.GetRover(a)
	assert.NoError(t, err, "Failed to get rover attribs")
	assert.NotZero(t, rover.Range, "Rover should not be spawned blind")
	assert.Contains(t, rover.Logs[len(rover.Logs)-1].Text, "created", "Rover logs should contain the creation")
}

func TestWorld_DestroyRover(t *testing.T) {
	world := NewWorld(1)
	a, err := world.SpawnRover("")
	assert.NoError(t, err)
	b, err := world.SpawnRover("")
	assert.NoError(t, err)

	err = world.destroyRover(a)
	assert.NoError(t, err, "Error returned from rover destroy")

	// Basic duplicate check
	if len(world.Rovers) != 1 {
		t.Error("Too many rovers left in world")
	} else if _, ok := world.Rovers[b]; !ok {
		t.Error("Remaining rover is incorrect")
	}
}

func TestWorld_GetSetMovePosition(t *testing.T) {
	world := NewWorld(4)
	a, err := world.SpawnRover("")
	assert.NoError(t, err)

	pos := maths.Vector{
		X: 0.0,
		Y: 0.0,
	}

	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	newPos, err := world.RoverPosition(a)
	assert.NoError(t, err, "Failed to set position for rover")
	assert.Equal(t, pos, newPos, "Failed to correctly set position for rover")

	b := roveapi.Bearing_North
	newPos, err = world.TryMoveRover(a, b)
	assert.NoError(t, err, "Failed to set position for rover")
	pos.Add(maths.Vector{X: 0, Y: 1})
	assert.Equal(t, pos, newPos, "Failed to correctly move position for rover")

	rover, err := world.GetRover(a)
	assert.NoError(t, err, "Failed to get rover information")
	assert.Contains(t, rover.Logs[len(rover.Logs)-1].Text, "moved", "Rover logs should contain the move")

	// Place a tile in front of the rover
	world.Atlas.SetObject(maths.Vector{X: 0, Y: 2}, Object{Type: roveapi.Object_RockLarge})
	newPos, err = world.TryMoveRover(a, b)
	assert.NoError(t, err, "Failed to move rover")
	assert.Equal(t, pos, newPos, "Failed to correctly not move position for rover into wall")
}

func TestWorld_RadarFromRover(t *testing.T) {
	// Create world that should have visible walls on the radar
	world := NewWorld(2)
	a, err := world.SpawnRover("")
	assert.NoError(t, err)
	b, err := world.SpawnRover("")
	assert.NoError(t, err)

	// Warp the rovers into position
	bpos := maths.Vector{X: -3, Y: -3}
	world.Atlas.SetObject(bpos, Object{Type: roveapi.Object_ObjectUnknown})
	assert.NoError(t, world.WarpRover(b, bpos), "Failed to warp rover")
	world.Atlas.SetObject(maths.Vector{X: 0, Y: 0}, Object{Type: roveapi.Object_ObjectUnknown})
	assert.NoError(t, world.WarpRover(a, maths.Vector{X: 0, Y: 0}), "Failed to warp rover")

	radar, objs, err := world.RadarFromRover(a)
	assert.NoError(t, err, "Failed to get radar from rover")
	fullRange := 4 + 4 + 1
	assert.Equal(t, fullRange*fullRange, len(radar), "Radar returned wrong length")
	assert.Equal(t, fullRange*fullRange, len(objs), "Radar returned wrong length")

	// Test the expected values
	assert.Equal(t, roveapi.Object_RoverLive, objs[1+fullRange])
	assert.Equal(t, roveapi.Object_RoverLive, objs[4+4*fullRange])

	// Check the radar results are stable
	radar1, objs1, err := world.RadarFromRover(a)
	assert.NoError(t, err)
	radar2, objs2, err := world.RadarFromRover(a)
	assert.NoError(t, err)
	assert.Equal(t, radar1, radar2)
	assert.Equal(t, objs1, objs2)
}

func TestWorld_RoverDamage(t *testing.T) {
	world := NewWorld(2)
	acc, err := world.Accountant.RegisterAccount("tmp")
	assert.NoError(t, err)
	a, err := world.SpawnRover(acc.Name)
	assert.NoError(t, err)

	pos := maths.Vector{
		X: 0.0,
		Y: 0.0,
	}

	world.Atlas.SetObject(pos, Object{Type: roveapi.Object_ObjectUnknown})
	err = world.WarpRover(a, pos)
	assert.NoError(t, err, "Failed to set position for rover")

	info, err := world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")

	world.Atlas.SetObject(maths.Vector{X: 0.0, Y: 1.0}, Object{Type: roveapi.Object_RockLarge})

	vec, err := world.TryMoveRover(a, roveapi.Bearing_North)
	assert.NoError(t, err, "Failed to move rover")
	assert.Equal(t, pos, vec, "Rover managed to move into large rock")

	newinfo, err := world.GetRover(a)
	assert.NoError(t, err, "couldn't get rover info")
	assert.Equal(t, info.Integrity-1, newinfo.Integrity, "rover should have lost integrity")
	assert.Contains(t, newinfo.Logs[len(newinfo.Logs)-1].Text, "collision", "Rover logs should contain the collision")

	// Keep moving to damage the rover
	for i := 0; i < info.Integrity-1; i++ {
		vec, err := world.TryMoveRover(a, roveapi.Bearing_North)
		assert.NoError(t, err, "Failed to move rover")
		assert.Equal(t, pos, vec, "Rover managed to move into large rock")
	}

	// Rover should have been destroyed now
	_, err = world.GetRover(a)
	assert.Error(t, err)

	_, obj := world.Atlas.QueryPosition(info.Pos)
	assert.Equal(t, roveapi.Object_RoverDormant, obj.Type)
}

func TestWorld_Daytime(t *testing.T) {
	world := NewWorld(1)

	a, err := world.SpawnRover("")
	assert.NoError(t, err)

	// Remove rover charge
	rover := world.Rovers[a]
	rover.Charge = 0
	world.Rovers[a] = rover

	// Try and recharge, should work
	_, err = world.RoverRecharge(a)
	assert.NoError(t, err)
	assert.Equal(t, 1, world.Rovers[a].Charge)

	// Loop for half the day
	for i := 0; i < world.TicksPerDay/2; i++ {
		assert.True(t, world.Daytime())
		world.Tick()
	}

	// Remove rover charge again
	rover = world.Rovers[a]
	rover.Charge = 0
	world.Rovers[a] = rover

	// Try and recharge, should fail
	_, err = world.RoverRecharge(a)
	assert.NoError(t, err)
	assert.Equal(t, 0, world.Rovers[a].Charge)

	// Loop for half the day
	for i := 0; i < world.TicksPerDay/2; i++ {
		assert.False(t, world.Daytime())
		world.Tick()
	}
}

func TestWorld_Broadcast(t *testing.T) {
	world := NewWorld(8)

	a, err := world.SpawnRover("")
	assert.NoError(t, err)

	b, err := world.SpawnRover("")
	assert.NoError(t, err)

	// Warp rovers near to eachother
	world.Atlas.SetObject(maths.Vector{X: 0, Y: 0}, Object{Type: roveapi.Object_ObjectUnknown})
	world.Atlas.SetObject(maths.Vector{X: 1, Y: 0}, Object{Type: roveapi.Object_ObjectUnknown})
	assert.NoError(t, world.WarpRover(a, maths.Vector{X: 0, Y: 0}))
	assert.NoError(t, world.WarpRover(b, maths.Vector{X: 1, Y: 0}))

	// Broadcast from a
	assert.NoError(t, world.RoverBroadcast(a, []byte{'A', 'B', 'C'}))

	// Check if b heard it
	ra, err := world.GetRover(a)
	assert.NoError(t, err)
	assert.Equal(t, ra.MaximumCharge-1, ra.Charge, "A should have used a charge to broadcast")
	assert.Contains(t, ra.Logs[len(ra.Logs)-1].Text, "ABC", "Rover B should have heard the broadcast")

	// Check if a logged it
	rb, err := world.GetRover(b)
	assert.NoError(t, err)
	assert.Contains(t, rb.Logs[len(rb.Logs)-1].Text, "ABC", "Rover A should have logged it's broadcast")

	// Warp B outside of the range of A
	world.Atlas.SetObject(maths.Vector{X: ra.Range, Y: 0}, Object{Type: roveapi.Object_ObjectUnknown})
	assert.NoError(t, world.WarpRover(b, maths.Vector{X: ra.Range, Y: 0}))

	// Broadcast from a again
	assert.NoError(t, world.RoverBroadcast(a, []byte{'X', 'Y', 'Z'}))

	// Check if b heard it
	ra, err = world.GetRover(b)
	assert.NoError(t, err)
	assert.NotContains(t, ra.Logs[len(ra.Logs)-1].Text, "XYZ", "Rover B should not have heard the broadcast")

	// Check if a logged it
	rb, err = world.GetRover(a)
	assert.NoError(t, err)
	assert.Contains(t, rb.Logs[len(rb.Logs)-1].Text, "XYZ", "Rover A should have logged it's broadcast")

	// Warp B outside of the range of A
	world.Atlas.SetObject(maths.Vector{X: ra.Range + 1, Y: 0}, Object{Type: roveapi.Object_ObjectUnknown})
	assert.NoError(t, world.WarpRover(b, maths.Vector{X: ra.Range + 1, Y: 0}))

	// Broadcast from a again
	assert.NoError(t, world.RoverBroadcast(a, []byte{'H', 'J', 'K'}))

	// Check if b heard it
	ra, err = world.GetRover(b)
	assert.NoError(t, err)
	assert.NotContains(t, ra.Logs[len(ra.Logs)-1].Text, "HJK", "Rover B should have heard the broadcast")

	// Check if a logged it
	rb, err = world.GetRover(a)
	assert.NoError(t, err)
	assert.Contains(t, rb.Logs[len(rb.Logs)-1].Text, "HJK", "Rover A should have logged it's broadcast")
}

func TestWorld_Sailing(t *testing.T) {
	world := NewWorld(8)
	world.Tick()                       // One initial tick to set the wind direction the first time
	world.Wind = roveapi.Bearing_North // Set the wind direction to north

	name, err := world.SpawnRover("")
	assert.NoError(t, err)

	// Warp the rover to 0,0 after clearing it
	world.Atlas.SetObject(maths.Vector{X: 0, Y: 0}, Object{Type: roveapi.Object_ObjectUnknown})
	assert.NoError(t, world.WarpRover(name, maths.Vector{X: 0, Y: 0}))

	s, err := world.RoverToggle(name)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.SailPosition_CatchingWind, s)

	b, err := world.RoverTurn(name, roveapi.Bearing_North)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.Bearing_North, b)

	// Clear the space to the north
	world.Atlas.SetObject(maths.Vector{X: 0, Y: 1}, Object{Type: roveapi.Object_ObjectUnknown})

	// Tick the world and check we've moved not moved
	world.Tick()
	info, err := world.GetRover(name)
	assert.NoError(t, err)
	assert.Equal(t, maths.Vector{Y: 0}, info.Pos)

	// Loop a few more times
	for i := 0; i < ticksPerNormalMove-2; i++ {
		world.Tick()
		info, err := world.GetRover(name)
		assert.NoError(t, err)
		assert.Equal(t, maths.Vector{Y: 0}, info.Pos)
	}

	// Now check we've moved (after the TicksPerNormalMove number of ticks)
	world.Tick()
	info, err = world.GetRover(name)
	assert.NoError(t, err)
	assert.Equal(t, maths.Vector{Y: 1}, info.Pos)

	// Reset the world ticks back to stop any wind changes etc.
	world.CurrentTicks = 1

	// Face the rover south, into the wind
	b, err = world.RoverTurn(name, roveapi.Bearing_South)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.Bearing_South, b)

	// Tick a bunch, we should never move
	for i := 0; i < ticksPerNormalMove*2; i++ {
		world.Tick()
		info, err := world.GetRover(name)
		assert.NoError(t, err)
		assert.Equal(t, maths.Vector{Y: 1}, info.Pos)
	}

	// Reset the world ticks back to stop any wind changes etc.
	world.CurrentTicks = 1
	world.Wind = roveapi.Bearing_SouthEast // Set up a south easternly wind

	// Turn the rover perpendicular
	b, err = world.RoverTurn(name, roveapi.Bearing_NorthEast)
	assert.NoError(t, err)
	assert.Equal(t, roveapi.Bearing_NorthEast, b)

	// Clear a space
	world.Atlas.SetObject(maths.Vector{X: 1, Y: 2}, Object{Type: roveapi.Object_ObjectUnknown})

	// Now check we've moved immediately
	world.Tick()
	info, err = world.GetRover(name)
	assert.NoError(t, err)
	assert.Equal(t, maths.Vector{X: 1, Y: 2}, info.Pos)
}
