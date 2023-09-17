# Solution


## Initial observations

- Check file type using _file net.cap_: net.cap: pcap capture file, microsecond ts (little-endian) - version 2.4 (Ethernet, capture length 1514)
- Found pcap file format at https://datatracker.ietf.org/doc/id/draft-gharris-opsawg-pcap-00.html#name-general-file-structure
- _xxd -c 4 net.cap_ gives us data, matching with pcap format
- We get something that seems to match 
```
00000000: d4c3 b2a1  .... (Magic Number)
00000004: 0200 0400  .... (Version minor/major)
00000008: 0000 0000  ....
0000000c: 0000 0000  ....
00000010: ea05 0000  ....
00000014: 0100 0000  ....
00000018: 4098 d057  @..W
0000001c: 0a1f 0300  ....
00000020: 4e00 0000  N...
00000024: 4e00 0000  N...

packet data starts below.
```

- Next step is to parse packets and display data, something akin to Wireshark.


## Parse pcap

- Created a struct PacketRecord
- Parsed pcap header
- Parsed pcap as an array of PacketRecord[]
- Checked if sizes match by opening the packet in wireshark
- Seem to match, so the data should PROBABLY match, no test cases created here


## Parse Ethernet frame
- Use TDD
- 20th byte - 0-4th bits are not set, means 0 bytes of FCS is at the end of each Ethernet frame
- As a result will not consider FCS bits in our computation for now
- Reading up, also realize that physical layer data such as preamble et al are not recorded in this capture
- Thus, for Ethernet frames, we have structure as (Big-Endian, not sure why? is it by convention? Since pcap was little endian):
    - 6 bytes MAC destination
    - 6 bytes MAC source
    - 2 bytes check if tag 0x8100 or 0x8000 (EtherType field/TPID)
    - If 0x8100, we have 2 more bytes (TPID extended), else payload
    - Payload size -> Captured Len - (6 + 6 + 2 + TPIDExtendedLen)
    - Next (payload size) bytes are all payload of the Ethernet frame
- Steps:
    - Define an ethernet struct
    - Define a function to parse ethernet frame and its input/output
    - Directly check EtherType tag, i.e 12th and 13th byte, since we know its an ethernet frame, let's simply check 12th byte
    - If 0x81, else


## Conclusion

- Didn't really write down how I solved rest of the problem but some observations are as follows
- Followed TDD, had to get correct results using wireshark first
- Used Wireshark as a general guide to help me in the right direction
- Have ignored certain things at times Eg CRC, options section in packets, because of the fact that I saw the data beforehand in Wireshark
- Only checked out the solution on the website in glimpses
- Slight inconsistency in how I've used structs eg. TPID in Ethernet Frame could have been uint16 instead of []byte
- Also _PacketRecords_ have been parsed ugly-ly, would do it a little cleaner ideally