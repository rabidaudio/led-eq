`include "LEDMatrix.sv"

module LEDMatrix_Test (
    input t_clk,
    input t_reset
);
    logic [1:0] t_row_addr;
    logic [3:0] t_pixel_addr;

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
        .DWELL_CYCLES(16),
        .CLOCK_DIVIDER(4)
    ) dut (
        .clk(t_clk),
        .reset(t_reset),

        .row_addr(t_row_addr),
        .pixel_addr(t_pixel_addr),
        .r_data(t_pixel_addr[1:0]),
        .g_data(t_pixel_addr[2:1]),
        .b_data('1),

        .brightness(5'h08), // half brightness
        // .brightness(16),
        // .brightness(0),

        .red(t_parallel[5:4]),
        .green(t_parallel[3:2]),
        .blue(t_parallel[1:0]),
        .row_select(t_row_select),
        .lat(t_lat),
        .oe_n(t_oe_n)
    );
endmodule
