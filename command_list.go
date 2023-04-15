package main

import "fmt"

/*
	type Command struct {
		command        string
		parameters     map[string]float64
		subCommandList []interface{}
	}
*/
type CommandList struct {
	arr               []interface{}
	previous_position Vector3d
}

// Add a Movement
func (c *CommandList) addMovement(movement Movement) {
	movement.setStartPosition(c.previous_position)
	c.previous_position = movement.getEndPosition()

	c.arr = append(c.arr, movement)

}

// New Command List
func NewCommandList() *CommandList {
	return &CommandList{arr: []interface{}{}, previous_position: Vector3d{X: 0, Y: 0, Z: 0}}
}

// Return Movement list
func (c *CommandList) GetMovementList() []Movement {
	var movement_list []Movement
	for _, command := range c.arr {
		switch command.(type) {
		case Movement:
			movement_list = append(movement_list, command.(Movement))
		}
	}

	return movement_list
}

// Print the command list
func (c *CommandList) print() {
	for _, command := range c.arr {
		switch command.(type) {
		case Movement:
			fmt.Println(command.(Movement))
		}
	}

	fmt.Println("Done ===========================")
}
