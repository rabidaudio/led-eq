`include "LEDMatrix.sv"

module LEDMatrix_Test (
    input t_clk,
    input t_reset
);

    logic [5:0] t_parallel;
    logic [1:0] t_row_select;
    logic t_out_clk;
    logic t_lat;
    logic t_oe_n;

    LEDMatrix #(
        .WIDTH(16),
        .HEIGHT(8),
        .SCAN_RATE(2),
        .COLOR_WIDTH(2),
        .BIT_DEPTH(3),
        .CLOCK_DIVIDER(4)
    ) dut (
        .clk(t_clk),
        .reset(t_reset),
        .red(t_parallel[5:4]),
        .green(t_parallel[3:2]),
        .blue(t_parallel[1:0]),
        .row_select(t_row_select),
        .lat(t_lat),
        .oe_n(t_oe_n)
    );
endmodule
