package main

import (
	"fmt"
	"math"
)

// Circular movement
type ArcMovement struct {
	start_position  Vector3d
	end_position    Vector3d
	start_velocity  float64
	end_velocity    float64
	target_velocity float64
	gcodeVelocity   float64
	clockwise       bool
	center_offset   Vector3d
	axis            Axis
}

// Create a new circular movement
func newArcMovement(end_position Vector3d, center_offset Vector3d, gcodeVelocity float64, clockwise bool, axis Axis) *ArcMovement {
	return &ArcMovement{
		start_position:  Vector3d{X: 0, Y: 0, Z: 0},
		end_position:    end_position,
		start_velocity:  gcodeVelocity,
		target_velocity: gcodeVelocity,
		gcodeVelocity:   gcodeVelocity,
		end_velocity:    0,
		clockwise:       clockwise,
		center_offset:   center_offset,
		axis:            axis,
	}
}

// Get the start position
func (m *ArcMovement) getStartPosition() Vector3d {
	return m.start_position
}

// Get the end position
func (m *ArcMovement) getEndPosition() Vector3d {
	return m.end_position
}

// Get the start velocity
func (m *ArcMovement) getStartVelocity() float64 {
	return m.start_velocity
}

// Get the end velocity
func (m *ArcMovement) getEndVelocity() float64 {
	return m.end_velocity
}

// Get the target velocity
func (m *ArcMovement) getTargetVelocity() float64 {
	return m.target_velocity
}

// Set the start position
func (m *ArcMovement) setStartPosition(start_position Vector3d) {
	m.start_position = start_position
}

// Set the end position
func (m *ArcMovement) setEndPosition(end_position Vector3d) {
	m.end_position = end_position
}

// Set the start velocity
func (m *ArcMovement) setStartVelocity(start_velocity float64) {
	m.start_velocity = start_velocity
}

// Set the end velocity
func (m *ArcMovement) setEndVelocity(end_velocity float64) {
	m.end_velocity = end_velocity
}

// Set the target velocity
func (m *ArcMovement) setTargetVelocity(target_velocity float64) {
	m.target_velocity = target_velocity
}

func (m *ArcMovement) getGcodeVelocity() float64 {
	return m.gcodeVelocity
}

func (m *ArcMovement) setGcodeVelocity(gcodeVelocity float64) {
	m.gcodeVelocity = gcodeVelocity
}

// Get the start direction
func (m *ArcMovement) getStartDirection() Vector3d {
	return m.center_offset.Rotate90(m.axis, !m.clockwise)
}

// Get the center
func (m *ArcMovement) getCenter() Vector3d {
	return m.start_position.Add(m.center_offset)
}

// Get the end direction
func (m *ArcMovement) getEndDirection() Vector3d {
	return m.end_position.subtract(m.getCenter()).normalize().Rotate90(m.axis, m.clockwise)
}

// Get the angle of the movement
func (m *ArcMovement) angle() float64 {
	angle := m.getStartDirection().AngleWith(m.getEndDirection())

	if m.clockwise {
		angle = 2*math.Pi - angle
	}

	return angle
}

// Limit the velocity to the maximum velocity
func (m *ArcMovement) limitVelocity(max_velocity Vector3d) {

	if m.axis != XAxis {
		if m.target_velocity > max_velocity.X {
			m.target_velocity = max_velocity.X
		}
	}
	if m.axis != YAxis {
		if m.target_velocity > max_velocity.Y {
			m.target_velocity = max_velocity.Y
		}
	}
	if m.axis != ZAxis {
		if m.target_velocity > max_velocity.Z {
			m.target_velocity = max_velocity.Z
		}
	}
}

func (m *ArcMovement) getMaxAcceleractionAlongMovement(maxAcceleration Vector3d) float64 {

	// Get the axis with the highest acceleration within maxAcceleration vector
	maxAccel := maxAcceleration.max()

	if m.axis != XAxis {
		if maxAcceleration.X < maxAccel {
			maxAccel = maxAcceleration.X
		}
	}
	if m.axis != YAxis {
		if maxAcceleration.Y < maxAccel {
			maxAccel = maxAcceleration.Y
		}
	}
	if m.axis != ZAxis {
		if maxAcceleration.Z < maxAccel {
			maxAccel = maxAcceleration.Z
		}
	}

	return maxAccel

}

// Get the length of the movement
func (m *ArcMovement) getLength() float64 {
	return m.center_offset.length() * m.angle()
}

// Return a string representation of the movement
func (m *ArcMovement) String() string {
	return fmt.Sprintf("Arc move:    Pos: %7.3f -> %7.3f  Velocity: %7.3f m/s -> %7.3f m/s -> %7.3f m/s", m.start_position, m.end_position, m.start_velocity, m.target_velocity, m.end_velocity)
}
