package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func analyzeBytes(rawData string, version int8) bool {
	fmt.Println("_____________________________________________________________________________________________________________________")
	//fmt.Println("rawData ", len(rawData))

	if len(rawData) == 24 {
		fmt.Println("Sensit v", version, "Downlink Frame")
		return false
	}

	fmt.Println("Sensit v", version, "Uplink Frame")

	parsed, err := strconv.ParseUint(rawData, 16, 32)
	if err != nil {
		log.Fatal(err)
	}
	data := fmt.Sprintf("%08b", parsed)
	byte1 := data[0:8]
	byte2 := data[8:16]
	byte3 := data[16:24]
	byte4 := data[24:32]

	if version == 2 {
		if len(data) == 25 { //Low battery MSB
			fmt.Println("Sensit Low battery")
			//TODO: Handle low battery bit shift
			return false
		}

		//Byte 1
		batteryMsb := data[0:1]
		eventType, _ := strconv.ParseInt(data[1:3], 2, 8)
		timeframe, _ := strconv.ParseInt(data[3:5], 2, 8)
		mode, _ := strconv.ParseInt(data[5:8], 2, 8)

		//Byte 2
		temperatureMsb := data[8:12]
		batteryLsb := data[12:16]
		battData := []string{batteryMsb, batteryLsb}
		battery, _ := strconv.ParseInt(strings.Join(battData, ""), 2, 8)
		batVal := (float64(battery) * 0.05) + 2.7

		//Byte 3
		temperature := int64(0)
		tempVal := float64(0)

		reedSwitch := false
		if mode == 0 || mode == 1 {
			temperatureLsb := data[18:24]
			tempData := []string{temperatureMsb, temperatureLsb}
			temperature, _ := strconv.ParseInt(strings.Join(tempData, ""), 2, 16)
			tempVal = (float64(temperature) - 200) / 8
			if data[17] == 1 {
				reedSwitch = true
			}
		} else {
			temperature, _ = strconv.ParseInt(temperatureMsb, 2, 16)
			tempVal = (float64(temperature) - 200) / 8
		}

		modeStr := ""
		swRev := ""
		humidity := float64(0.0)
		light := float64(0.0)

		switch mode {
		case 0:
			modeStr = "Button"
			majorSwRev, _ := strconv.ParseInt(data[24:28], 2, 8)
			minorSwRev, _ := strconv.ParseInt(data[28:32], 2, 8)
			swRev = fmt.Sprintf("%d.%d", majorSwRev, minorSwRev)
		case 1:
			modeStr = "Temperature + Humidity"
			humi, _ := strconv.ParseInt(data[24:32], 2, 16)
			humidity = float64(humi) * 0.5
		case 2:
			modeStr = "Light"
			lightVal, _ := strconv.ParseInt(data[18:24], 2, 8)
			lightMulti, _ := strconv.ParseInt(data[17:18], 2, 8)
			light = float64(lightVal) * 0.01
			if lightMulti == 1 {
				light = light * 8
			}
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

		// Outputing and printing values
		fmt.Println("---------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Raw data :", byte1, byte2, byte3, byte4)
		fmt.Println("Mode", mode, ":", modeStr, "\t\t", "Event type", eventType, ":", typeStr, "\t\t", "Timeframe", timeframe, ":", timeStr)
		fmt.Println("Battery :", batVal, "V\t\t")
		switch mode {
		case 0:
			fmt.Println("Temperature :", tempVal, "°C")
			fmt.Println("v" + swRev)
		case 1:
			fmt.Println("Temperature :", tempVal, "°C")
			fmt.Println(humidity, "% RH")
		case 2:
			fmt.Println(light, "lux")
			alerts, _ := strconv.ParseInt(data[24:32], 2, 16)
			fmt.Println("Number of alerts :", alerts)
		case 3, 4, 5:
			alerts, _ := strconv.ParseInt(data[24:32], 2, 16)
			fmt.Println("Number of alerts :", alerts)
		}
		if reedSwitch {
			fmt.Println("Reed switch on")
		}
		fmt.Println("_____________________________________________________________________________________________________________________")
	} else { //version == 3
		modeStr := ""
		//Byte 1 : 5b Battery & 3b reserved (0b110)
		battery, _ := strconv.ParseInt(data[0:5], 2, 8)
		// reserved, _ := strconv.ParseInt(data[5:8], 2, 8) //Should be 0b110

		//Byte 2 : 5b Mode, 1b Alert Button, 2b data
		buttonFlag := data[13:14]
		mode, _ := strconv.ParseInt(data[8:13], 2, 8)
		switch mode {
		case 0:
			modeStr = "Standby"
			/*majorSwRev, _ := strconv.ParseInt(data[24:28], 2, 8)
			minorSwRev, _ := strconv.ParseInt(data[28:32], 2, 8)
			swRev = fmt.Sprintf("%d.%d", majorSwRev, minorSwRev)*/
		case 1:
			modeStr = "Temperature + Humidity"
			/*humi, _ := strconv.ParseInt(data[24:32], 2, 16)
			humidity = float64(humi) * 0.5*/
		case 2:
			modeStr = "Light"
			/*lightVal, _ := strconv.ParseInt(data[18:24], 2, 8)
			lightMulti, _ := strconv.ParseInt(data[17:18], 2, 8)
			light = float64(lightVal) * 0.01
			if lightMulti == 1 {
				light = light * 8
			}*/
		case 3:
			modeStr = "Door"
		case 4:
			modeStr = "Vibration"
		case 5:
			modeStr = "Magnet"
		default:
			modeStr = ""
		}
		/*temperatureMsb := data[8:12]
		batteryLsb := data[12:16]
		battData := []string{batteryMsb, batteryLsb}
		battery, _ := strconv.ParseInt(strings.Join(battData, ""), 2, 8)
		batVal := (float64(battery) * 0.05) + 2.7*/

		//Byte 3
		/*temperature := int64(0)
		tempVal := float64(0)*/
		// Outputing and printing values
		fmt.Println("---------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Raw data :", byte1, byte2, byte3, byte4)
		fmt.Println("Mode", mode, ":", modeStr, "\t\t", "Event type") /*, eventType, ":", typeStr, "\t\t", "Timeframe", timeframe, ":", timeStr)*/
		fmt.Println("Battery :", battery, "V\t\t Button: ", buttonFlag)
		/*switch mode {
		case 0:
			fmt.Println("Temperature :", tempVal, "°C")
			fmt.Println("v" + swRev)
		case 1:
			fmt.Println("Temperature :", tempVal, "°C")
			fmt.Println(humidity, "% RH")
		case 2:
			fmt.Println(light, "lux")
			alerts, _ := strconv.ParseInt(data[24:32], 2, 16)
			fmt.Println("Number of alerts :", alerts)
		case 3, 4, 5:
			alerts, _ := strconv.ParseInt(data[24:32], 2, 16)
			fmt.Println("Number of alerts :", alerts)
		}
		if reedSwitch {
			fmt.Println("Reed switch on")
		}*/
		fmt.Println("_____________________________________________________________________________________________________________________")
	}

	/*} else { //len: 24 exactly, 12 bytes
		fmt.Println("Sensit Daily Downlink Message")
		//TODO: Decode sensit downlink message
	}*/

	return true
}

func main() {
	var version int8
	if os.Args[1] == "v3" {
		version = 3
	} else {
		version = 2
	}

	for i := 2; i <= len(os.Args[1:]); i++ {
		analyzeBytes(os.Args[i], version)
	}
}
