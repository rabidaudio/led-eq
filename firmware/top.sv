`include "ResetGenerator.sv"
`include "LEDMatrix.sv"

module top #(
    parameter RESET_AFTER = 'hFFFF
) (
    input clk,
    output logic [15:0] hub_75_o
);
    logic reset;

    ResetGenerator #(.AFTER(RESET_AFTER)) reset_gen (.clk(clk), .reset(reset));

    LEDMatrix_32x16_1to8 matrix (
        .clk(clk),
        .reset(reset),
        .brightness(8'h7f), // half brightness
        .hub_75(hub_75_o)
    );
endmodule
