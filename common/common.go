package common

import (
	"gocv.io/x/gocv"
)

type Config struct {
	Visualizer Visualizer
}

type ObjectLocator interface {
	Locate(frame *gocv.Mat) ([]FoundObject, error)
}

type Visualizer interface {
	Update(captured gocv.Mat, objects []VisObject)
}

type Object struct {
	Name  string
	Lower gocv.Scalar
	Upper gocv.Scalar
}

type FoundObject struct {
	Object Object
	Rect   gocv.RotatedRect
}

type VisObject struct {
	FoundObject

	Distance float64
	Selected bool
}
