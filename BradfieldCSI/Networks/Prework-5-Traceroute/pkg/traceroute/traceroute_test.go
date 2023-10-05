package traceroute

import (
	"bytes"
	"testing"
)

func TestICMP(t *testing.T) {

	t.Run("Test ICMP Checksum", func(t *testing.T) {
		data := []struct {
			Packet      ICMPPacket
			ReqChecksum uint16
		}{
			{ICMPPacket{Type: 0x02, Code: 0x44, Data: []byte{0x01, 0x03, 0x04}}, 0xf8b8},
			{ICMPPacket{Type: 0x02, Code: 0x44, Data: []byte{0xF1, 0x03, 0x14}}, 0xf8b7},
		}

		for _, packet := range data {
			want := packet.ReqChecksum
			got := packet.Packet.ComputeChecksum()
			if want != got {
				t.Errorf("error computing checksum, wanted %x, got %x", want, got)
			}
		}
	})

	// 0x0244

	// 0x0244
	// 0x0103
	// ------
	// 0x0347
	// 0x0400
	// ------
	// 0x0747

	// 0x0244

	// 0x0244
	// 0xF103
	// ------
	// 0xF347
	// 0x1400
	// ------
	// 0x0748

	// 0000 0111 0100 1000
	// 1111 1000 1011 0111

	t.Run("Test icmp encode", func(t *testing.T) {
		packet := ICMPPacket{Type: 0x55, Code: 0x13, Checksum: uint16(20), ID: uint16(118), SequenceNo: uint16(66), Data: []byte{0x55, 0x32, 0x22, 0x11}}
		want := []byte{0x55, 0x13, 0x00, 0x14, 0x00, 0x76, 0x00, 0x42, 0x55, 0x32, 0x22, 0x11}
		got, _ := packet.Encode()

		if !bytes.Equal(want, got) {
			t.Errorf("error encoding checksum, wanted %d, got %d", want, got)
		}
	})

	t.Run("Test icmp decode", func(t *testing.T) {
		got, _ := DecodeICMPPacket([]byte{0x00, 0x00, 0xff, 0xd5, 0x00, 0x16, 0x00, 0x15, 0x00, 0x19})
		want := ICMPPacket{Type: 0x00, Code: 0x00, Checksum: 0xffd5, ID: 0x0016, SequenceNo: 0x0015, Data: []byte{0x00, 0x19}}

		if !bytes.Equal(got.Data, want.Data) || got.Type != want.Type || got.Code != want.Code || got.Checksum != want.Checksum || got.ID != want.ID || got.SequenceNo != want.SequenceNo {
			t.Errorf("error decoding checksum, wanted %d, got %d", want, got)
		}
	})

}

func TestIP(t *testing.T) {

	t.Run("Test IP Header Checksum", func(t *testing.T) {
		data := []struct {
			Packet      IPv4Packet
			ReqChecksum uint16
		}{
			{IPv4Packet{VersionAndIHL: 0x02, TOS: 0x44, TotalLen: 0x0103, ID: 0x0400}, 0x0747},
			{IPv4Packet{VersionAndIHL: 0x02, TOS: 0x44, TotalLen: 0xF103, ID: 0x1400}, 0x0748},
		}

		for i, packet := range data {
			want := packet.ReqChecksum
			got := packet.Packet.computeHeaderChecksum()
			if want != got {
				t.Errorf("error index %d computing checksum, wanted %d, got %d", i, want, got)
			}
		}
	})

	// 0x0244

	// 0x0244
	// 0x0103
	// ------
	// 0x0347
	// 0x0400
	// ------
	// 0x0747

	// 0x0244

	// 0x0244
	// 0xF103
	// ------
	// 0xF347
	// 0x1400
	// ------
	// 0x0748

	t.Run("Test IP encode", func(t *testing.T) {
		packet := NewIPv4(0, 0, 1, []byte{})
		want := []byte{0x45, 0x00, 0x00, 0x14, 0x00, 0x00, 0x40, 0x00, 0x01, 0x01, 0x86, 0x15, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		got, _ := packet.Encode()

		if !bytes.Equal(want, got) {
			t.Errorf("error encoding checksum, wanted %d, got %d", want, got)
		}
	})

	t.Run("Test ip decode", func(t *testing.T) {
		got, _ := DecodeIPv4Packet([]byte{0x45, 0x00, 0x00, 0x14, 0x00, 0x00, 0x40, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		want := IPv4Packet{
			VersionAndIHL:         0x45,
			TOS:                   0x00,
			TotalLen:              uint16(20),
			ID:                    0x0000,
			FlagsAndFragmentation: 0x4000, /* don't fragment */
			TTL:                   0x01,
			ULProto:               0x0001, /* ICMP */
			HeaderChecksum:        0x0000,
			SourceIP:              0,
			DestIP:                0,
		}

		if !bytes.Equal(got.Data, want.Data) || got.VersionAndIHL != want.VersionAndIHL || got.FlagsAndFragmentation != want.FlagsAndFragmentation || got.ID != want.ID {
			t.Errorf("error decoding checksum, wanted %d, got %d", want, got)
		}
	})

}
