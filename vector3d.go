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

// normalize
func (v Vector3d) normalize() Vector3d {
	length := v.length()
	return Vector3d{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

// length
func (v Vector3d) length() float64 {
	length_squared := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	return math.Sqrt(length_squared)
}

// dot product
func (v Vector3d) Dot(v2 Vector3d) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

// Angle with
func (v Vector3d) AngleWith(v2 Vector3d) float64 {
	ratio := v.Dot(v2) / (v.length() * v2.length())

	if ratio > 1 {
		ratio = 1
	}

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

func (v Vector3d) max() float64 {
	return math.Max(math.Max(v.X, v.Y), v.Z)
}

// Add two vectors
func (v Vector3d) Add(v2 Vector3d) Vector3d {
	return Vector3d{X: v.X + v2.X, Y: v.Y + v2.Y, Z: v.Z + v2.Z}
}

// subtract two vectors
func (v Vector3d) subtract(v2 Vector3d) Vector3d {
	return Vector3d{X: v.X - v2.X, Y: v.Y - v2.Y, Z: v.Z - v2.Z}
}

// Cross
func (v Vector3d) Cross(v2 Vector3d) Vector3d {
	return Vector3d{X: v.Y*v2.Z - v.Z*v2.Y, Y: v.Z*v2.X - v.X*v2.Z, Z: v.X*v2.Y - v.Y*v2.X}
}
