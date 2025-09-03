/*
#include <SPI.h>

// #define USE_W5100               false // using W5500

// #define USING_CUSTOM_SPI false
// #define USING_CUSTOM_SPI            true
// #define USING_SPI2                  true
// #define CUR_PIN_MISO              PA6
// #define CUR_PIN_MOSI              PA7
// #define CUR_PIN_SCK               PA5
// #define CUR_PIN_SS                PA4

// #define SPI_NEW_INITIALIZED       true
// SPIClass SPI_New(CUR_PIN_MOSI, CUR_PIN_MISO, CUR_PIN_SCK);
// #include <Ethernet_Generic.h>

#define ETHERNET_LARGE_BUFFERS

#include <Ethernet.h>

uint8_t MAC_ADDRESS[6] = { 'C', 'O', 'V', 'A', 'E', 'Q' };
EthernetServer server = EthernetServer(1000);

void connect() {
    Ethernet.init(PA4);
    while (true) {
        if (Ethernet.begin(MAC_ADDRESS) == 0) {
            Serial.println("unable to acquire IP!");
            if (Ethernet.hardwareStatus() == EthernetNoHardware) {
                Serial.println("W5x00 chip not found");
            } else if (Ethernet.linkStatus() == LinkOFF) {
                Serial.println("Ethernet cable is not connected");
            }
            
            delay(1000);
        } else {
            break;
        }
    }
    Serial.print("Connected! ");
    Serial.println(Ethernet.localIP());

    server.begin();
}
    */
   