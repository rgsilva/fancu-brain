package ragnarok

import (
	"fancubrain/common"
	"gocv.io/x/gocv"
	"image"
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

	// Define the RGB threshold for the desired color range
	lowerThreshold := gocv.NewScalar(255, 255, 255, 0) // Example lower threshold (BGR order)
	upperThreshold := gocv.NewScalar(255, 255, 255, 0) // Example upper threshold (BGR order)

	// Define the RGB threshold for the second color range
	lowerThreshold2 := gocv.NewScalar(210, 170, 180, 0)
	upperThreshold2 := gocv.NewScalar(255, 240, 235, 0)

	// Create a mask based on the threshold
	mask := gocv.NewMat()
	defer mask.Close()
	gocv.InRangeWithScalar(img, lowerThreshold, upperThreshold, &mask)

	// Find contours in the mask
	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	// Draw bounding boxes around detected elements
	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		boundingRect := gocv.BoundingRect(contour)

		// Define a range of X pixels around the white element
		rangeX := 50 // Adjust the range as needed

		// Determine the region to search for the second color range
		searchRegion := image.Rect(
			max(boundingRect.Min.X, 0), boundingRect.Min.Y,
			min(boundingRect.Max.X+rangeX, img.Cols()), boundingRect.Max.Y,
		)

		// Create a sub-image representing the search region
		searchRegionImg := img.Region(searchRegion)

		// Create a mask based on the second threshold (greenish range)
		mask2 := gocv.NewMat()
		defer mask2.Close()
		gocv.InRangeWithScalar(searchRegionImg, lowerThreshold2, upperThreshold2, &mask2)

		// Check if the mask2 has non-zero pixels (greenish color range)
		if gocv.CountNonZero(mask2) > 0 {
			rect := gocv.RotatedRect{
				Center: boundingRect.Max,
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

		// Release the searchRegionImg to free resources
		searchRegionImg.Close()
	}

	return []common.FoundObject{}, nil
}
