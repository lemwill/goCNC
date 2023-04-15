package main

import (
	"fmt"
	"math"
)

// Create an array of movements
type MotionPlanner struct {
	commandList           CommandList
	machine_configuration *MachineConfiguration
}

// Create a new motion planner
func newMotionPlanner(machineConfiguration *MachineConfiguration) *MotionPlanner {
	return &MotionPlanner{commandList: CommandList{}, machine_configuration: machineConfiguration}
}

// Calculate radius according to the machine configuration path deviation tolerance
func (m *MotionPlanner) calculateRadius(angle float64) float64 {
	radius := m.machine_configuration.path_deviation_tolerance * math.Sin(angle/2) / (1 - math.Sin(angle/2))
	return radius
}

func (m *MotionPlanner) getMaxCornerVelocity(movement1 Movement, movement2 Movement) float64 {
	movement1_vector := movement1.getEndDirection()
	movement2_vector := movement2.getStartDirection()

	if movement1_vector.length() == 0 || movement2_vector.length() == 0 {
		return 0
	}

	angle := movement1_vector.AngleWith(movement2_vector)

	radius := m.calculateRadius(angle)

	acceleration := m.machine_configuration.getMaxAcceleractionForTwoVectors(movement1_vector, movement2_vector)

	max_cornering_velocity := math.Sqrt(acceleration * radius)

	return max_cornering_velocity
}

func (m *MotionPlanner) calculateMaxEndVelocity(movement Movement) float64 {

	// Verify the maximum velocity at the end of the current movement
	v_initial := movement.getStartVelocity()
	distance := movement.getLength()

	acceleration := movement.getMaxAcceleractionAlongMovement(m.machine_configuration.maxAcceleraction)
	max_end_velocity := math.Sqrt(math.Pow(v_initial, 2) + acceleration*distance)

	return max_end_velocity
}

func (m *MotionPlanner) calculateJunctionVelocity(movement Movement, next_movement Movement) float64 {

	// get the angle between the two movements
	max_cornering := m.getMaxCornerVelocity(movement, next_movement)

	// Find the minimum between the target velocity, the next move start velocity, the max cornering velocity and the max end feedrate
	max_junction_velocity := math.Min(movement.getTargetVelocity(), next_movement.getStartVelocity())
	max_junction_velocity = math.Min(max_junction_velocity, max_cornering)
	max_junction_velocity = math.Min(max_junction_velocity, movement.getEndVelocity())

	//fmt.Println("Max Junction Velocity: ", max_junction_velocity, " Max Cornering Velocity: ", max_cornering, " Max End Velocity: ", movement.getEndVelocity())

	return max_junction_velocity
}

func (m *MotionPlanner) calculateMaxStartVelocity(movement Movement, max_junction_velocity float64) float64 {
	v_initial := max_junction_velocity
	distance := movement.getLength()

	acceleration := movement.getMaxAcceleractionAlongMovement(m.machine_configuration.maxAcceleraction)
	max_start_velocity := math.Sqrt(math.Pow(v_initial, 2) + acceleration*distance)
	return max_start_velocity
}

func (m *MotionPlanner) fromParsedGcode(gcodeList []GCodeCommand) {
	// loop through the gcode commands
	for _, gcodeLine := range gcodeList {
		if gcodeLine.command == "G0" {
			// Create a new movement
			movement := newLinearMovement(Vector3d{X: gcodeLine.params["X"], Y: gcodeLine.params["Y"], Z: gcodeLine.params["Z"]}, m.machine_configuration.rapidVelocity)

			m.commandList.addMovement(movement)
		} else if gcodeLine.command == "G1" {
			// Create a new movement
			movement := newLinearMovement(
				Vector3d{X: gcodeLine.params["X"], Y: gcodeLine.params["Y"], Z: gcodeLine.params["Z"]},
				gcodeLine.params["F"])

			// Add the movement to the command list
			m.commandList.addMovement(movement)
		} else if gcodeLine.command == "G2" || gcodeLine.command == "G3" {
			clockwise := gcodeLine.command == "G2"
			// Create a new movement
			movement := newArcMovement(

				Vector3d{X: gcodeLine.params["X"], Y: gcodeLine.params["Y"], Z: gcodeLine.params["Z"]},
				Vector3d{X: gcodeLine.params["I"], Y: gcodeLine.params["J"], Z: gcodeLine.params["K"]},
				gcodeLine.params["F"],
				clockwise,
				ZAxis)

			// Add the movement to the command list
			m.commandList.addMovement(movement)
		}
	}
}

