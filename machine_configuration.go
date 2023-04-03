package main

import "math"

type MachineConfiguration struct {
	max_acceleration         Vector3d
	max_velocity             Vector3d
	path_deviation_tolerance float64
}

// NewMachineConfiguration creates a new machine configuration
func NewMachineConfiguration(max_acceleration Vector3d, max_velocity Vector3d, path_deviation_tolerance float64) *MachineConfiguration {
	return &MachineConfiguration{max_acceleration: max_acceleration, max_velocity: max_velocity, path_deviation_tolerance: path_deviation_tolerance}
}

func (m *MachineConfiguration) get_max_acceleration_scalar() float64 {
	// Get the max acceleration scalar
	max_acceleration_scalar := m.max_acceleration.X
	if m.max_acceleration.Y > max_acceleration_scalar {
		max_acceleration_scalar = m.max_acceleration.Y
	}
	if m.max_acceleration.Z > max_acceleration_scalar {
		max_acceleration_scalar = m.max_acceleration.Z
	}
	return max_acceleration_scalar
}

// Get max acceleration for a given direction
func (m *MachineConfiguration) get_max_acceleration(start_direction Vector3d, end_direction Vector3d) float64 {
	// project the max acceleration vector onto the direction vector
	acceleration := math.Abs(m.max_acceleration.Dot(start_direction) / start_direction.Length())
	return acceleration
}
