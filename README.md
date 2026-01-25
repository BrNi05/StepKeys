# StepKeys

StepKeys is a cross-platform pedal-based input system that maps physical footswitches to keyboard actions, sequences, and combos.

> [!WARNING]
> StepKeys is currently in beta and may have some bugs. Please take this into account when using it.

## Motivation

Human interaction with computers is almost always limited to hands - keyboards, mice, MIDI controllers, touchscreens, all follow the same paradigm.
But interaction design shouldn’t follow hardware traditions - it should follow human capabilities.

Pedals (footswitches) introduce a third interaction channel. An independent and parallel one. This enables actions that would otherwise compete with typing, aiming, or mouse movement. It creates new interaction patterns through a new form of input, greatly enhancing accessibility and productivity.

**This unlocks entirely new ways to interact with a computer:**

- Trigger complex key combinations with a single press

- Type letters, words, or full strings instantly

- Hold keys till later pedal input

- Bind game mechanics, abilities, or macros to pedals

- Create parallel workflows without interrupting primary input

**StepKeys mean real multitasking. StepKeys mean interface evolution.**

## Key Features

- Multi-pedal support (up to 128 pedals)

- Minimal hardware requirements and largely hardware-agnostic architecture

- Fast and efficient Go backend server

- Pedal behavior modes: hold, toggle, and oneshot

- Trigger single keys, key combinations, or sequences of key presses

- Graphical interface for pedal assignment configuration

- System tray menu for quick access

- Detailed logging of errors and operations

- Public API with documentation for third-party integrations

- Cross-platform: Windows, Linux, macOS

- Lightweight and portable binary builds

## Hardware

StepKeys is designed in a way that the (external) hardware can be easily replaced or substituted with products from other vendors. The software monitors a serial port and listens for incoming bytes that identify a footswitch and press/release event.

**The original implementation:**

- **External MCU:** Arduino Leonardo

> [!TIP]
> Other Arduino models, a Raspberry Pi, or ESP boards can also be used. Basically, any device capable of serial communication can be used as the MCU.

- **Foot switches:** 3-pin momentary switches

> [!TIP]
> This is a practical choice because no pull-down resistors are required. Wire as follows: COM - Aruino pin, NC (normally closed) - GND, NO (normally open) - +5V.
>
> StepKeys expects momentary switches. For other types, the MCU code must simulate momentary presses. Even the cheapest ones from Aliexpress will do the job.

- **Wiring between the MCU and footswitches:** CAT cable (preferably CAT6 for better shielding)

- The **MCU code** must:

  - Follow the protocol described [here](https://github.com/BrNi05/StepKeys/blob/main/arduino/code/stepkeys.ino#L44)
  - Send data as single, raw bytes
  - Include some kind of debounce mechanism
  - Specify a baud rate (default: 115200), which should match the value in the `.env` file.

> [!WARNING]
> Due to the protocol used, the maximum number of pedals supported is 128 (0-127).

## How to set up StepKeys server

Use the installer scripts, which will guide you through the entire installation and setup process.

### Linux / macOS

``` bash
curl -fsSL https://raw.githubusercontent.com/BrNi05/StepKeys/main/release/posix.sh | bash
```

### Windows

``` bash
iwr -useb https://raw.githubusercontent.com/BrNi05/StepKeys/main/release/windows.ps1 | iex
```

## How to use StepKeys

### Tray menu

### webGUI

## StepKeys API

The shipped binaries contain both the API code and the API documentation that is available from the **tray menu** or [here](http://localhost:18000/api/docs/index.html).

The API provides an alternative way to interact with StepKeys, but it is not intended to satisfy every technical expectation one might have of an API. StepKeys is primarily controlled through the system tray menu and the included GUI, which handle certain tasks that would otherwise fall under the API’s responsibilities.

For example, toggle operations handle errors silently, meaning a HTTP200 response code is sent even if the toggling failed. Since such errors are expected to be rare and the API response includes the boolean value that was supposed to be toggled (so success can be infered from historic and current value), this approach is not a limitation, just a different way of working with an API.

StepKeys enforces limitations of **RobotGo**, that is an amazing third party lib which interacts with the OS to press keys. The GUI will provide visual aid for supported keys, while API users can refer to the [official docs](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys) or the StepKeys [support list](https://github.com/BrNi05/StepKeys/blob/main/server/pedal/supported_keys.go).

## Set up project

1. Clone the repo:

``` bash
git clone https://github.com/BrNi05/StepKeys.git
```

2. Make sure you have **Go** installed and **GOPATH** set.

> [!TIP]
> Check the current **GOPATH** with: `go env GOPATH`.

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

go mod tidy

cd gui
npm install # the GUI needs a node and npm to be installed
```

> [!IMPORTANT]
> Restart VSC or your terminal session for changes to take effect.

4. Create the `.env` in **server** directory

5. Start working with StepKeys

- Use the VS Code tasks and launch config to start or build StepKeys.

## Attributions

StepKeys icon is from Vecteezy ([Shapes Vectors by Vecteezy](https://www.vecteezy.com/free-vector/shapes)).
