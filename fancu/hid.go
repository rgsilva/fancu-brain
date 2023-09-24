package fancu

import (
	"fmt"
	"log"
	"net"
)

type HID interface {
	Mouse(x, y int8, left, middle, right bool)
}

type hid struct {
	conn *net.UDPConn
}

func NewFANCU(address string, port uint16) HID {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}

	return &hid{
		conn: conn,
	}
}

func (s *hid) Mouse(x, y int8, left, middle, right bool) {
	p := packet{
		mouseLeft:   left,
		mouseMiddle: middle,
		mouseRight:  right,
		mouseX:      x,
		mouseY:      y,
	}
	pb := p.ToBytes()
	_, err := s.conn.Write(pb)
	if err != nil {
		log.Printf("Got an error writing data: %v", err)
	}
}
