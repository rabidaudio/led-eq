#include <Arduino.h>

#include <WiFi.h>

const char* WIFI_NET = "The Coven";
const char* WIFI_PASS = "heatfromfire";

WiFiServer server = WiFiServer(23);

void setup() {
    pinMode(LED_BUILTIN, OUTPUT);

    Serial.begin(115200);
    while (!Serial) ;

    WiFi.begin(WIFI_NET, WIFI_PASS);
    bool led = true;
    while (WiFi.status() != WL_CONNECTED) 
    {
        // Blink LED while we're connecting:
        digitalWrite(LED_BUILTIN, led ? HIGH : LOW);
        led = !led;
        delay(100);
        Serial.print(".");
    }
    digitalWrite(LED_BUILTIN, LOW);
    Serial.println();
    Serial.println("WiFi connected!");
    Serial.print("IP address: ");
    Serial.println(WiFi.localIP());

    server.begin();
    Serial.println("Server running.");
}

void loop() {
    WiFiClient client = server.available();   // listen for incoming clients

    if (client) {                             // if you get a client,
        Serial.println("New Client.");           // print a message out the serial port
        String currentLine = "";                // make a String to hold incoming data from the client
        while (client.connected()) {            // loop while the client's connected
            if (client.available()) {             // if there's bytes to read from the client,
                char c = client.read();             // read a byte, then
                Serial.write(c);                    // print it out the serial monitor
            }
        }
        client.stop();
        Serial.println("Client Disconnected.");
    }
}