func (m *MotionPlanner) calculateFeedrateProfile(movement Movement) {

	// Calculate distance to accelerate to target feedrate

	maxAcceleration := movement.getMaxAcceleractionAlongMovement(m.machine_configuration.maxAcceleraction)
	accelerationDistance := (math.Pow(movement.getTargetVelocity(), 2) - math.Pow(movement.getStartVelocity(), 2)) / (2 * maxAcceleration)

	// Calculate distance to deccelerate to end feedrate
	deccelerationDistance := (math.Pow(movement.getTargetVelocity(), 2) - math.Pow(movement.getEndVelocity(), 2)) / (2 * maxAcceleration)

	// Calculate distance to accelerate and deccelerate
	accelerationDeccelerationDistance := accelerationDistance + deccelerationDistance

	// Verify that this is not longer than the move
	if accelerationDeccelerationDistance > movement.getLength() {
		// Calculate the new target feedrate, knowing the previous feedrate cannot be reached
		// v^2 = v_0^2 + 2ad
		feedrateDeltaDistance := math.Abs(math.Pow(movement.getEndVelocity(), 2)-math.Pow(movement.getStartVelocity(), 2)) / (2 * maxAcceleration)

		remainingDistance := math.Abs(movement.getLength()-feedrateDeltaDistance) / 2.0

		if movement.getEndVelocity() > movement.getStartVelocity() {
			accelerationDistance = remainingDistance + feedrateDeltaDistance
			deccelerationDistance = remainingDistance
		} else {
			accelerationDistance = remainingDistance
			deccelerationDistance = remainingDistance + feedrateDeltaDistance
		}

		// Calculate the new target feedrate
		reducedTargetFeedrate := math.Sqrt(2.0*maxAcceleration*accelerationDistance) + math.Pow(movement.getStartVelocity(), 2)

		movement.setTargetVelocity(reducedTargetFeedrate)
	}

	//constantAccelerationDistance := movement.getLength() - accelerationDistance - deccelerationDistance
	//fmt.Printf("Distance: %7.3f   Acc Dist: %7.3f   Const Dist: %7.3f   Dec Dist: %7.3f\r\n", movement.getLength(), accelerationDistance, constantAccelerationDistance, deccelerationDistance)
}

func (m *MotionPlanner) run() {

	movements := m.commandList.GetMovementList()

	movements[0].setStartVelocity(0)

	i := 0
	for i < len(movements)-1 {

		movement := movements[i]
		next_movement := movements[i+1]

		movement.limitVelocity(m.machine_configuration.maxVelocity)

		// Calculate the maximum velocity at the end of the current movement given the start velocity, distance and acceleration
		max_end_velocity := m.calculateMaxEndVelocity(movement)
		movement.setEndVelocity(max_end_velocity)

		// Calculate the maximum junction velocity given the machine deviation tolerance an angle between the two movements
		junction_velocity := m.calculateJunctionVelocity(movement, next_movement)

		// Calculate the maximum start velocity given the calculated junction velocity, distance and acceleration
		max_start_feedrate := m.calculateMaxStartVelocity(movement, junction_velocity)

		///fmt.Println("[", i, "] ", movement)

		// Verify the velocity at the start of the movement is not too fast
		if max_start_feedrate < movement.getStartVelocity() {
			movement.setStartVelocity(max_start_feedrate)
			// TODO : Recalculate the previous movements
			//fmt.Println("Start velocity too fast, recalculating previous movements")
			i -= 1
		} else {
			movement.setEndVelocity(junction_velocity)
			next_movement.setStartVelocity(junction_velocity)
			i = i + 1
		}
	}

	for i := 0; i < len(movements); i++ {
		movement := movements[i]

		m.calculateFeedrateProfile(movement)

	}

	fmt.Println("=====================================")
	// Traverse the movement list
	for i := 0; i < len(movements); i++ {
		fmt.Println("[", i, "] ", movements[i])
	}

}
