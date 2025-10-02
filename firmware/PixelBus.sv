`ifndef __PIXELBUS__
`define __PIXELBUS__

/**
 * `PixelBus` is an interface that `LEDMatrix` uses to access pixel data to display.
 * The matrix will request the `bitplane`, `row_addr`, and `pixel_addr` desired and
 * pull `req_read` high.
 * When the data is available, it will be placed on `red`, `green`, and `blue`, and
 * `read_ready` will go high.
 */
interface PixelBus #(
    parameter WIDTH = 64, // in pixels
    parameter HEIGHT = 64, // in pixels
    parameter SCAN_RATE = 8,
    parameter BIT_DEPTH = 8
) (input clk);
    localparam COLOR_WIDTH = HEIGHT/SCAN_RATE;

    logic req_read;
    logic [$clog2(BIT_DEPTH-1)-1:0] bitplane;
    logic [$clog2(SCAN_RATE-1)-1:0] row_addr;
    logic [$clog2(WIDTH-1)-1:0] pixel_addr;

    logic read_ready;
    logic [COLOR_WIDTH-1:0] red;
    logic [COLOR_WIDTH-1:0] green;
    logic [COLOR_WIDTH-1:0] blue;

    modport Canvas (
        input req_read, input bitplane, input row_addr, input pixel_addr,
        output read_ready, output red, output green, output blue
    );

    modport Matrix (
        output req_read, output bitplane, output row_addr, output pixel_addr,
        input read_ready, input red, input green, input blue
    );
endinterface

`endif // __PIXELBUS__
