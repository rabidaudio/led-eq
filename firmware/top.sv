`include "ResetGenerator.sv"
`include "LEDMatrix.sv"
`include "BrightnessRamp.sv"
`include "FlagCanvas.sv"
`include "uart.sv"

module top #(
    parameter RESET_AFTER = 'hFFFF,
    parameter BIT_DEPTH = 8
) (
    input clk,
    output logic [15:0] hub_75_o,

    // uart connection
    output uart_tx_o,
    input uart_rx_i
);
    logic reset;

    ResetGenerator #(.AFTER(RESET_AFTER)) reset_gen (.clk(clk), .reset(reset));

    logic req_read;
    logic [$clog2(BIT_DEPTH-1)-1:0] bitplane_addr;
    logic [2:0] row_addr;
    logic [4:0] pixel_addr;
    logic read_ready;
    logic [1:0] r_data;
    logic [1:0] g_data;
    logic [1:0] b_data;

    FlagCanvas cv (
        .clk(clk),
        .reset(reset),

        .req_read(req_read),
        .bitplane_addr(bitplane_addr),
        .row_addr(row_addr),
        .pixel_addr(pixel_addr),
        .r_data(r_data),
        .g_data(g_data),
        .b_data(b_data)
    );

    logic frame_end;
    logic [4:0] brightness;

    BrightnessRamp ramp (
        .clk(clk),
        .reset(reset),
        .inc(frame_end),
        .brightness(brightness)
    );

    LEDMatrix_32x16_1to8 #(.BIT_DEPTH(8)) matrix (
        .clk(clk),
        .reset(reset),

        // pixel bus
        .req_read(req_read),
        .bitplane_addr(bitplane_addr),
        .row_addr(row_addr),
        .pixel_addr(pixel_addr),
        .r_data(r_data),
        .g_data(g_data),
        .b_data(b_data),

        .frame_complete(frame_end),
        .brightness(brightness),
        // .brightness(4'h07),

        .hub_75(hub_75_o)
    );

    // TODO: hook me up.
    // This block of hardware recieves bytes over UART.
    // Whenever a new frame starts, a '\r' byte is sent over the UART TX. The program on the other
    // end should use this signal to figure out when to start sending a new frame.
    // At the start of a new frame, data reception is restarted from the top-left corner of the
    // display.
    // Image bytes should be sent over UART in xbgr format, where 'x' is a don't care bit.
    // Each byte contains a single bitplane from 2 pixels - one from the bottom half of the screen
    // and the other from the top half.
    logic [7:0] rx_byte;
    logic [15:0] rx_index;
    logic rx_byte_valid;
    localparam UART_BAUD = 2_000_000;
    uart_tx #(.CLK_DIVIDER(12_000_000/UART_BAUD)) uart_vsync (
        .clk_i(clk), .reset_i(reset),
        .txdata_valid_i(frame_end), .txdata_i(8'h0d),
        .uart_busy_o(), .uart_tx_o
    );

    UartPixelReceiver #(.CLK_DIVIDER(12_000_000/UART_BAUD)) uart_receiver (
        .clk_i(clk), .reset_i(reset),
        .uart_rx_i, .led_vsync(frame_end),
        .rx_byte, .rx_index, .rx_byte_valid
    );
endmodule


module UartPixelReceiver #(
    parameter CLK_DIVIDER = 104
)  (
    input clk_i,
    input reset_i,

    // uart connection
    input uart_rx_i,

    // vsync signal. Tells us when we should reset the pixel index.
    input led_vsync,

    // When this signal is strobed high, we have a new byte that must be read.
    output logic rx_byte_valid,

    // Data out to pixel memory.
    // We expect to recieve "pre-scrambled" data.
    output logic [7:0] rx_byte,

    // index for rx'd data
    output logic [15:0] rx_index,
);
    logic uart_rxdata_valid;
    logic [7:0] uart_rxdata;
    uart_rx #(.CLK_DIVIDER(CLK_DIVIDER)) uart_rx_inst (
        .clk_i, .reset_i,
        .uart_rx_i,
        .rxdata_valid_o(uart_rxdata_valid), .rxdata(uart_rxdata), .uart_busy_o()
    );

    always_ff @(posedge clk_i) begin
        if (uart_rxdata_valid) begin
            rx_byte_valid <= 1;
            rx_byte <= uart_rxdata;
            rx_index <= rx_index + 1;
        end else begin
            rx_byte_valid <= 0;
            rx_byte <= 'x;
        end

        if (led_vsync) rx_index <= 0;

        if (reset_i) begin
            rx_byte_valid <= 0;
            rx_byte <= 'x;
            rx_index <= 0;
        end
    end
endmodule
