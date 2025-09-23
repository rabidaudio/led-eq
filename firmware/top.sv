`include "ResetGenerator.sv"
`include "LEDMatrix.sv"

module ClockDivider #(
    parameter WIDTH = 2
) (input clk, input reset, output logic slow_clk);

    logic [WIDTH-1:0] counter;

    always_ff @(posedge clk) begin
        if (counter == 0) begin
            slow_clk <= 1;
        end else begin
            counter <= counter - 1;
            slow_clk <= 0;
        end

        if (reset) begin
            slow_clk <= 0;
            counter <= 0;
        end
    end
endmodule

module top (
    input clk,
    output [15:0] hub_75_o
);
    logic reset;

    ResetGenerator reset_gen (.clk(clk), .reset(reset));

    // logic led_clk;

    ClockDivider #(.WIDTH(22)) clk_div (.clk(clk), .reset(reset), .slow_clk(hub_75_o[12]));

    LEDMatrix_32x16_1to8 matrix (
        .clk(hub_75_o[12]),
        .reset(reset),
        .hub_75(hub_75_o)
    );

    // always_comb hub_75_o[12] = led_clk; // STOPSHIP

    always_comb begin : ground
        hub_75_o[3] <= 0;
        hub_75_o[7] <= 0;
        hub_75_o[11] <= 0;
        hub_75_o[15] <= 0;
    end
endmodule
