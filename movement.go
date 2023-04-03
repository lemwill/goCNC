package main

import (
	"fmt"
	"math"
)

// Create an interface for a movement object
type Movement interface {
	// Move the object
	get_start_position() Vector3d
	get_end_position() Vector3d
	get_start_velocity() float64
	get_end_velocity() float64
	get_target_velocity() float64
	set_start_position(Vector3d)
	set_end_position(Vector3d)
	set_start_velocity(float64)
	set_end_velocity(float64)
	set_target_velocity(float64)
	get_start_direction() Vector3d
	get_end_direction() Vector3d
	get_length() float64
}

// Linear movement
type LinearMovement struct {
	start_position  Vector3d
	end_position    Vector3d
	start_velocity  float64
	end_velocity    float64
	target_velocity float64
}

// Create a new linear movement
func NewLinearMovement(end_position Vector3d, target_velocity float64) *LinearMovement {
	return &LinearMovement{
		start_position:  Vector3d{X: 0, Y: 0, Z: 0},
		end_position:    end_position,
		start_velocity:  target_velocity,
		target_velocity: target_velocity,
		end_velocity:    0,
	}
}

// Get the start position
func (m *LinearMovement) get_start_position() Vector3d {
	return m.start_position
}

// Get the end position
func (m *LinearMovement) get_end_position() Vector3d {
	return m.end_position
}

// Get the start velocity
func (m *LinearMovement) get_start_velocity() float64 {
	return m.start_velocity
}

// Get the end velocity
func (m *LinearMovement) get_end_velocity() float64 {
	return m.end_velocity
}

// Get the target velocity
func (m *LinearMovement) get_target_velocity() float64 {
	return m.target_velocity
}

// Set the start position
func (m *LinearMovement) set_start_position(start_position Vector3d) {
	m.start_position = start_position
}

// Set the end position
func (m *LinearMovement) set_end_position(end_position Vector3d) {
	m.end_position = end_position
}

// Set the start velocity
func (m *LinearMovement) set_start_velocity(start_velocity float64) {
	m.start_velocity = start_velocity
}

// Set the end velocity
func (m *LinearMovement) set_end_velocity(end_velocity float64) {
	m.end_velocity = end_velocity
}

// Set the target velocity
func (m *LinearMovement) set_target_velocity(target_velocity float64) {
	m.target_velocity = target_velocity
}

// Get the start direction
func (m *LinearMovement) get_start_direction() Vector3d {
	return m.end_position.Subtract(m.start_position).Normalize()
}

// Get the end direction
func (m *LinearMovement) get_end_direction() Vector3d {
	return m.start_position.Subtract(m.end_position).Normalize()
}

// Get the length of the movement
func (m *LinearMovement) get_length() float64 {
	return m.end_position.Subtract(m.start_position).Length()
}

// Return a string representation of the movement
func (m *LinearMovement) String() string {
	return fmt.Sprintf("Linear movement:   Pos: %7.3f -> %7.3f  Velocity: %7.3f m/s -> %7.3f m/s -> %7.3f m/s", m.start_position, m.end_position, m.start_velocity, m.target_velocity, m.end_velocity)

}

// Circular movement
type CircularMovement struct {
	start_position  Vector3d
	end_position    Vector3d
	start_velocity  float64
	end_velocity    float64
	target_velocity float64
	clockwise       bool
	center_offset   Vector3d
	axis            Axis
}

// Create a new circular movement
func NewCircularMovement(end_position Vector3d, center_offset Vector3d, target_velocity float64, clockwise bool, axis Axis) *CircularMovement {
	return &CircularMovement{
		start_position:  Vector3d{X: 0, Y: 0, Z: 0},
		end_position:    end_position,
		start_velocity:  target_velocity,
		target_velocity: target_velocity,
		end_velocity:    0,
		clockwise:       clockwise,
		center_offset:   center_offset,
		axis:            axis,
	}
}

// Get the start position
func (m *CircularMovement) get_start_position() Vector3d {
	return m.start_position
}

// Get the end position
func (m *CircularMovement) get_end_position() Vector3d {
	return m.end_position
}

// Get the start velocity
func (m *CircularMovement) get_start_velocity() float64 {
	return m.start_velocity
}

// Get the end velocity
func (m *CircularMovement) get_end_velocity() float64 {
	return m.end_velocity
}

// Get the target velocity
func (m *CircularMovement) get_target_velocity() float64 {
	return m.target_velocity
}

// Set the start position
func (m *CircularMovement) set_start_position(start_position Vector3d) {
	m.start_position = start_position
}

// Set the end position
func (m *CircularMovement) set_end_position(end_position Vector3d) {
	m.end_position = end_position
}

// Set the start velocity
func (m *CircularMovement) set_start_velocity(start_velocity float64) {
	m.start_velocity = start_velocity
}

// Set the end velocity
func (m *CircularMovement) set_end_velocity(end_velocity float64) {
	m.end_velocity = end_velocity
}

// Set the target velocity
func (m *CircularMovement) set_target_velocity(target_velocity float64) {
	m.target_velocity = target_velocity
}

// Get the start direction
func (m *CircularMovement) get_start_direction() Vector3d {
	return m.center_offset.Rotate90(m.axis, !m.clockwise)
}

// Get the center
func (m *CircularMovement) get_center() Vector3d {
	return m.start_position.Add(m.center_offset)
}

// Get the end direction
func (m *CircularMovement) get_end_direction() Vector3d {
	return m.end_position.Subtract(m.get_center()).Normalize().Rotate90(m.axis, m.clockwise)
}

// Get the angle of the movement
func (m *CircularMovement) angle() float64 {
	angle := m.get_start_direction().AngleWith(m.get_end_direction())

	if m.clockwise {
		angle = 2*math.Pi - angle
	}

	return angle
}

// Get the length of the movement
func (m *CircularMovement) get_length() float64 {
	return m.center_offset.Length() * m.angle()
}

// Return a string representation of the movement
func (m *CircularMovement) String() string {
	return fmt.Sprintf("Circular Movement: Pos: %7.3f -> %7.3f  Velocity: %7.3f m/s -> %7.3f m/s -> %7.3f m/s", m.start_position, m.end_position, m.start_velocity, m.target_velocity, m.end_velocity)
}
