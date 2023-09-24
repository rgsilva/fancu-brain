package common

import (
	"fmt"
	"gocv.io/x/gocv"
)

type CaptureConfig struct {
	Device string
	Width  float64
	Height float64
}

type CaptureDevice struct {
	config  CaptureConfig
	capture *gocv.VideoCapture
}

func NewVideoCaptureDevice(config CaptureConfig) CaptureDevice {
	return CaptureDevice{
		config: config,
	}
}

func (v *CaptureDevice) Open() error {
	if v.capture != nil {
		return fmt.Errorf("device already open")
	}

	c, err := gocv.OpenVideoCapture(v.config.Device)
	if err != nil {
		return err
	}

	c.Set(gocv.VideoCaptureFrameWidth, v.config.Width)
	c.Set(gocv.VideoCaptureFrameHeight, v.config.Height)

	v.capture = c
	return nil
}

func (v *CaptureDevice) Capture(dst *gocv.Mat) (bool, error) {
	if v.capture == nil || !v.capture.IsOpened() {
		return false, fmt.Errorf("device is not open")
	}

	frame := gocv.NewMat()
	defer frame.Close()
	v.capture.Read(&frame)

	// Got no frame for some reason.
	if frame.Empty() {
		return false, nil
	}

	frame.CopyTo(dst)
	return true, nil
}

func (v *CaptureDevice) Close() error {
	if v.capture == nil {
		return fmt.Errorf("device is not open")
	}

	if err := v.capture.Close(); err != nil {
		return err
	}

	v.capture = nil
	return nil
}
