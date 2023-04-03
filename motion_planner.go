package main

import (
	"fmt"
	"math"
)

// Create an array of movements
type MotionPlanner struct {
	command_list          CommandList
	machine_configuration *MachineConfiguration
}

// Create a new motion planner
func NewMotionPlanner(machineConfiguration *MachineConfiguration) *MotionPlanner {
	return &MotionPlanner{command_list: CommandList{}, machine_configuration: machineConfiguration}
}

// Calculate radius according to the machine configuration path deviation tolerance
func (m *MotionPlanner) calculate_radius(angle float64) float64 {
	radius := m.machine_configuration.path_deviation_tolerance * math.Sin(angle/2) / (1 - math.Sin(angle/2))
	return radius
}

func (m *MotionPlanner) get_max_corner_velocity(movement1 Movement, movement2 Movement) float64 {
	movement1_vector := movement1.get_end_direction()
	movement2_vector := movement2.get_start_direction()

	angle := movement1_vector.AngleWith(movement2_vector)

	radius := m.calculate_radius(angle)

	acceleration := m.machine_configuration.get_max_acceleration_scalar()

	max_cornering_velocity := math.Sqrt(acceleration * radius)

	return max_cornering_velocity
}

func (m *MotionPlanner) calculate_max_end_velocity(movement Movement) float64 {

	// Verify the maximum velocity at the end of the current movement
	v_initial := movement.get_start_velocity()
	distance := movement.get_length()

	// TODO: fix for arcs
	acceleration := m.machine_configuration.get_max_acceleration(movement.get_start_direction(), movement.get_end_direction())
	max_end_velocity := math.Sqrt(math.Pow(v_initial, 2) + acceleration*distance)

	return max_end_velocity
}

func (m *MotionPlanner) calculate_junction_velocity(movement Movement, next_movement Movement) float64 {

	// get the angle between the two movements
	max_cornering := m.get_max_corner_velocity(movement, next_movement)

	// Find the minimum between the target velocity, the next move start velocity, the max cornering velocity and the max end feedrate
	max_junction_velocity := math.Min(movement.get_target_velocity(), next_movement.get_start_velocity())
	max_junction_velocity = math.Min(max_junction_velocity, max_cornering)
	max_junction_velocity = math.Min(max_junction_velocity, movement.get_end_velocity())

	//fmt.Println("Max Junction Velocity: ", max_junction_velocity, " Max Cornering Velocity: ", max_cornering, " Max End Velocity: ", movement.get_end_velocity())

	return max_junction_velocity
}

func (m *MotionPlanner) calculate_max_start_velocity(movement Movement, max_junction_velocity float64) float64 {
	v_initial := max_junction_velocity
	distance := movement.get_length()
	acceleration := m.machine_configuration.get_max_acceleration(movement.get_start_direction(), movement.get_end_direction())
	max_start_velocity := math.Sqrt(math.Pow(v_initial, 2) + acceleration*distance)
	return max_start_velocity
}

func (m *MotionPlanner) run() {

	movements := m.command_list.GetMovementList()

	movements[0].set_start_velocity(0)

	i := 0
	for i < len(movements)-1 {

		movement := movements[i]
		next_movement := movements[i+1]

		// Calculate the maximum velocity at the end of the current movement given the start velocity, distance and acceleration
		max_end_velocity := m.calculate_max_end_velocity(movement)
		movement.set_end_velocity(max_end_velocity)

		// Calculate the maximum junction velocity given the machine deviation tolerance an angle between the two movements
		junction_velocity := m.calculate_junction_velocity(movement, next_movement)

		// Calculate the maximum start velocity given the calculated junction velocity, distance and acceleration
		max_start_feedrate := m.calculate_max_start_velocity(movement, junction_velocity)

		fmt.Println("[", i, "] ", movement)

		// Verify the velocity at the start of the movement is not too fast
		if max_start_feedrate < movement.get_start_velocity() {
			movement.set_start_velocity(max_start_feedrate)
			// TODO : Recalculate the previous movements
			fmt.Println("Start velocity too fast, recalculating previous movements")
			i -= 1
		} else {
			movement.set_end_velocity(junction_velocity)
			next_movement.set_start_velocity(junction_velocity)
			i = i + 1
		}
	}

	fmt.Println("=====================================")
	// Traverse the movement list
	for i := 0; i < len(movements); i++ {
		movement := movements[i]
		fmt.Println("[", i, "] ", movement)
	}
}
