#include <Adafruit_Protomatter.h>

#define DISPLAY_WIDTH 32
#define DISPLAY_HEIGHT 16

// R0   G0
// B0   gnd
// R1   G1
// B1   gnd
// A    B
// C    gnd
// CLK  LAT/SCLK
// ~OE  gnd
uint8_t rgbPins[]  = {
    4, 12, 13, // R0 B0 G0
    14, 15, 21 // R1 B1 G1
};
uint8_t addrPins[] = {16, 17, 25}; // A B C
uint8_t clockPin   = 27; // Must be on same port as rgbPins
uint8_t latchPin   = 32;
uint8_t oePin      = 33;

#define HUE_RANGE 0
#define HUE_STEP 32

Adafruit_Protomatter matrix_(
                DISPLAY_WIDTH, // Width of matrix (or matrix chain) in pixels
                5,             // Bit depth, 1-6
                1, rgbPins,    // # of matrix chains, array of 6 RGB pins for each
                3, addrPins,   // # of address pins (height is inferred), array of pins
                clockPin, latchPin, oePin, // Other matrix control pins
                true);         // Enable double-buffering

uint16_t hue_;
uint16_t step_;
uint16_t HUE_OFFSETS[DISPLAY_WIDTH*DISPLAY_HEIGHT];

bool displayBegin(){ 
    // compute hue offsets
    double displayRadius = sqrt(DISPLAY_WIDTH*DISPLAY_WIDTH + DISPLAY_HEIGHT*DISPLAY_HEIGHT);
    for (int x = 0; x < DISPLAY_WIDTH; x++) {
        for (int y = 0; y < DISPLAY_HEIGHT; y++) {
            double radius = sqrt(x*x + y*y);
            double ratio = 1.0 - radius/displayRadius;
            int i = y*DISPLAY_HEIGHT + (DISPLAY_WIDTH-x-1);
            HUE_OFFSETS[i] = (uint16_t)(ratio*HUE_RANGE);
        }
    }

    ProtomatterStatus status = matrix_.begin();
    return status == PROTOMATTER_OK;
}

void showConnecting() {
    matrix_.println("WIFI");
    matrix_.show();
}

void displayTick() {
    for (int x = 0; x < DISPLAY_WIDTH; x++) {
        for (int y = 0; y < DISPLAY_HEIGHT; y++) {
            uint16_t color = matrix_.colorHSV(hue_ + HUE_OFFSETS[y*DISPLAY_HEIGHT + x]);
            matrix_.writePixel(x, y, color);
        }
    }
    matrix_.show();
    
    hue_ += HUE_STEP;
    step_ += 1;
    if (step_ >= 100) {
        Serial.print("Refresh FPS = ~");
        Serial.println(matrix_.getFrameCount());
        step_ = 0;
    }
}