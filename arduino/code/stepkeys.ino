// The first digital pin a pedal is connected to
constexpr int FIRST_PIN = 2;

// The number of pedals connected
// So the last pin used will be FIRST_PIN + NUM_PEDALS - 1
constexpr int NUM_PEDALS = 5;

uint8_t pedalState[NUM_PEDALS];
uint8_t prevPedalState[NUM_PEDALS];

void setup() {
  // Baud rate: 115200
  Serial.begin(115200);

  for (uint8_t i = 0; i < NUM_PEDALS; i++) {
    pinMode(FIRST_PIN + i, INPUT);
    pedalState[i] = LOW;
    prevPedalState[i] = HIGH; // logically correct this way
  }
}

void loop() {
  for (uint8_t i = 0; i < NUM_PEDALS; i++) {
    pedalState[i] = digitalRead(FIRST_PIN + i);

    // Transition: LOW -> HIGH (PRESS)
    // MSB: 1
    if (pedalState[i] == HIGH && prevPedalState[i] == LOW) {
      Serial.write(i | 0x80);
      delay(10); // debounce
    }

    // Transition: HIGH -> LOW (RELEASE)
    // MSB: 0
    if (pedalState[i] == LOW && prevPedalState[i] == HIGH) {
      Serial.write(i);
      delay(10); // debounce
    }

    prevPedalState[i] = pedalState[i];
  }
}

// The protocol:
// Arduino sends a single byte when a pedal is pressed or released 
// The byte is used like: first bit (MSB) for the event type (1=press, 0=release), following bits for pedal ID (in decimal: 0 to NUM_PEDALS-1)
// Example: 00000010 = Pedal 2 released
//          10000001 = Pedal 1 pressed
// This means the maximum number of pedals supported is 128 (0-127).