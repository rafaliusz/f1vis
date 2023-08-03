package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/gopacket"
	pcap "github.com/google/gopacket/pcap"
)

type PacketType uint8

const (
	F1Motion              PacketType = 0
	F1Session             PacketType = 1
	F1LapData             PacketType = 2
	F1Event               PacketType = 3
	F1Participants        PacketType = 4
	F1CarSetups           PacketType = 5
	F1CarTelemetry        PacketType = 6
	F1CarStatus           PacketType = 7
	F1FinalClassification PacketType = 8
	F1LobbyInfo           PacketType = 9
	F1CarDamage           PacketType = 10
	F1SessionHistory      PacketType = 11
)
const (
	packetSnaplen       = 2 * 1024
	netInterfaceTimeout = 5 * time.Second
)

var EventCodes map[string]string
var buttonMap []string
var lapAnalyser LapAnalyser

func init() {
	initEventCodes()
	initButtons()
}

var other int

func main() {
	other = 0

	if len(os.Args) != 3 {
		log.Fatalln("neither pcap nor interface provided")
	}
	var pcapHandle *pcap.Handle
	var err error
	if os.Args[1] == "pcap" {
		filepath := os.Args[2]
		if pcapHandle, err = pcap.OpenOffline(filepath); err != nil {
			log.Fatalln("error opening or reading pcap: " + err.Error())
		}
	} else if os.Args[1] == "if" {
		netInterface := os.Args[2]
		if pcapHandle, err = pcap.OpenLive(netInterface, packetSnaplen, false, netInterfaceTimeout); err != nil {
			log.Fatalln("error opening interface " + netInterface + " : " + err.Error())
		}
	} else {
		log.Fatalln("invalid arguments")
	}

	packetSource := gopacket.NewPacketSource(pcapHandle, pcapHandle.LinkType())
	for packet := range packetSource.Packets() {
		handlePacket(packet)
	}

	saveResults(&lapAnalyser)
}

func handlePacket(packet gopacket.Packet) {
	appLayer := packet.ApplicationLayer()
	if appLayer == nil {
		return
	}

	var header PacketHeader
	reader := bytes.NewReader(appLayer.Payload())

	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		fmt.Println("Error reading header:", err.Error())
		return
	}

	if header.SessionTime < 0 || header.PacketFormat != 2022 {
		return
	}

	switch PacketType(header.PacketId) {
	case F1Motion:
		handleMotionData(reader, &header)
	case F1LapData:
		handleLapData(reader, &header)
	case F1Event:
		handlePacketEventData(reader)
	case F1CarTelemetry:
		handleCarTelemetryData(reader, &header)
	default:
		other++
	}
}

func handleMotionData(reader io.Reader, header *PacketHeader) {
	var motionData MotionData
	if err := binary.Read(reader, binary.LittleEndian, &motionData); err != nil {
		log.Println("Error reading motion data:", err.Error())
		return
	}
	lapAnalyser.analyzeMotionData(&motionData, uint(header.FrameIdentifier), uint(header.PlayerCarIndex))
}

func handleLapData(reader io.Reader, header *PacketHeader) {
	var packetLapData PacketLapData
	if err := binary.Read(reader, binary.LittleEndian, &packetLapData); err != nil {
		log.Printf("error reading LapData: %s\n", err.Error())
	}
	playerLapData := packetLapData.LapData[header.PlayerCarIndex]
	lapAnalyser.analyzeLapData(&playerLapData)
}

func handleCarTelemetryData(reader io.Reader, header *PacketHeader) {
	var packetCarTelemetryData PacketCarTelemetryData
	if err := binary.Read(reader, binary.LittleEndian, &packetCarTelemetryData); err != nil {
		log.Printf("error reading LapData: %s\n", err.Error())
		return
	}
	lapAnalyser.analyzeCarTelemetryData(&packetCarTelemetryData, uint(header.FrameIdentifier), uint(header.PlayerCarIndex))
}

func handlePacketEventData(reader io.Reader) {
	var eventCode [4]uint8
	if err := binary.Read(reader, binary.LittleEndian, &eventCode); err != nil {
		log.Println("error reading event data")
		return
	}
	printEventDataDetails(string(eventCode[:]), reader)
}

func printEventDataDetails(eventCode string, r io.Reader) {
	var err error
	switch eventCode {
	case "FTLP":
		var details FastestLap
		err = binary.Read(r, binary.LittleEndian, &details)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		lap := strconv.FormatFloat(float64(details.LapTime), 'g', -1, 32)
		fmt.Println("Fastest lap: " + strconv.Itoa(int(details.VehicleIdx)) + " " + lap)
	case "RTMT":
		var details Retirement
		fmt.Println("Retirement: " + strconv.Itoa(int(details.VehicleIdx)))
	case "TMPT":
		var details TeamMateInPits
		fmt.Println("TeamMateInPits: " + strconv.Itoa(int(details.VehicleIdx)))
	case "RCWN":
		var details RaceWinner
		fmt.Println("RaceWinner: " + strconv.Itoa(int(details.VehicleIdx)))
	case "PENA":
		var details Penalty
		fmt.Println("Penalty: " + strconv.Itoa(int(details.VehicleIdx)))
	case "SPTP":
		var details SpeedTrap
		fmt.Printf("SpeedTrap: %d, speed: %f \n", details.VehicleIdx, details.Speed)
	case "STLG":
		var details StartLights
		fmt.Printf("StartLights: %d \n", details.NumLights)
	case "BUTN":
		var details Buttons
		err = binary.Read(r, binary.LittleEndian, &details)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var buttons []string
		for i, v := range buttonMap {
			if details.ButtonStatus&(1<<i) != 0 {
				buttons = append(buttons, v)
			}
		}
		//fmt.Printf("Buttons being used: %v\n", buttons)
	case "SSTA":
		fmt.Println("Session started")
		lapAnalyser.sessionStarted()
	case "SEND":
		fmt.Println("Session ended")
	default:
		log.Printf("Unknown event code: %s\n", eventCode)
	}
}

func saveResults(lapAnalyser *LapAnalyser) {
	for i, lap := range lapAnalyser.pastLaps {
		carLine := getCarLine(&lap)
		saveCarLine(&carLine, fmt.Sprintf("lap%d.png", i))
	}
}
