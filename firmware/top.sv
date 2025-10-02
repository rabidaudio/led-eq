`include "ResetGenerator.sv"
`include "LEDMatrix.sv"
`include "PixelBus.sv"
`include "MockCanvas.sv"

module top #(
    parameter RESET_AFTER = 'hFFFF
) (
    input clk,
    output logic [15:0] hub_75_o
);
    logic reset;

    ResetGenerator #(.AFTER(RESET_AFTER)) reset_gen (.clk(clk), .reset(reset));

    PixelBus #(
        .WIDTH(32), 
        .HEIGHT(16),
        .SCAN_RATE(8),
        // it's actually 1 but we'll accept higher bit depths and always return 1
        .BIT_DEPTH(8)
    ) bus (.clk(clk));

    MockCanvas cv (.reset(reset), .bus(bus.Canvas));

    LEDMatrix_32x16_1to8 matrix (
        .clk(clk),
        .reset(reset),
        // .brightness(4'h8), // half brightness
        .hub_75(hub_75_o)
    );
endmodule
