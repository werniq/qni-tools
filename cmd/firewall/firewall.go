package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
)

func main() {
	// Specify the network interface to capture packets from (e.g., "eth0" or "en0")
	iface := "enp6s0"

	// Open the network interface for packet capturing
	handle, err := pcap.OpenLive(iface, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Loop to capture and process packets
	for packet := range packetSource.Packets() {
		// Print the packet details
		fmt.Println(packet)
	}
}
