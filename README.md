# StepKeys

**StepKeys is a fully-featured, cross-platform pedal-based input system that maps physical footswitches to keyboard actions, sequences, and combos.**

## Motivation

Human interaction with computers is almost always limited to hands - keyboards, mice, MIDI controllers, touchscreens, etc. All follow the same paradigm. But interaction design shouldn’t follow hardware traditions - it should follow human capabilities.

Pedals (footswitches) introduce a third interaction channel. An independent and parallel one. This enables actions that would otherwise compete with typing, aiming, or mouse movement. It creates new interaction patterns through a new form of input, greatly enhancing accessibility and productivity.

**This unlocks entirely new ways to interact with a computer:**

- Trigger complex key combinations with a single press

- Type letters, words, or full strings instantly with foot

- Hold keys till later pedal input or release

- Bind game mechanics, abilities, or macros to pedals

- Create parallel workflows without interrupting primary input

**StepKeys mean real multitasking. StepKeys mean interface evolution.**

## Key Features

- Multi-pedal support (up to 128 pedals)

- Minimal hardware requirements and largely hardware-agnostic architecture

- Fast and efficient Go backend server

- Pedal behavior modes: hold, toggle, and oneshot

- Trigger single keys, key combinations, or sequences of key presses

- Modern graphical interface for pedal assignment configuration

- System tray menu for quick access

- Detailed logging of errors and operations

- Public API and WebSockets with detailed documentation for third-party integrations

- Cross-platform: Windows, Linux, macOS

- Lightweight and portable binary builds

- Easy-to-set-up project repo

## Hardware

StepKeys is designed in a way that the (external) hardware can be easily replaced or substituted with products from other vendors. The software monitors a serial port and listens for incoming bytes that identify a footswitch and a press or release event.

**The original (reference) implementation:**

- **External MCU:** Arduino Leonardo

> [!TIP]
> Other Arduino models, a Raspberry Pi, or ESP boards can also be used. Basically, any device capable of serial communication can be used as the MCU.

- **Foot switches:** 3-pin momentary switches

> [!TIP]
> This is a practical choice because no pull-down resistors are required. Wire as follows: COM - Aruino pin, NC (normally closed) - GND, NO (normally open) - +5V.
>
> StepKeys expects momentary switches. For other types, the MCU code must simulate momentary presses. Even the cheapest ones from Aliexpress will do the job.

- **Wiring between the MCU and footswitches:** CAT cable (preferably CAT6 for better shielding). +5V and GND should be distributed from a common split point (star topology).

