package main

import "fmt"

type CommandList struct {
	arr               []interface{}
	previous_position Vector3d
}

// Add a Movement
func (c *CommandList) AddMovement(movement Movement) {
	movement.set_start_position(c.previous_position)
	c.previous_position = movement.get_end_position()

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
