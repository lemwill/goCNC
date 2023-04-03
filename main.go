// hello world
package main

// Calculate accelerations
// Transform to constant acceleration
// Interpolate to constant velocity
// GCode parser
func main() {

	// Create a Motion Planner
	// New machine configuration
	gcode_parser := NewGCodeParser()
	gcode_parser.Parse("test.gcode")

	machineConfiguration := NewMachineConfiguration(Vector3d{X: 100, Y: 100, Z: 100}, Vector3d{X: 100, Y: 100, Z: 100}, 0.1)
	motionPlanner := NewMotionPlanner(machineConfiguration)

	motionPlanner.command_list.AddMovement(NewLinearMovement(Vector3d{X: 8, Y: 0, Z: 0}, 50))
	motionPlanner.command_list.AddMovement(NewLinearMovement(Vector3d{X: 10, Y: 0, Z: 0}, 50))
	motionPlanner.command_list.AddMovement(NewLinearMovement(Vector3d{X: 0, Y: 0, Z: 0}, 10))
	motionPlanner.command_list.AddMovement(NewCircularMovement(Vector3d{X: 2, Y: 0, Z: 0}, Vector3d{X: 1, Y: 0, Z: 0}, 10, true, ZAxis))

	// Add the movement to the array

	//test := NewLinearMovement(Vector3d{X: 1, Y: 0, Z: 0}, 10)

	motionPlanner.run()

}
