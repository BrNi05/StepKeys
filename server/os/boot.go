package os

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

const name = "StepKeys"

var pathToStartipTrigger string

// Initialize path to autostart trigger based on OS
func init() {
	switch runtime.GOOS {
	case "windows":
		pathToStartipTrigger = filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup", name+".lnk")
	case "darwin":
		usr, _ := user.Current()
		pathToStartipTrigger = filepath.Join(usr.HomeDir, "Library", "LaunchAgents", name+".plist")
		os.MkdirAll(filepath.Dir(pathToStartipTrigger), 0755)
	case "linux":
		usr, _ := user.Current()
		pathToStartipTrigger = filepath.Join(usr.HomeDir, ".config", "autostart", name+".desktop")
		os.MkdirAll(filepath.Dir(pathToStartipTrigger), 0755)
	default:
		pathToStartipTrigger = "" // unsupported OS
	}
}

func EnableStartOnBoot() error {
	if isEnabledOnBoot() {
		return nil // already enabled
	}

	switch runtime.GOOS {
	case "windows":
		return enableWindows()
	case "darwin":
		return enableMac()
	case "linux":
		return enableLinux()
	default:
		return errors.New("Autostart not supported on this OS")
	}
}

// Checks if autostart is enabled on OS level
func isEnabledOnBoot() bool {
	if _, err := os.Stat(pathToStartipTrigger); err == nil {
		return true
	}
	return false
}

func DisableStartOnBoot() error {
	if runtime.GOOS == "darwin" {
		exec.Command("launchctl", "unload", pathToStartipTrigger).Run()
	}

	if _, err := os.Stat(pathToStartipTrigger); err == nil {
		return os.Remove(pathToStartipTrigger)
	}
	return nil
}

// Windows: create shortcut in Startup folder
func enableWindows() error {
	cmd := fmt.Sprintf(`$s=(New-Object -COM WScript.Shell).CreateShortcut("%s");$s.TargetPath="%s";$s.Save()`, pathToStartipTrigger, GetExecPath())
	return exec.Command("powershell", "-Command", cmd).Run()
}

// macOS: create LaunchAgent plist
func enableMac() error {
	content := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key><string>%s</string>
    <key>ProgramArguments</key><array><string>%s</string></array>
    <key>RunAtLoad</key><true/>
</dict>
</plist>`, name, GetExecPath())
	return os.WriteFile(pathToStartipTrigger, []byte(content), 0644)
}

// Linux: create .desktop file in ~/.config/autostart
func enableLinux() error {
	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Exec=%s
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Name=%s
Comment=Start %s on login
`, GetExecPath(), name, name)
	return os.WriteFile(pathToStartipTrigger, []byte(content), 0644)
}
