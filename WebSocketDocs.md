
# WebSocket Endpoints

## Live logs

**URL:** `/ws/logs`

**Method:** GET (upgrade to WebSocket)

**Description:** Connect to receive live log lines from StepKeys.

**Message format:** text (string). **Eg.:** _2026/01/29 14:47:10 Loading app and pedal config._

Example usage can be found [here](https://github.com/BrNi05/StepKeys/blob/main/gui/src/components/LogViewer.vue).

## Settings states

**URL:** `/ws/settings`

**Method:** GET (upgrade to WebSocket)

**Description:** Connect to receive real-time updates whenever the StepKeys enabled or start on boot state changes.

**Message format:** JSON object **Eg.:** _{ "event": "enabled", "value": true }_

- Accepted event types: **enabled** and **boot**.
- Accepted values: booleans.

Example usage can be found [here](https://github.com/BrNi05/StepKeys/blob/main/gui/src/components/TopBar.vue).
