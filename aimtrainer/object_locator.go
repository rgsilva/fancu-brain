package aimtrainer

import (
	"fancubrain/common"
	"gocv.io/x/gocv"
	"image"
	"sort"
)

const rectResizeFactor = 0.05

type objectLocator struct {
	oldFrame gocv.Mat
	zeroes   gocv.Mat
}

func NewObjectLocator() common.ObjectLocator {
	return &objectLocator{
		oldFrame: gocv.NewMat(),
		zeroes:   gocv.Zeros(5, 5, gocv.MatTypeCV8U),
	}
}

func (l *objectLocator) Locate(frame *gocv.Mat) ([]common.FoundObject, error) {
	img := gocv.NewMat()
	defer img.Close()

	frame.CopyTo(&img)

	// Define the BGR threshold for the desired color range
	lowerThreshold := gocv.NewScalar(0, 70, 254, 0)
	upperThreshold := gocv.NewScalar(5, 85, 255, 0)

	// Create a mask based on the threshold
	mask := gocv.NewMat()
	defer mask.Close()
	gocv.InRangeWithScalar(img, lowerThreshold, upperThreshold, &mask)

	// Find contours in the mask
	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	// Mask all possible mobs.
	foundObjects := make([]common.FoundObject, 0)

	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		boundingRect := gocv.BoundingRect(contour)

		center := image.Point{
			X: boundingRect.Min.X + boundingRect.Dx()/2,
			Y: boundingRect.Min.Y + boundingRect.Dy()/2,
		}
		center = boundingRect.Min

		if center.X < MinX || center.Y < MinY {
			continue
		}

		alreadyFound := false
		for _, fo := range foundObjects {
			if isWithinRect(fo.Rect.BoundingRect, 100, center) {
				alreadyFound = true
				break
			}
		}
		if alreadyFound {
			continue
		}

		rect := gocv.RotatedRect{
			Center:       center,
			BoundingRect: boundingRect,
		}
		rect.Center.Y -= 10
		foundObjects = append(foundObjects, common.FoundObject{
			Object: common.Object{
				Name:  "Target",
				Lower: gocv.Scalar{},
				Upper: gocv.Scalar{},
			},
			Rect: rect,
		})
	}

	return foundObjects, nil
}

func isWithinRect(rect image.Rectangle, margin int, point image.Point) bool {
	return point.X >= rect.Min.X-margin && point.Y >= rect.Min.Y-margin && point.X <= rect.Max.X+margin && point.Y <= rect.Max.Y+margin
}

func sorFoundObjects(fms []common.FoundObject, sorters []func(left common.FoundObject, right common.FoundObject) bool) {
	for _, sorter := range sorters {
		sort.SliceStable(fms, func(i, j int) bool {
			return sorter(fms[i], fms[j])
		})
	}
}
