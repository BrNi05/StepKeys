// The first digital pin a pedal is connected to
constexpr int FIRST_PIN = 2;

// The number of pedals connected
// So the last pin used will be FIRST_PIN + NUM_PEDALS - 1
constexpr int NUM_PEDALS = 5;

// Used to detect LOW → HIGH transitions
uint8_t prevPedalState[NUM_PEDALS];

void setup() {
  // Baud rate: 115200
  Serial.begin(115200);

  for (uint8_t i = 0; i < NUM_PEDALS; i++) {
    pinMode(FIRST_PIN + i, INPUT);
    prevPedalState[i] = LOW;
  }
}

void loop() {
  for (uint8_t i = 0; i < NUM_PEDALS; i++) {
    uint8_t state = digitalRead(FIRST_PIN + i);

    // Detect LOW → HIGH transition
    if (state == HIGH && prevPedalState[i] == LOW) {
      Serial.println(i);
    }

    prevPedalState[i] = state;
  }
}