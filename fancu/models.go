package fancu

type packet struct {
	mouseLeft   bool
	mouseMiddle bool
	mouseRight  bool
	mouseX      int8
	mouseY      int8
}

func (sd packet) ToBytes() []byte {
	data := make([]byte, 16)
	data[0] = 0xFF

	// HID bytes
	if sd.mouseLeft {
		data[1] |= 0b00000001
	}
	if sd.mouseMiddle {
		data[1] |= 0b00000010
	}
	if sd.mouseRight {
		data[1] |= 0b00000100
	}
	data[2] = byte(sd.mouseX)
	data[3] = byte(sd.mouseY)

	// TODO: keyboard.

	return data
}
