package main

import (
	"fancubrain/aimtrainer"
	"fancubrain/common"
	"fancubrain/fancu"
	"gocv.io/x/gocv"
	"image"
	"log"
	"time"
)

func main() {
	fc := fancu.NewFANCU("172.16.20.20", 9999)

	video := common.NewVideoCaptureDevice(common.CaptureConfig{
		Device: "0",
		Width:  1366,
		Height: 768,
	})
	if err := video.Open(); err != nil {
		log.Fatal(err)
	}

	// For color testing/tunning, enable this:
	//vis := common.NewTunnerVisualizer()

	// For preview window, enable this:
	vis := common.NewScreenVisualizer()

	// For ragnarok, enable these:
	//cursorLoc := ragnarok.NewCursorLocator()
	//objectLoc := ragnarok.NewObjectLocator()

	// For Aim Trainer, enable these:
	cursorLoc := aimtrainer.NewCursorLocator()
	objectLoc := aimtrainer.NewObjectLocator()

	var (
		currentCursor *common.FoundObject
		currentTarget *common.VisObject
	)

	for {
		frame := gocv.NewMat()

		// Grab a frame.
		if success, err := video.Capture(&frame); err != nil {
			log.Printf(err.Error())
			frame.Close()
			continue
		} else if !success {
			log.Printf("Got no frame, skipping.")
			frame.Close()
			continue
		}

		// Find the cursor.
		cursor, err := cursorLoc.Locate(&frame)
		if err != nil {
			log.Printf(err.Error())
		}
		if len(cursor) > 0 {
			currentCursor = &cursor[0]
		} else {
			currentCursor = nil
		}

		// If we don't have a cursor, don't do anything.
		if currentCursor == nil {
			vis.Update(frame, nil)
			continue
		}

		// Find all objects.
		found, err := objectLoc.Locate(&frame)
		if err != nil {
			log.Printf(err.Error())
		}

		allObjects := make([]common.VisObject, 0)
		for _, o := range found {
			allObjects = append(allObjects, common.VisObject{
				FoundObject: o,
			})
		}

		if len(allObjects) > 0 {
			// Sets the distance on all objects.
			for i := range allObjects {
				allObjects[i].Distance = common.PointDistance(allObjects[i].Rect.Center, currentCursor.Rect.Center)
			}

			// Sort them by distance to cursor.
			common.SortFoundObjects(allObjects, []func(left common.VisObject, right common.VisObject) bool{
				func(left common.VisObject, right common.VisObject) bool {
					return left.Distance < right.Distance
				},
			})

			// If there is a current target, we need to draw it too. Add it to the beginning of the slice.
			if currentTarget != nil {
				allObjects = append([]common.VisObject{*currentTarget}, allObjects...)
			}

			// Update screen.
			vis.Update(frame, allObjects)
		} else {
			// Update screen with no objects and carry one.
			vis.Update(frame, nil)
			continue
		}

		// If there is no current target, pick the closest one.
		if currentTarget == nil {
			currentTarget = &allObjects[0]
			currentTarget.Selected = true
		}

		// Distance to target on the X axis.
		target := currentTarget.Rect.Center
		distanceX := common.PointDistance(target, image.Point{
			X: currentCursor.Rect.Center.X,
			Y: target.Y,
		})

		// Distance to target on the Y axis.
		distanceY := common.PointDistance(target, image.Point{
			X: target.X,
			Y: currentCursor.Rect.Center.Y,
		})

		deltaX := deltaForDistance(distanceX)
		deltaY := deltaForDistance(distanceY)
		if deltaX == 0 && deltaY == 0 {
			fc.Mouse(0, 0, true, false, false)
			time.Sleep(100 * time.Millisecond)
			fc.Mouse(0, 0, false, false, false)

			// Resets current target after clicking on it.
			currentTarget = nil
		} else {
			mX := int8(0)
			if target.X > currentCursor.Rect.Center.X {
				mX = deltaX
			} else {
				mX = -deltaX
			}

			mY := int8(0)
			if target.Y > currentCursor.Rect.Center.Y {
				mY = deltaY
			} else {
				mY = -deltaY
			}
			fc.Mouse(mX, mY, false, false, false)
		}

		frame.Close()
	}
}

func deltaForDistance(distance float64) int8 {
	if distance <= 10 {
		return 0
	} else if distance <= 20 {
		return 3
	} else if distance <= 50 {
		return 5
	} else if distance <= 100 {
		return 10
	} else if distance <= 300 {
		return 15
	} else if distance <= 400 {
		return 20
	} else if distance <= 500 {
		return 25
	} else {
		return 35
	}
}
