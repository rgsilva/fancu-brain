package common

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"time"
)

const windowWidth = 1920
const windowHeight = 1080

type screenVisualizer struct {
	window       *gocv.Window
	oldFrameTime int64
}

func NewScreenVisualizer() Visualizer {
	v := &screenVisualizer{}
	v.initialize()
	return v
}

func (v *screenVisualizer) initialize() {
	v.window = gocv.NewWindow("Preview")
	v.window.ResizeWindow(windowWidth, windowHeight)
}

func (v *screenVisualizer) Update(captured gocv.Mat, objects []VisObject) {
	if !captured.Empty() {
		fps := v.calculateFps()

		processed := gocv.NewMat()
		defer processed.Close()

		// Copy it, but darker
		gocv.ConvertScaleAbs(captured, &processed, 0.75, 0)
		//captured.CopyTo(&processed)

		// Draw FPS info.
		gocv.PutText(&processed, fmt.Sprintf("FPS: %.2f", fps), image.Point{X: 0, Y: 25}, gocv.FontHersheySimplex, 1.0, color.RGBA{G: 50, A: 255}, 2)

		// Show debug info.
		for i, fm := range objects {
			objectColor := color.RGBA{R: 50, G: 50, B: 50}
			if fm.Selected {
				objectColor = color.RGBA{R: 255, G: 0, B: 0}
			}

			// Draw a circle around the mob.
			gocv.Circle(&processed, fm.Rect.Center, 25, objectColor, 2)

			// Draw debug info around the mob.
			debug := []string{
				fmt.Sprintf("%d - %s", i, fm.Object.Name),
				fmt.Sprintf("Distance: %.02f", fm.Distance),
			}
			if fm.Selected {
				debug = append(debug, "*** SELECTED ***")
			}
			drawDebugInfo(&processed, image.Point{
				X: fm.Rect.Center.X,
				Y: fm.Rect.Center.Y + 50,
			}, objectColor, debug)
		}

		v.window.IMShow(processed)
	}
	v.window.WaitKey(1)
}

func (v *screenVisualizer) calculateFps() float64 {
	now := time.Now().UnixMilli()
	fps := 1000 / float64(now-v.oldFrameTime)
	v.oldFrameTime = now
	return fps
}

func drawDebugInfo(dst *gocv.Mat, position image.Point, color color.RGBA, lines []string) {
	yOffset := 0
	for _, line := range lines {
		ts := gocv.GetTextSize(line, gocv.FontHersheySimplex, 1.0, 2)
		gocv.PutText(dst, line, image.Point{
			X: position.X,
			Y: position.Y + yOffset,
		}, gocv.FontHersheySimplex, 1.0, color, 2)
		yOffset += ts.Y + 10
	}
}
