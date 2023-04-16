// hello world
package main

import (
	"fmt"
	"goCNC_protocol"
	"log"

	"google.golang.org/protobuf/proto"
)

// Transform to constant acceleration segments
// Interpolate to linear movements by steps of 1 ms
// Manage axis in Arc Movements
// Convert velocities to steps/s and acceleration to step/s2
// Make the algorithm real time (allow recalculation if feedrate is overridden during execution)

// Solution 1
//    Decelerate from 100% to 0
//    + Simple
//    - Deceleration can be faster
// Solution 2
//    Decelerate with the machine max acceleration
//    + Deceleration is optimal
//    - MCU needs to be able to calculate the deceleratio
// Solution 3
//    Move calculation to the machine

func main() {

	person := &goCNC_protocol.Person{
		Name: "John Doe",
		Age:  30,
	}

	data, err := proto.Marshal(person)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	fmt.Println(data)
	return

	// Create a Motion Planner
	// New machine configuration
	gcode_parser := newGCodeParser()
	parsedGCode := gcode_parser.fromFile("test.gcode")

	for _, command := range parsedGCode {
		fmt.Println(command.description)
	}

	machineConfiguration := newMachineConfiguration(
		Vector3d{X: 80, Y: 90, Z: 100},
		Vector3d{X: 50, Y: 40, Z: 100},
		100,
		0.1)

	motionPlanner := newMotionPlanner(machineConfiguration)

	motionPlanner.fromParsedGcode(parsedGCode)

	// Add the movement to the array

	//test := newLinearMovement(Vector3d{X: 1, Y: 0, Z: 0}, 10)

	motionPlanner.run()

}
