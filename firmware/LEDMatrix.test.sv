`include "LEDMatrix.sv"

module LEDMatrix_Test (
    input t_clk,
    input t_reset
);

    logic [5:0] t_parallel;
    logic [2:0] t_row_select;
    logic t_out_clk;
    logic t_lat;
    logic t_oe;

    LEDMatrix #(.CLOCK_DIV_WIDTH(1)) dut (
        .clk(t_clk),
        .reset(t_reset),
        .red(t_parallel[5:4]),
        .green(t_parallel[3:2]),
        .blue(t_parallel[1:0]),
        .row_select(t_row_select),
        .out_clk(t_out_clk),
        .lat(t_lat),
        .oe(t_oe)
    );
endmodule
