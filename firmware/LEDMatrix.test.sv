`include "LEDMatrix.sv"

module LEDMatrix_Test (
    input t_clk,
    input t_reset
);

    logic [5:0] t_parallel;
    logic [2:0] t_row_select;
    logic t_out_clk;
    logic t_lat;
    logic t_oe_n;

    LEDMatrix #(
        .WIDTH(16),
        .HEIGHT(32),
        .SCAN_RATE(8),
        .COLOR_WIDTH(2)
    ) dut (
        .clk(t_clk),
        .reset(t_reset),
        .red(t_parallel[5:4]),
        .green(t_parallel[3:2]),
        .blue(t_parallel[1:0]),
        .row_select(t_row_select),
        // .out_clk(t_out_clk),
        .lat(t_lat),
        .oe_n(t_oe_n)
    );
endmodule
