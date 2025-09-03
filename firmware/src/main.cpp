#include <Arduino.h>

#include "Display.h"
#include "Net.h"

void setup() {
    Serial.begin(115200);
    while (!Serial) ;

    if (!displayBegin()) {
        Serial.println("Unable to initialize display");
        for(;;);
    }

    showConnecting();
    netBegin();
}

void loop() {
    if (frameAvailable()) {
        // TODO
    }
    displayTick();
    delay(10);
}
