package common

import (
	"gocv.io/x/gocv"
)

func NewTunnerVisualizer() Visualizer {
	v := &tunnerVisualizer{}
	v.initialize()
	return v
}

type tunnerVisualizer struct {
	window    *gocv.Window
	wControls *gocv.Window
	hh        *gocv.Trackbar
	hs        *gocv.Trackbar
	hv        *gocv.Trackbar
	lh        *gocv.Trackbar
	ls        *gocv.Trackbar
	lv        *gocv.Trackbar
}

func (v *tunnerVisualizer) initialize() {
	v.window = gocv.NewWindow("Preview")
	v.window.ResizeWindow(windowWidth, windowHeight)

	v.wControls = gocv.NewWindow("Controls")
	v.wControls.ResizeWindow(700, 100)
	v.lh = v.wControls.CreateTrackbar("Low H", 255)
	v.hh = v.wControls.CreateTrackbar("High H", 255)
	v.hh.SetPos(255)
	v.ls = v.wControls.CreateTrackbar("Low S", 255)
	v.hs = v.wControls.CreateTrackbar("High S", 255)
	v.hs.SetPos(255)
	v.lv = v.wControls.CreateTrackbar("Low V", 255)
	v.hv = v.wControls.CreateTrackbar("High V", 255)
	v.hv.SetPos(255)
}

func (v *tunnerVisualizer) Update(captured gocv.Mat, _ []VisObject) {
	processed := gocv.NewMat()
	defer processed.Close()

	captured.CopyTo(&processed)

	gocv.InRangeWithScalar(processed,
		gocv.Scalar{Val1: getPosFloat(v.lh), Val2: getPosFloat(v.ls), Val3: getPosFloat(v.lv)},
		gocv.Scalar{Val1: getPosFloat(v.hh), Val2: getPosFloat(v.hs), Val3: getPosFloat(v.hv)},
		&processed)

	v.window.IMShow(processed)
	v.window.WaitKey(1)
	v.wControls.WaitKey(1)
}

func getPosFloat(t *gocv.Trackbar) float64 {
	return float64(t.GetPos())
}
