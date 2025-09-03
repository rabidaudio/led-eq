#include <WiFi.h>

const char* WIFI_NET = "The Coven";
const char* WIFI_PASS = "heatfromfire";

#define PORT 23

WiFiServer server_;
WiFiClient client_; // TODO: can only support one connection at a time

void netBegin() {
    pinMode(LED_BUILTIN, OUTPUT);
    server_ = WiFiServer(PORT);

    WiFi.begin(WIFI_NET, WIFI_PASS);
    WiFi.setAutoReconnect(true);
    bool led = true;
    wl_status_t status = WL_IDLE_STATUS;
    while (status != WL_CONNECTED) 
    {
        wl_status_t newStatus = WiFi.status();
        if (newStatus != status) {
            switch (newStatus) {
                case WL_IDLE_STATUS: Serial.println("idle"); break;
                case WL_NO_SSID_AVAIL: Serial.println("ssid not found"); break;
                case WL_SCAN_COMPLETED: Serial.println("scan completed"); break;
                case WL_CONNECTED: Serial.println("connected"); break;
                case WL_CONNECT_FAILED: Serial.println("connection failed"); break;
                case WL_CONNECTION_LOST: Serial.println("connection lost"); break;
                case WL_DISCONNECTED: Serial.println("disconnected"); break;
                default:
                    Serial.println(newStatus, DEC);
            }
            status = newStatus;
        }
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

    server_.begin();
    Serial.println("Server running.");
}

bool frameAvailable() {
    if (!client_) {
        client_ = server_.available();
    } else if (!client_.connected()) {
        client_.stop();
    }
    if (!client_) {
        return false; // no one is connected
    }

    while (client_.available()) {
        // TODO: state machine to parse message
        char c = client_.read();
        Serial.write(c);
    }
    return true;
}
