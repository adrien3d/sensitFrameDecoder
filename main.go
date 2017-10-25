package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func analyzeBytes(data string) {
	// TODO : Fix shift when battery MSB=0
	// TODO : Handle modes 2, 3, 4 & 5
	byte1 := data[0:8]
	byte2 := data[8:16]
	byte3 := data[16:24]
	byte4 := data[24:32]

	//Byte 1
	mode, _ := strconv.ParseInt(data[5:8], 2, 16)
	timeframe, _ := strconv.ParseInt(data[3:5], 2, 16)
	eventType, _ := strconv.ParseInt(data[1:3], 2, 16)
	batteryMsb := data[0:1]

	//Byte 2
	temperatureMsb := data[8:12]
	batteryLsb := data[12:16]
	battData := []string{batteryMsb, batteryLsb}
	battery, _ := strconv.ParseInt(strings.Join(battData, ""), 2, 16)
	batVal := float32(battery) * 0.05 * 2.7

	//Byte 3
	temperature := int64(0)
	tempVal := float32(0)

	reedSwitch := false
	if mode == 0 || mode == 1 {
		temperatureLsb := data[18:24]
		tempData := []string{temperatureMsb, temperatureLsb}
		temperature, _ := strconv.ParseInt(strings.Join(tempData, ""), 2, 16)
		tempVal = (float32(temperature) - 200) / 8
		if data[17] == 1 {
			reedSwitch = true
		}
	} else {
		temperature, _ = strconv.ParseInt(temperatureMsb, 2, 16)
		tempVal = (float32(temperature) - 200) / 8
	}

	modeStr := ""
	switch mode {
	case 0:
		modeStr = "Button"
	case 1:
		modeStr = "Temperature + Humidity"
	case 2:
		modeStr = "Light"
	case 3:
		modeStr = "Door"
	case 4:
		modeStr = "Move"
	case 5:
		modeStr = "Reed switch"
	default:
		modeStr = ""
	}

	timeStr := ""
	switch timeframe {
	case 0:
		timeStr = "10 mins"
	case 1:
		timeStr = "1 hour"
	case 2:
		timeStr = "6 hours"
	case 3:
		timeStr = "24 hours"
	default:
		timeStr = ""
	}

	typeStr := ""
	switch eventType {
	case 0:
		typeStr = "Regular, no alert"
	case 1:
		typeStr = "Button call"
	case 2:
		typeStr = "Alert"
	case 3:
		typeStr = "New mode"
	default:
		timeStr = ""
	}

	//fmt.Println(data)
	fmt.Println("Raw data :", byte1, byte2, byte3, byte4)
	fmt.Println("Mode", mode, ":", modeStr, "\t\t", "Event type", eventType, ":", typeStr, "\t\t", "Timeframe", timeframe, ":", timeStr)
	fmt.Println("Battery :", batVal, "V\t\t", "Temperature :", tempVal, "°C")
	if reedSwitch {
		fmt.Println("Reed switch on")
	}
}

func formatData(data string) {
	//decoded, err := hex.DecodeString(data)
	parsed, err := strconv.ParseUint(data, 16, 32)
	if err != nil {
		log.Fatal(err)
	}
	pars := fmt.Sprintf("%08b", parsed)
	analyzeBytes(pars)
}

func main() {
	frameBits := "e86e1a21"
	formatData(frameBits)
}
