`ifndef __PIXELBUS__
`define __PIXELBUS__

/**
 * `PixelBus` is an interface that `LEDMatrix` uses to access pixel data to display.
 * The matrix will request the `bitplane_addr`, `row_addr`, and `pixel_addr` desired and
 * pull `req_read` high.
 * On the next clock cycle, it will be placed on `r_data`, `g_data`, and `b_data`.
 */
interface PixelBus #(
    parameter WIDTH = 64, // in pixels
    parameter HEIGHT = 64, // in pixels
    parameter SCAN_RATE = 8,
    parameter BIT_DEPTH = 8
) (input clk);
    localparam COLOR_WIDTH = HEIGHT/SCAN_RATE;

    logic req_read;
    logic [$clog2(BIT_DEPTH-1)-1:0] bitplane_addr;
    logic [$clog2(SCAN_RATE-1)-1:0] row_addr;
    logic [$clog2(WIDTH-1)-1:0] pixel_addr;

    logic [COLOR_WIDTH-1:0] r_data;
    logic [COLOR_WIDTH-1:0] g_data;
    logic [COLOR_WIDTH-1:0] b_data;

    modport Canvas (
        input req_read, input bitplane_addr, input row_addr, input pixel_addr,
        output read_ready, output r_data, output g_data, output b_data
    );

    modport Matrix (
        output req_read, output bitplane_addr, output row_addr, output pixel_addr,
        input read_ready, input r_data, input g_data, input b_data
    );
endinterface

`endif // __PIXELBUS__
