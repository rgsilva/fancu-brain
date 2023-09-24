package aimtrainer

import (
	"fancubrain/common"
	"gocv.io/x/gocv"
)

func NewCursorLocator() common.ObjectLocator {
	return &cursorLocator{}
}

type cursorLocator struct {
	handTemplate   gocv.Mat
	attackTemplate gocv.Mat
}

func (l *cursorLocator) Locate(frame *gocv.Mat) ([]common.FoundObject, error) {
	img := gocv.NewMat()
	defer img.Close()

	frame.CopyTo(&img)

	// Define the BGR threshold for the desired color range
	lowerThreshold := gocv.NewScalar(0, 255, 0, 0)
	upperThreshold := gocv.NewScalar(30, 255, 30, 0)

	// Create a mask based on the threshold
	mask := gocv.NewMat()
	defer mask.Close()
	gocv.InRangeWithScalar(img, lowerThreshold, upperThreshold, &mask)

	// Find contours in the mask
	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		boundingRect := gocv.BoundingRect(contour)

		rect := gocv.RotatedRect{
			Center: boundingRect.Min,
		}
		return []common.FoundObject{
			{
				Object: common.Object{
					Name:  "Cursor",
					Lower: gocv.Scalar{},
					Upper: gocv.Scalar{},
				},
				Rect: rect,
			},
		}, nil
	}

	return []common.FoundObject{}, nil
}