- The **MCU code** must:

  - Follow the protocol described [here](https://github.com/BrNi05/StepKeys/blob/main/arduino/code/stepkeys.ino#L44)
  - Send data as single, raw bytes
  - Include some kind of debounce mechanism (even if the footswitches supposedly have it built-in)
  - Specify a baud rate (default: 115200), which should match the value in the `.env` file.

> [!WARNING]
> Due to the protocol used, the maximum number of pedals supported is 128 (0-127).

> [!WARNING]
> Wire the pedals to the MCU using consecutive Arduino pins (the pins must follow each other). You can use StepKeys with phantom pedals in your configuration, but it’s cleaner to wire them in order.

## How to set up StepKeys

Use the installer scripts, which will guide you through the entire installation and setup process.

### Linux / macOS

> [!IMPORTANT]
> On Linux, depending on your distro, you may need to grant your user permission to access serial devices. Restart your device for changes to take effect.
>
> Use: `sudo usermod -aG dialout <username>` or `sudo usermod -aG uucp <username>`.

``` bash
bash <(curl -sL https://raw.githubusercontent.com/BrNi05/StepKeys/main/release/posix.sh)
```

### Windows

``` bash
irm -useb https://raw.githubusercontent.com/BrNi05/StepKeys/main/release/windows.ps1 | iex
```

> [!TIP]
> During installation, you will be prompted to select the serial device to be used by StepKeys. The installer will list all detected serial devices.
>
> If you are unsure which one to choose:
>  
> On Windows, open **Device Manager** and look for entries like **USB Serial Device (COMX)** under **Ports (COM & LPT)**.
>
> On POSIX systems, you can identify devices using platform-specific tools. On Linux, use **udevadm** to inspect `/dev/ttyACM*` or `/dev/ttyUSB*` devices. On macOS, use **ioreg** to inspect `/dev/cu.usbmodem*` or `/dev/cu.usbserial*`.
>
> Alternatively, the **Arduino Cloud Agent** together with the **Arduino Cloud Editor** can be used to identify connected Arduino devices.

## How to update StepKeys

StepKeys includes a built-in version manager and will notify you whenever an update is available. You will be redirected to the release page, where you will find links to this section of the README.

### Linux / macOS

``` bash
bash <(curl -sL https://raw.githubusercontent.com/BrNi05/StepKeys/main/release/posix.sh) update
```

### Windows

``` bash
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/BrNi05/StepKeys/main/release/windows.ps1))) update
```

## How to use StepKeys

### Tray menu

StepKeys comes with a tray menu, which is the main interface you’ll use most of the time - aside from the pedals, of course.

- **Open:** opens StepKeys webGUI.

> [!TIP]
> If an update is available, this menu item will show as: **Open (Update Available)**.

- **Enabled:** ticked if StepKeys is enabled ankvd listening on the serial port.

> [!IMPORTANT]
> StepKeys cannot be enabled if there is no pedal configuration created yet.
>
> If StepKeys cannot open the serial port, the enabled status won't have an effect on functionality as serial listening will be disabled until a restart.

- **Start on boot:** toggles whether StepKeys should start on boot or not.

- **Docs:** opens the GitHub page of Stepkeys and shows this README.

- **API Docs:** opens the StepKeys API documentation (Swagger).

- **WebSocket Docs:** opens the [documentation](https://github.com/BrNi05/StepKeys/blob/main/WebSocketDocs.md) of WebSockets used by StepKeys.

- **Quit:** terminates StepKeys process.

### webGUI

StepKeys includes a built-in GUI that opens in your preferred browser.

<br />

<img width="1708" height="909" alt="webGUI snippet"
     src="https://github.com/user-attachments/assets/50c2ae58-6df0-4398-830e-0243c451dfa2" />

<br />

> [!TIP]
> One can create a custom GUI, since the built-in is powered entirely by the public API and WebSockets. It wouldn't take much effort to create a native appearance for StepKeys ([WebView](https://github.com/webview/webview_go)), but I consider opening in browser more reliable.

> [!IMPORTANT]
> By default, StepKeys server uses port **18000**. In case of a conflit, you can modfiy your existing port mappings or assign an other port for StepKeys in the **config.json** file.

You can open the GUI from the tray or [here](http://localhost:18000/). You will notice it follows a fairly standard and clean approach. There is a top and a bottom (lower) menu bar, and two side-by-side windows.

### Top menu bar

- **Enabled toggle:** shows the enabled status. It will not toggle on, if there is no (internal) pedal map set.

- **Start on boot toggle:** shows and toggles start on boot state.

- **GitHub icon:** opens the GitHub repo of StepKeys.

### Bottom menu bar

- **Serial display:** shows the serial port (Windows) or device file (macOS/Linux) that StepKeys is using or attempted to use.

- **Check for updates / Update available:** depending on availability, one button will appear. **Check for updates** forces StepKeys to look for updates again (it automatically checks on app startup).

### Log Viewer

This is an extremely useful feature of the webGUI. It displays all server logs that were generated after app startup. These logs are persistent, as they are logged to file as well.

> [!IMPORTANT]
> The webGUI may sometimes log to browser console. These are not displayed here as such logs are supposed to be rare and technical.

> [!TIP]
> If you are unsure about the physical order of your pedals, just use the **Log Viewer**, as it will display all pedal actions in realtime. You will need to enable StepKeys for this.

### Pedal Configurator

This is probably the most important feature of the GUI and has quite a few UI elements.

The configurator window has a dedicated lower menu bar:

- **Reset:** discard all changes made on the GUI and reload the saved pedal config.

- **Input field (Profile name):** assign a name for the current config.

- **Save Config:** save the current webGUI config (with the assigned name).

- **Load config:** browse and load a saved StepKeys pedal map config.

- **Apply:** send the current changes to backend and apply them.

Once you have at least one pedal in the configurator, you will see the pedal cards. These represent a pedal and an assigned action. The **Add Pedal** button always adds the smallest indexed pedal possible, this is the reason for the consecutive wiring.

On a pedal card you can see its ID, a button to **remove** it from the configurator, a **mode** and **behaviour** dropdown and the keys assignment input field. Added keys appear to the left to it, with an **X** button to remove them.

There are quite a few assistive mechanisms in place. If you start typing a key, suggestions will appear. You can navigate with the **arrow keys** and accept the selected with **Space** or **Enter**. When no characters are typed, use **Backspace** to remove the last added key.

> [!TIP]
> You may notice that key repetition is not allowed in **combo** mode. While StepKeys can handle configurations with repeated keys, the GUI experience is more streamlined with this restriction in place. In **sequence** mode, repeated keys are allowed.

Have a look at the [list of supported keys](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys) and the StepKeys [implementation](https://github.com/BrNi05/StepKeys/blob/main/server/pedal/supported_keys.go).

> [!IMPORTANT]
> StepKeys server tracks and knows about one config (profile). It does not natively include profile management. However, the webGUI has such feature. When saving a profile, the state of the webGUI is saved, which might not match the loaded profile (internal state).

## StepKeys API

The shipped binaries contain both the API code and the API documentation that is available from the **tray menu** or [here](http://localhost:18000/api/docs/index.html).

The API provides an alternative way to interact with StepKeys, but it is not intended to satisfy every technical expectation one might have of an API. StepKeys is primarily controlled through the system tray menu and the included GUI, which handle certain tasks that would otherwise fall under the API’s responsibilities.

For example, toggle operations handle errors silently, meaning a **HTTP200** response code is sent even if the toggling failed. Since such errors are expected to be rare, and the API response includes the boolean value that was supposed to be toggled (so success can be infered from historic and current value), this approach is not a limitation, just a different way of working with an API.

StepKeys enforces limitations of **RobotGo**, that is an amazing third party lib which interacts with the OS to press keys. The GUI will provide visual aid for supported keys, while API users can refer to the [official docs](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys) or the StepKeys [support list](https://github.com/BrNi05/StepKeys/blob/main/server/pedal/supported_keys.go).

> [!TIP]
> For more complex operations, you may need to use the WebSockets provided by the StepKeys server. See the [documentation](https://github.com/BrNi05/StepKeys/blob/main/WebSocketDocs.md) for details.

## Technical notes

StepKeys logs with moderate verbosity and handles most runtime errors. It assumes the user is familiar with the app and does NOT guard against deliberate misuse.

StepKeys behaves differently from other programs regarding its lifecycle. It always exits with code 0, even if execution is halted due to a handled error. During app startup, the shutdown process can be fragile for a few milliseconds, since it depends on the system tray, which may not yet be initialized. However, given StepKeys’ small scope and the testing that has been performed, such issues are expected to be extremely rare.

## Set up project

1. Clone the repo:

``` bash
git clone https://github.com/BrNi05/StepKeys.git
```

2. Make sure you have **Go** installed and **GOPATH** set.

> [!TIP]
> Check the current **GOPATH** with: `go env GOPATH`.
>
> Example: add it to **PATH**: `echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc`

3. Install dependencies

``` bash
go install github.com/swaggo/swag/cmd/swag@latest # this is a global install

go get go.bug.st/serial
go get github.com/go-vgo/robotgo
go get github.com/getlantern/systray
go get github.com/pkg/browser
go get github.com/joho/godotenv
go get github.com/swaggo/swag
go get github.com/swaggo/http-swagger
go get github.com/gorilla/websocket

go mod tidy

go env -w CGO_ENABLED=1

cd gui
npm install # the GUI needs node and npm to be installed
```

**Robotgo** needs `gcc` to work. On Windows, use MSYS2. On POSIX, use your preferred package manager to install it.

> [!IMPORTANT]
> Restart VSC or your terminal session for changes to take effect.

4. Create the `.env` in **server** directory

5. Start working with StepKeys

- Use the VS Code tasks and launch config to start or build StepKeys.

> [!TIP]
> On Linux, **Systray** needs a few Linux libs to compile. You will encounter build errors if these are missing. Example solutions:
>
> Fedora: `sudo dnf install -y pkg-config libayatana-appindicator-gtk3-devel`
>
> Ubuntu: `sudo apt-get install -y pkg-config libayatana-appindicator3-dev`
>
> Some deprecation warnings may still appear during build or runtime, but these are expected and can safely be ignored.

> [!IMPORTANT]
> On Linux, depending on your distro, you may need to grant your user permission to access serial devices. Restart your device for changes to take effect.
>
> Use: `sudo usermod -aG dialout <username>` or `sudo usermod -aG uucp <username>`.

## Attributions

StepKeys icon is from Vecteezy ([Shapes Vectors by Vecteezy](https://www.vecteezy.com/free-vector/shapes)). The original (white) background was removed by me.
