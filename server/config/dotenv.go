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

func LoadEnv() (baudRate int, serialPort string) {
	// Load .env file
	// Ignore error if missing
	_ = godotenv.Load(".env")
	log.Println("Environment file loaded.")

	// SERIAL_PORT
	serialPort = os.Getenv("SERIAL_PORT")
	if serialPort == "" {
		serialPort = defaultSerialPort()
		log.Println("Using default for SERIAL_PORT env var:", serialPort)
	}

	// BAUD_RATE
	if b := os.Getenv("BAUD_RATE"); b != "" {
		if val, err := strconv.Atoi(b); err == nil {
			baudRate = val
		} else {
			baudRate = 115200
			log.Println("Invalid BAUD_RATE env var, using default:", baudRate)
		}
	} else {
		baudRate = 115200
		log.Println("Using default for BAUD_RATE env var:", baudRate)
	}

	return baudRate, serialPort
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
