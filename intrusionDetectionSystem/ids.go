package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"strings"
)

type IDS struct {
	rules map[string][]func(packet gopacket.Packet)
}

func NewIDS(interfaceName string, ruleBlock func(*IDS)) (*IDS, error) {
	ids := &IDS{rules: make(map[string][]func(packet gopacket.Packet))}
	ruleBlock(ids)

	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		ids.processPacket(packet)
	}

	return ids, nil
}

func (ids *IDS) processPacket(packet gopacket.Packet) {
	ethLayer := packet.LinkLayer()
	if ethLayer == nil {
		return
	}

	ethType := ethLayer.LinkFlow().EndpointType()
	header := ethType.String()

	if blocks, ok := ids.rules[header]; ok {
		for _, block := range blocks {
			block(packet)
		}
	}
}

func (ids *IDS) Rule(header string, block func(packet gopacket.Packet)) {
	if _, ok := ids.rules[header]; ok {
		ids.rules[header] = append(ids.rules[header], block)
	} else {
		ids.rules[header] = []func(packet gopacket.Packet){block}
	}
}

func main() {
	interfaceName := "enp6s0"
	ids, err := NewIDS(interfaceName, func(ids *IDS) {
		ids.Rule("Ethernet", func(packet gopacket.Packet) {
			log.Println("Ethernet packet detected")
		})

		ids.Rule("TCP", func(packet gopacket.Packet) {
			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer != nil {
				fmt.Println("TCP Layer detected")
				tcp, _ := tcpLayer.(*layers.TCP)

				// checks both the PSH and ACK TCP flags are set
				if tcp.PSH && tcp.ACK {
					out := tcp.Payload
					if strings.Contains(string(out), "cgi-bin/phf") {
						// verification passed
						fmt.Println("Flags are ok. Packet contains cgi-bin/phf")
					}
				}
			}
		})

		// the Push flag tells the receiver's network stack to "push" the data straight to the receiving socket,
		// and not to wait for any more packets before doing so
		ids.Rule("UDP", func(packet gopacket.Packet) {
			log.Printf("UDP packet detected")
		})
	})

	handle, err := pcap.OpenLive("enp6s0", 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatalln(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		ids.processPacket(packet)
	}
}
