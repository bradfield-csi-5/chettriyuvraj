# README

- Basic protocol implemented on top of UDP (called YDP)
    - Uses a 32-bit hash in the header as an ID, if packet with the same is received back from the receiver with ack bit set (unecessary actually), then the packet is assumed to be successfully transmitted.
    - Uses a timeout to wait for packet before it retransmits, only packet/ack for a single ID can be in transit at any given point.
    - Haven't implement a checksum, so corruption cannot be detected (except if there is corruption in the 'Hash' bits, which would lead to failed ACK and retranmission from sender)
    - Encoded as binary to send over protocol

- Unreliable proxy which periodically drops/corrupts bits used to demonstrate reliability of the protocol. Messages routed as follows: Client -> Proxy -> Server; Server->Proxy->Client
