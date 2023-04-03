package main

import (
	"fmt"
	"math"
)

// Axis enum containing the three axis (AxisX, AxisY, AxisZ)
type Axis int

const (
	XAxis Axis = iota
	YAxis
	ZAxis
)

// Vector 3d struct
type Vector3d struct {
	X float64
	Y float64
	Z float64
}

// Custom string representation of the vector
func (v Vector3d) String() string {
	// Print with a padding of 10 characters
	return fmt.Sprintf("(%7.3f, %7.3f, %7.3f)", v.X, v.Y, v.Z)
}

// GoString is the same as String
func (v Vector3d) GoString() string {
	return v.String()
}

// Normalize
func (v Vector3d) Normalize() Vector3d {
	length := v.Length()
	return Vector3d{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

// Length
func (v Vector3d) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// dot product
func (v Vector3d) Dot(v2 Vector3d) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

// Angle with
func (v Vector3d) AngleWith(v2 Vector3d) float64 {
	ratio := v.Dot(v2) / (v.Length() * v2.Length())

	return math.Acos(ratio)
}

// Rotate 90 degrees around the given axis
func (v Vector3d) Rotate90(axis Axis, clockwise bool) Vector3d {
	switch axis {
	case XAxis:
		if clockwise {
			return Vector3d{X: v.X, Y: -v.Z, Z: v.Y}
		}
		return Vector3d{X: v.X, Y: v.Z, Z: -v.Y}
	case YAxis:
		if clockwise {
			return Vector3d{X: v.Z, Y: v.Y, Z: -v.X}
		}
		return Vector3d{X: -v.Z, Y: v.Y, Z: v.X}
	case ZAxis:
		if clockwise {
			return Vector3d{X: -v.Y, Y: v.X, Z: v.Z}
		}
		return Vector3d{X: v.Y, Y: -v.X, Z: v.Z}
	}

	return v
}

// Add two vectors
func (v Vector3d) Add(v2 Vector3d) Vector3d {
	return Vector3d{X: v.X + v2.X, Y: v.Y + v2.Y, Z: v.Z + v2.Z}
}

// Subtract two vectors
func (v Vector3d) Subtract(v2 Vector3d) Vector3d {
	return Vector3d{X: v.X - v2.X, Y: v.Y - v2.Y, Z: v.Z - v2.Z}
}
