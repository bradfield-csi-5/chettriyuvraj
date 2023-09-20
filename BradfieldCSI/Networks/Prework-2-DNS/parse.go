package main

import "encoding/binary"

func ParseIPv4(data []byte) IPv4Packet {
	packet := IPv4Packet{}
	packet.Version = data[0] & 0xF0 >> 4                                /* bits 0-3 */
	packet.IHL = data[0] & 0x0F                                         /* bits 4-7 */
	packet.DSCP = (data[1] & 0xFC) >> 2                                 /* bits 8-13 */
	packet.ECN = data[1] & 0x03                                         /* bits 14-15 */
	packet.TotalLen = binary.BigEndian.Uint16(data[2:4])                /* bits 16-31 bytes 2,3 */ // (uint16(data[2]) & 0xFF00) & (0x00FF & uint16(data[3]))
	packet.Id = binary.BigEndian.Uint16(data[4:6])                      /* bytes 4,5 */
	packet.Flags = (data[6] & 0xE0) >> 5                                /* byte 6 bits 0-3 */
	packet.FragmentOffset = binary.BigEndian.Uint16(data[6:8]) & 0x1FFF /* byte 6, 7, last 13 bits */
	packet.TTL = data[8]                                                /* byte 8 */
	packet.Protocol = data[9]                                           /* byte 9 */
	packet.HeaderChecksum = binary.BigEndian.Uint16(data[10:12])        /* bytes 10,11 */
	packet.SourceIP = binary.BigEndian.Uint32(data[12:16])              /* bytes 12-15 */
	packet.DestIP = binary.BigEndian.Uint32(data[16:20])                /* bytes 16-19 */
	packet.Data = data[packet.IHL*4:]
	return packet
}

func ParseEthernet(data []byte) EthernetFrame {
	frame := EthernetFrame{MACdest: data[0:6], MACsource: data[6:12], TPID: data[12:14], TPIDExtended: nil, Data: nil}
	if frame.TPID[0] == 0x81 {
		frame.TPIDExtended = data[14:16]
		frame.Data = data[16:]
	} else {
		frame.Data = data[14:]
	}
	return frame
}
