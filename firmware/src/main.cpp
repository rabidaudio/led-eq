#include <Arduino.h>

#include <Adafruit_GFX.h>
#include <DMD_RGB.h>

#define WIDTH 32
#define HEIGHT 16

#define ENABLE_DUAL_BUFFER true

// mux pins - A, B, C... all mux pins must be selected from same port!
#define DMD_PIN_A PB6
#define DMD_PIN_B PB5
#define DMD_PIN_C PB4
uint8_t mux_list[] = { DMD_PIN_A , DMD_PIN_B , DMD_PIN_C };

// pin OE must be one of PB0 PB1 PA6 PA7
#define DMD_PIN_nOE PB0
#define DMD_PIN_SCLK PB7

// Pins for R0, G0, B0, R1, G1, B1 channels and for clock.
// By default the library uses RGB color order.
// If you need to change this - reorder the R0, G0, B0, R1, G1, B1 pins.
// All this pins also must be selected from same port!
uint8_t custom_rgbpins[] = { PA15, PA0,PA1,PA2,PA3,PA4,PA5 }; // CLK, R0, G0, B0, R1, G1, B1

DMD_RGB<RGB32x16plainS8, COLOR_4BITS> dmd(mux_list, DMD_PIN_nOE, DMD_PIN_SCLK, custom_rgbpins, 1, 1, ENABLE_DUAL_BUFFER);

void setup() {
    pinMode(PC13, OUTPUT);
    
    Serial.begin(115200);
    while(!Serial);

    Serial.println("Hello!");

    dmd.init();
    dmd.setBrightness(50); // recommended 30-100

    // Serial.println("display initialized!");
    digitalWrite(PC13, HIGH);
}

uint8_t i = 0;
void loop() {
    uint16_t color = dmd.Color888((i%3) == 0 ? 255 : 0, (i%3) == 1 ? 255 : 0, (i%3) == 2 ? 255 : 0);
    dmd.clearScreen(true);
    dmd.drawFilledBox(0, 0, 4, 4, color);
    dmd.swapBuffers(true);

    delay(1000);
    i += 1;
}