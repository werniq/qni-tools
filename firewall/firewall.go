package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	_ "github.com/lib/pq"
	"log"
)

type NetworkLayerPayload struct {
	SrcMac string `json:"srcMac"`
	DstMac string `json:"dstMac"`
}

const (
	DbHost     = "localhost"
	DbPort     = 5432
	DbUsername = "YOUR_NAME"
	DbUserPw   = "YOUR_PW"
	DbName     = "DATABAS"
)

func init() {
	err := insertMacs()
	if err != nil {
		log.Fatalf("Error inserting mac addresses in database: %v\n", err)
	}
}

func insertMacs() error {
	macAddresses := []string{}
	connStr := fmt.Sprintf(`host=%s, port=%d, db_name=%s, user=%s, password=%s`,
		DbHost, DbPort, DbName, DbUsername, DbUserPw,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	for _, v := range macAddresses {
		res, err := db.Exec("INSERT INTO mac_addresses(mac) VALUES($1);", v)
		if err != nil {
			return err
		}

		if id, err := res.LastInsertId(); err != nil && id != 0 {
			return err
		}
	}

	return nil
}

func macInDb(mac string) (bool, error) {
	connStr := fmt.Sprintf(`host=%s, port=%d, db_name=%s, user=%s, password=%s`,
		DbHost, DbPort, DbName, DbUsername, DbUserPw,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return false, err
	}

	res := db.QueryRow("SELECT * FROM mac_addresses WHERE mac = $1", mac)
	if res.Err() != nil {
		return false, err
	}

	var result struct {
		id  int    `json:"id"`
		mac string `json:"mac"`
	}
	if err := res.Scan(&result); err != nil {
		return false, err
	}

	if result.mac == "" {
		return false, nil
	}

	return true, nil
}

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

		var netPayload *NetworkLayerPayload
		//ethLayer := packet.LinkLayer()
		//if ethLayer == nil {
		//	log.Println("No Ethernet layer found in the packet")
		//	continue
		//}

		//ethLayer.LinkFlow().Dst()

		if err = json.Unmarshal(packet.NetworkLayer().LayerContents(), &netPayload); err != nil {
			ok, err := macInDb(netPayload.DstMac)
			if err != nil {
				log.Printf("Error querying to the database: %v\n", err.Error())
				continue
			}

			if !ok {
				log.Printf("Given Mac is not in the database")

				continue
			}
		}
	}
}
