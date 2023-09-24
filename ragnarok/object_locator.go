package ragnarok

import (
	"fancubrain/common"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"log"
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
	iCapture := gocv.NewMat()
	defer iCapture.Close()
	iProcessed := gocv.NewMat()
	defer iProcessed.Close()

	//frame.CopyTo(&iCapture)
	gocv.CvtColor(*frame, &iCapture, gocv.ColorBGRToHSV)

	filteredFrame := gocv.NewMat()
	defer filteredFrame.Close()

	l.preFilter(frame, &filteredFrame)
	mobs := l.process(&filteredFrame, &iProcessed)

	return mobs, nil
}

func (l *objectLocator) preFilter(src *gocv.Mat, dst *gocv.Mat) {
	if l.oldFrame.Empty() {
		log.Printf("Ignoring first frame.")
		src.CopyTo(&l.oldFrame)
		src.CopyTo(dst)
		return
	}

	frame1 := gocv.NewMat()
	defer frame1.Close()
	frame2 := gocv.NewMat()
	defer frame2.Close()

	// Frame 1 <- previous frame
	// Frame 2 <- current frame
	// Previous frame <- current frame
	l.oldFrame.CopyTo(&frame1)
	src.CopyTo(&frame2)
	src.CopyTo(&l.oldFrame)

	// Preprocess frames
	gocv.CvtColor(frame1, &frame1, gocv.ColorRGBAToGray)
	gocv.GaussianBlur(frame1, &frame1, image.Point{X: 5, Y: 5}, 0, 0, gocv.BorderDefault)
	gocv.CvtColor(frame2, &frame2, gocv.ColorRGBAToGray)
	gocv.GaussianBlur(frame2, &frame2, image.Point{X: 5, Y: 5}, 0, 0, gocv.BorderDefault)

	// Calculate frame differences
	frameDiff := gocv.NewMat()
	defer frameDiff.Close()
	gocv.AbsDiff(frame1, frame2, &frameDiff)

	// Dilate the differences a bit
	gocv.Dilate(frameDiff, &frameDiff, l.zeroes)

	// Take only relevant differences
	thresholdFrame := gocv.NewMat()
	defer thresholdFrame.Close()
	gocv.Threshold(frameDiff, &thresholdFrame, 5, 255, gocv.ThresholdBinary)

	// Prepare the final image mask.
	mask := gocv.NewMatWithSize(src.Rows(), src.Cols(), gocv.MatTypeCV8U)
	defer mask.Close()

	// Get the contours of the areas with movement.
	contours := gocv.FindContours(thresholdFrame, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	for i := 0; i < contours.Size(); i++ {
		cpv := contours.At(i)
		if gocv.ContourArea(cpv) < 50 {
			continue
		}
		cr := gocv.BoundingRect(cpv)
		gocv.Rectangle(&mask, resizeRect(cr, rectResizeFactor), color.RGBA{uint8(255), uint8(255), uint8(255), uint8(255)}, -1)
	}

	// Blank the destination image.
	z := gocv.Zeros(src.Rows(), src.Cols(), src.Type())
	defer z.Close()
	z.CopyTo(dst)

	// Copy the final image with the mask.
	src.CopyToWithMask(dst, mask)
}

func (l *objectLocator) process(src *gocv.Mat, dst *gocv.Mat) []common.FoundObject {
	// Mask all possible mobs.
	foundObjects := make([]common.FoundObject, 0)
	for _, obj := range KnownObjects {
		mask := gocv.NewMat()

		// Mark the damn thing.
		gocv.InRangeWithScalar(*src, obj.Lower, obj.Upper, &mask)

		// If they are valid (have min area), write their centers down and calculate their distance.
		contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for i := 0; i < contours.Size(); i++ {
			cnt := contours.At(i)
			if area := gocv.ContourArea(cnt); area > 200 {
				// Write down the Object.
				rect := gocv.MinAreaRect(cnt)
				foundObjects = append(foundObjects, common.FoundObject{
					Object: obj,
					Rect:   rect,
				})
			}
		}

		mask.Close()
	}
	if len(foundObjects) > 0 {
		log.Printf("Found %d objects.\n", len(foundObjects))
	}

	// Clean up final image.
	zeros := gocv.NewMatWithSize(src.Rows(), src.Cols(), src.Type())
	zeros.CopyTo(dst)
	zeros.Close()

	// Copy relevant parts of the image.
	finalMask := gocv.NewMatWithSize(src.Rows(), src.Cols(), gocv.MatTypeCV8U)
	defer finalMask.Close()
	for _, fm := range foundObjects {
		gocv.Rectangle(&finalMask, resizeRect(fm.Rect.BoundingRect, rectResizeFactor), color.RGBA{255, 255, 255, 255}, -1)
		src.CopyToWithMask(dst, finalMask)
	}

	return foundObjects
}

func resizeRect(rect image.Rectangle, factor float64) image.Rectangle {
	return image.Rect(
		rect.Min.X-int(factor*float64(rect.Min.X)),
		rect.Min.Y-int(factor*float64(rect.Min.Y)),
		rect.Max.X+int(factor*float64(rect.Max.X)),
		rect.Max.Y+int(factor*float64(rect.Max.Y)),
	)
}
