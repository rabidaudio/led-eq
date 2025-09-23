`include "ResetGenerator.sv"
`include "LEDMatrix.sv"

module top (
    input clk,
    output [15:0] hub_75_o
);
    logic reset;

    ResetGenerator reset_gen (.clk(clk), .reset(reset));

    LEDMatrix matrix (
        .clk(clk),
        .reset(reset),
        .red({ hub_75_o[0], hub_75_o[4] }),
        .green({ hub_75_o[1], hub_75_o[5] }),
        .blue({ hub_75_o[2], hub_75_o[6] }),
        .row_select(hub_75_o[10:8]),
        .out_clk(hub_75_o[12]),
        .lat(hub_75_o[13]),
        .oe(hub_75_o[14])
    );

    always_comb begin : ground
        hub_75_o[3] <= 0;
        hub_75_o[7] <= 0;
        hub_75_o[11] <= 0;
        hub_75_o[15] <= 0;
    end
endmodule
