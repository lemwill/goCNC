package main

import (
	"fmt"
	"math"
)

// Linear movement
type LinearMovement struct {
	start_position  Vector3d
	end_position    Vector3d
	start_velocity  float64
	end_velocity    float64
	target_velocity float64
	gcodeVelocity   float64
}

// Create a new linear movement
func newLinearMovement(end_position Vector3d, gcodeVelocity float64) *LinearMovement {
	return &LinearMovement{
		start_position:  Vector3d{X: 0, Y: 0, Z: 0},
		end_position:    end_position,
		start_velocity:  gcodeVelocity,
		target_velocity: gcodeVelocity,
		gcodeVelocity:   gcodeVelocity,
		end_velocity:    0,
	}
}

// Get the start position
func (m *LinearMovement) getStartPosition() Vector3d {
	return m.start_position
}

// Get the end position
func (m *LinearMovement) getEndPosition() Vector3d {
	return m.end_position
}

// Get the start velocity
func (m *LinearMovement) getStartVelocity() float64 {
	return m.start_velocity
}

// Get the end velocity
func (m *LinearMovement) getEndVelocity() float64 {
	return m.end_velocity
}

func (m *LinearMovement) getGcodeVelocity() float64 {
	return m.gcodeVelocity
}

func (m *LinearMovement) setGcodeVelocity(gcodeVelocity float64) {
	m.gcodeVelocity = gcodeVelocity
}

// Get the target velocity
func (m *LinearMovement) getTargetVelocity() float64 {
	return m.target_velocity
}

// Set the start position
func (m *LinearMovement) setStartPosition(start_position Vector3d) {
	m.start_position = start_position
}

// Set the end position
func (m *LinearMovement) setEndPosition(end_position Vector3d) {
	m.end_position = end_position
}

// Set the start velocity
func (m *LinearMovement) setStartVelocity(start_velocity float64) {
	m.start_velocity = start_velocity
}

// Set the end velocity
func (m *LinearMovement) setEndVelocity(end_velocity float64) {
	m.end_velocity = end_velocity
}

// Set the target velocity
func (m *LinearMovement) setTargetVelocity(target_velocity float64) {
	m.target_velocity = target_velocity
}

// Get the start direction
func (m *LinearMovement) getStartDirection() Vector3d {
	start_direction := m.end_position.subtract(m.start_position)
	if start_direction.length() == 0 {
		return Vector3d{X: 0, Y: 0, Z: 0}
	} else {
		return start_direction.normalize()
	}
}

// Get the end direction
func (m *LinearMovement) getEndDirection() Vector3d {
	return m.start_position.subtract(m.end_position).normalize()
}

// Get the length of the movement
func (m *LinearMovement) getLength() float64 {
	return m.end_position.subtract(m.start_position).length()
}

// Limit the velocity of the movement
func (m *LinearMovement) limitVelocity(maxVelocity Vector3d) {

	// Project the max velocity onto the direction of the movement
	direction := m.getStartDirection()

	maxVelocityAlongDirection := math.Abs(maxVelocity.Dot(direction) / direction.length())

	// Limit the velocity
	if m.target_velocity > maxVelocityAlongDirection {
		m.target_velocity = maxVelocityAlongDirection
	}
}

func (m *LinearMovement) getMaxAcceleractionAlongMovement(maxAcceleration Vector3d) float64 {
	// Project the max acceleration onto the direction of the movement
	direction := m.getStartDirection()

	maxAccelerationAlongDirection := math.Abs(maxAcceleration.Dot(direction) / direction.length())

	return maxAccelerationAlongDirection
}

// Return a string representation of the movement
func (m *LinearMovement) String() string {
	return fmt.Sprintf("Linear move: Pos: %7.3f -> %7.3f  Velocity: %7.3f m/s -> %7.3f m/s -> %7.3f m/s", m.start_position, m.end_position, m.start_velocity, m.target_velocity, m.end_velocity)

}
