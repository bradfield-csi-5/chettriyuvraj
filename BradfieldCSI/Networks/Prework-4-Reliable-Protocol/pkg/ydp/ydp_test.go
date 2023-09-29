package ydp

import (
	"bytes"
	"testing"

	"golang.org/x/sys/unix"
)

var TESTPROXYADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 64403, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}
var TESTCLIENTADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 5432, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}
var TESTSERVERADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 1234, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}

func TestYDPPacket(t *testing.T) {
	t.Run("test addr to byte slice", func(t *testing.T) {
		got, _ := AddrToByteSlice(TESTPROXYADDR)
		want := []byte{0x7f, 0x00, 0x00, 0x01, 0x00, 0x00, 0xFB, 0x93}
		// 1111 1011 1001 0011
		if !bytes.Equal(got, want) {
			t.Errorf("Error converting sockaddr to byte slice, Expected %x, Got %x", want, got)
		}
	})

	t.Run("test YDP packet message", func(t *testing.T) {
		got, _ := NewYDPPacket([]byte("Hi guys"), 0x01, TESTCLIENTADDR, TESTPROXYADDR)
		want := YDPPacket{Length: 7, Message: []byte("Hi guys")}

		if want.Length != got.Length || !bytes.Equal(want.Message, got.Message) {
			t.Errorf("Error generating YDP packet message field/length, Expected %x, Got %x", want.Message, got.Message)
		}
	})

	t.Run("test hash func", func(t *testing.T) {
		got, _ := Hash32([]byte{0x7f, 0x00, 0x00, 0x01, 0x00, 0x00, 0x15, 0x38, 0x7f, 0x00, 0x00, 0x01, 0x00, 0x00, 0xFB, 0x93})
		want := uint32(0x87779534)

		if got != want {
			t.Errorf("Error generating YDP packet message field/length, Expected %x, Got %x", want, got)
		}
	})

	t.Run("test YDP packet encoding to Binary", func(t *testing.T) {
		packet, _ := SampleYDPPacket([]byte("Hi guys"), 0x01, TESTCLIENTADDR, TESTPROXYADDR, 0x65157a3a)
		got, _ := packet.encode()
		want := []byte{0x55, 0x2b, 0xae, 0x4a, 0x01, 0x00, 0x07, 0x48, 0x69, 0x20, 0x67, 0x75, 0x79, 0x73}

		if !bytes.Equal(got, want) {
			t.Errorf("Error generating YDP packet message field/length, Expected %x, Got %x", want, got)
		}
	})

	t.Run("test decoding YDP packet from Binary", func(t *testing.T) {
		packet, _ := SampleYDPPacket([]byte("Hi guys"), 0x01, TESTCLIENTADDR, TESTPROXYADDR, 0x65157a3a)
		encodedPacket, _ := packet.encode()
		got, _ := DecodeYDPPacket(encodedPacket)
		want := packet

		if want.ID != got.ID || want.Flags != got.Flags || want.Length != got.Length || !bytes.Equal(want.Message, got.Message) {
			t.Errorf("Error decoding YDP packet from encoded bytes %x, Got %x", want, got)
		}
	})
}
