package ydp

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

type YDPClient struct {
}

const CLIENTTIMEOUTSECS = 3

/* client socket is non blocking */
func (client *YDPClient) Send(message []byte, selfaddr unix.SockaddrInet4, dstaddr unix.SockaddrInet4, socket int) error {
	ackdata := make([]byte, 4096)

	/* Send, sleep till timeout waiting for ack, timeout, repeat */
	for {

		packet, err := SendYDP(message, selfaddr, dstaddr)
		if err != nil {
			return fmt.Errorf("error in client sending packet %x", err)
		}

		fmt.Printf("\n\n----->Client packet: %x\n----->Client Message: %q", packet, packet.Message)
		fmt.Printf("\n----->Sleeping for %d seconds...", CLIENTTIMEOUTSECS)

		time.Sleep(CLIENTTIMEOUTSECS * time.Second)

		/* Try to receive from non blocking socket - if someone keeps sending on the socket, will fail since it will keep receiving*/
		fmt.Printf("\n----->Client Resuming...")
		n, _, err := unix.Recvfrom(socket, ackdata, 0)
		if err != nil {
			if err == unix.EAGAIN {
				fmt.Printf("\n----->Nothing in socket on checking after timeout, retransmitting...")
				continue
			}
			return fmt.Errorf("error recieving ack packet from server %x", err)
		}

		ackpacket, err := DecodeYDPPacket(ackdata[:n])
		if err != nil { // unable to decode packet
			fmt.Printf("\n----->Error decoding received packet from server, continuing...")
			continue
		}

		fmt.Printf("\n----->Ackpacket received from server %x", ackpacket)

		if ackpacket.ID == packet.ID && ackpacket.Flags == 0x01 {
			fmt.Printf("\n----->Ack'd %x", packet.ID)
			return nil
		}

		fmt.Printf("\n----->Retrying client packet... %s ", packet.Message)
	}

}
