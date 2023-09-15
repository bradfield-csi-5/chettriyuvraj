package main

type PacketRecord struct {
	Timestamp   uint32
	Timestamp2  uint32
	CapturedLen uint32
	OriginalLen uint32
	Data        []byte
}

type EthernetFrame struct {
	MACdest      []byte
	MACsource    []byte
	TPID         []byte
	TPIDExtended []byte
	Data         []byte
}

// type IPPacket struct {
// 	Version []byte
// 	IHL     []byte
// 	DSCP    []byte
// 	ECN     []byte
// 	Total
// 	data []byte
// }
