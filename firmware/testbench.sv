`timescale 1ns/10ps
`include "SpiSlaveSimplex.test.sv"

module testbench;
    localparam realtime PERIOD = 10;
    localparam TEST_CYCLES = 4*16*16*256 + 1000;
    
    logic t_clk = 0;

    always #(PERIOD/2) t_clk <= ~t_clk;
    initial #(TEST_CYCLES*PERIOD/2) $finish;

    logic t_reset;

    initial begin
        @(posedge t_clk);
        t_reset <= 1;
        @(posedge t_clk);
        t_reset <= 0;
        @(posedge t_clk);
    end

    initial begin
        $dumpfile("build/results.vcd");
        $dumpvars(0, testbench);
    end

    SpiSlaveSimplex_Tests sss_tests (t_clk, t_reset);
endmodule
