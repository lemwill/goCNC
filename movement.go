package main

// Create an interface for a movement object
type Movement interface {
	// Move the object
	getStartPosition() Vector3d
	getEndPosition() Vector3d
	getStartVelocity() float64
	getEndVelocity() float64
	getTargetVelocity() float64
	getGcodeVelocity() float64
	setGcodeVelocity(float64)
	setStartPosition(Vector3d)
	setEndPosition(Vector3d)
	setStartVelocity(float64)
	setEndVelocity(float64)
	setTargetVelocity(float64)
	getStartDirection() Vector3d
	getEndDirection() Vector3d
	getLength() float64
	limitVelocity(Vector3d)
	getMaxAcceleractionAlongMovement(Vector3d) float64
}
