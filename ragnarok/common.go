package ragnarok

import (
	"fancubrain/common"
	"gocv.io/x/gocv"
)

var KnownObjects = []common.Object{
	{
		Name:  "Poring",
		Lower: gocv.Scalar{Val1: 91, Val2: 137, Val3: 199},
		Upper: gocv.Scalar{Val1: 212, Val2: 169, Val3: 255},
	},
	//{
	//	Name:  "Lunatic",
	//	Lower: gocv.Scalar{Val1: 200, Val2: 200, Val3: 200},
	//	Upper: gocv.Scalar{Val1: 255, Val2: 255, Val3: 255},
	//},
	//{
	//	Name:  "Cursor",
	//	Lower: gocv.Scalar{Val1: 91, Val2: 137, Val3: 199},
	//	Upper: gocv.Scalar{Val1: 212, Val2: 169, Val3: 255},
	//},
}
