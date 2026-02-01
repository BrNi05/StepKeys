package config

import (
	"os"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"

	Log "stepkeys/server/logging"
)

var _serialPort string
var baudRate int
var appVersion string

func LoadEnv() (baudRate int, serialPort string) {
	// Load .env file
	// Ignore error if missing
	_ = godotenv.Load(".env")
	Log.WriteToLogFile("Environment file loaded.")

	// SERIAL_PORT
	serialPort = os.Getenv("SERIAL_PORT")
	if serialPort == "" {
		serialPort = defaultSerialPort()
		Log.WriteToLogFile("Using default for SERIAL_PORT env var: " + serialPort)
	}
	_serialPort = serialPort

	// BAUD_RATE
	if b := os.Getenv("BAUD_RATE"); b != "" {
		if val, err := strconv.Atoi(b); err == nil {
			baudRate = val
		} else {
			baudRate = 115200
			Log.WriteToLogFile("Invalid BAUD_RATE env var, using default: " + strconv.Itoa(baudRate))
		}
	} else {
		baudRate = 115200
		Log.WriteToLogFile("Using default for BAUD_RATE env var: " + strconv.Itoa(baudRate))
	}

	// VERSION
	appVersion = os.Getenv("VERSION")
	if appVersion == "" {
		appVersion = "1.0.0"
		Log.WriteToLogFile("Using default for VERSION env var: " + appVersion)
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
		Log.WriteToLogFile("Unsupported OS detected.")
		return "unknown"
	}
}

func GetAppVersion() string {
	return appVersion
}

// The API uses this as an additional way to get the serial port (for the GUI)
func GetSerialPort() string {
	return _serialPort
}
