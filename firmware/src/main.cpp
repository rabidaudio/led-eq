#include <Arduino.h>

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

Adafruit_Protomatter matrix(
  DISPLAY_WIDTH, // Width of matrix (or matrix chain) in pixels
  5,             // Bit depth, 1-6
  1, rgbPins,    // # of matrix chains, array of 6 RGB pins for each
  3, addrPins,   // # of address pins (height is inferred), array of pins
  clockPin, latchPin, oePin, // Other matrix control pins
  true);         // Enable double-buffering

#define HUE_RANGE 0
#define HUE_STEP 32

uint16_t hue;

double displayRadius = sqrt(DISPLAY_WIDTH*DISPLAY_WIDTH + DISPLAY_HEIGHT*DISPLAY_HEIGHT);


void setup() {
    Serial.begin(115200);
    while (!Serial) ;

    ProtomatterStatus status = matrix.begin();
    Serial.print("Protomatter begin() status: ");
    Serial.println((int)status);
    if(status != PROTOMATTER_OK) {
        // DO NOT CONTINUE if matrix setup encountered an error.
        for(;;);
    }
}

void loop() {    
    for (int x = 0; x < DISPLAY_WIDTH; x++) {
        for (int y = 0; y < DISPLAY_HEIGHT; y++) {
            double radius = sqrt(x*x + y*y);
            double ratio = 1.0 - radius/displayRadius;
            uint16_t h = hue + (uint16_t)(ratio*HUE_RANGE);
            uint8_t s = 255;
            uint8_t v = 255;
            uint16_t color = matrix.colorHSV(h, s, v);
            matrix.writePixel(DISPLAY_WIDTH-x-1, y, color);
        }
    }
    matrix.show();
    

    hue += HUE_STEP;
    Serial.println(hue);
    delay(10);
}

// #include <WiFi.h>

// const char* WIFI_NET = "The Coven";
// const char* WIFI_PASS = "heatfromfire";

// WiFiServer server = WiFiServer(23);

// void setup() {
//     pinMode(LED_BUILTIN, OUTPUT);

//     Serial.begin(115200);
//     while (!Serial) ;

//     WiFi.begin(WIFI_NET, WIFI_PASS);
//     bool led = true;
//     while (WiFi.status() != WL_CONNECTED) 
//     {
//         // Blink LED while we're connecting:
//         digitalWrite(LED_BUILTIN, led ? HIGH : LOW);
//         led = !led;
//         delay(100);
//         Serial.print(".");
//     }
//     digitalWrite(LED_BUILTIN, LOW);
//     Serial.println();
//     Serial.println("WiFi connected!");
//     Serial.print("IP address: ");
//     Serial.println(WiFi.localIP());

//     server.begin();
//     Serial.println("Server running.");
// }

// void loop() {
//     WiFiClient client = server.available();   // listen for incoming clients

//     if (client) {                             // if you get a client,
//         Serial.println("New Client.");           // print a message out the serial port
//         String currentLine = "";                // make a String to hold incoming data from the client
//         while (client.connected()) {            // loop while the client's connected
//             if (client.available()) {             // if there's bytes to read from the client,
//                 char c = client.read();             // read a byte, then
//                 Serial.write(c);                    // print it out the serial monitor
//             }
//         }
//         client.stop();
//         Serial.println("Client Disconnected.");
//     }
// }
