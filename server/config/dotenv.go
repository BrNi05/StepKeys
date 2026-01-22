package config

import (
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

var serialPort string
var baudRate int

func LoadEnv() {
	// Load .env file
	// Ignore error if missing
	_ = godotenv.Load(".env")

	// SERIAL_PORT
	serialPort = os.Getenv("SERIAL_PORT")
	if serialPort == "" {
		serialPort = defaultSerialPort()
	}

	// BAUD_RATE
	if b := os.Getenv("BAUD_RATE"); b != "" {
		if val, err := strconv.Atoi(b); err == nil {
			baudRate = val
		} else {
			baudRate = 115200
		}
	} else {
		baudRate = 115200
	}
}

// Fallback default per OS
func defaultSerialPort() string {
	switch runtime.GOOS {
	case "windows":
		return "COM3"
	case "linux":
		return "/dev/ttyACM0"
	case "darwin":
		return "/dev/cu.usbmodem11301"
	default:
		log.Fatal("Unsupported OS")
		return ""
	}
}

// Returns the baud rate
func GetBaudRate() int {
	return baudRate
}

// Returns the serial port
func GetSerialPort() string {
	return serialPort
}
