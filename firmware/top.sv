`include "ResetGenerator.sv"
`include "LEDMatrix.sv"
// `include "PixelBus.sv"
// `include "MockCanvas.sv"

module top #(
    parameter RESET_AFTER = 'hFFFF
) (
    input clk,
    output logic [15:0] hub_75_o
);
    logic reset;

    ResetGenerator #(.AFTER(RESET_AFTER)) reset_gen (.clk(clk), .reset(reset));

    // PixelBus #(
    //     .WIDTH(32), 
    //     .HEIGHT(16),
    //     .SCAN_RATE(8),
    //     // it's actually 1 but we'll accept higher bit depths and always return 1
    //     .BIT_DEPTH(8)
    // ) bus (.clk(clk));

    // PixelBus interface
    // logic req_read;
    // logic [$clog2(BIT_DEPTH-1)-1:0] bitplane;
    // logic [2:0] row_addr;
    // logic [4:0] pixel_addr;
    // logic read_ready;
    // logic [1:0] red;
    // logic [1:0] green;
    // logic [1:0] blue;
    
    // MockCanvas cv (.reset(reset), .bus(bus));

    LEDMatrix_32x16_1to8 #(.BIT_DEPTH(3)) matrix (
        .clk(clk),
        .reset(reset),

        // pixel bus
        // .req_read(req_read),
        // .bitplane(bitplane),
        // .row_addr(row_addr),
        // .pixel_addr(pixel_addr),
        // .read_ready(read_ready),
        // .red(red),
        // .green(green),
        // .blue(blue),

        // .brightness(4'h8), // half brightness

        .hub_75(hub_75_o)
    );
endmodule
