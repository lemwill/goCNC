package main

import "math"

type MachineConfiguration struct {
	maxAcceleraction         Vector3d
	maxVelocity              Vector3d
	rapidVelocity            float64
	path_deviation_tolerance float64
}

// newMachineConfiguration creates a new machine configuration
func newMachineConfiguration(maxAcceleraction Vector3d, maxVelocity Vector3d, rapidVelocity float64, path_deviation_tolerance float64) *MachineConfiguration {
	return &MachineConfiguration{maxAcceleraction: maxAcceleraction, maxVelocity: maxVelocity, rapidVelocity: rapidVelocity, path_deviation_tolerance: path_deviation_tolerance}
}

func (m *MachineConfiguration) getMaxVelocity(direction Vector3d) float64 {

	// Project the max velocity vector onto the direction vector
	maxVelocity := math.Abs(m.maxVelocity.Dot(direction) / direction.length())
	return maxVelocity
}

func (m *MachineConfiguration) getMaxAcceleraction_scalar() float64 {
	// Get the max acceleration scalar
	maxAcceleraction_scalar := m.maxAcceleraction.X
	if m.maxAcceleraction.Y < maxAcceleraction_scalar {
		maxAcceleraction_scalar = m.maxAcceleraction.Y
	}
	if m.maxAcceleraction.Z < maxAcceleraction_scalar {
		maxAcceleraction_scalar = m.maxAcceleraction.Z
	}
	return maxAcceleraction_scalar
}

func (m *MachineConfiguration) getMaxAcceleractionForTwoVectors(start_direction Vector3d, end_direction Vector3d) float64 {
	maxAcceleration := m.maxAcceleraction.max()
	if start_direction.X != 0 || end_direction.X != 0 {
		if m.maxAcceleraction.X < maxAcceleration {
			maxAcceleration = m.maxAcceleraction.X
		}
	}
	if start_direction.Y != 0 || end_direction.Y != 0 {
		if m.maxAcceleraction.Y < maxAcceleration {
			maxAcceleration = m.maxAcceleraction.Y
		}
	}
	if start_direction.Z != 0 || end_direction.Z != 0 {
		if m.maxAcceleraction.Z < maxAcceleration {
			maxAcceleration = m.maxAcceleraction.Z
		}
	}

	return maxAcceleration
}

// Get max acceleration for a given direction
func (m *MachineConfiguration) getMaxAcceleraction(start_direction Vector3d, end_direction Vector3d) float64 {
	// project the max acceleration vector onto the direction vector
	acceleration := math.Abs(m.maxAcceleraction.Dot(start_direction) / start_direction.length())
	return acceleration
}
