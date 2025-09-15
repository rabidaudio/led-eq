`timescale 1ns/10ps
`include "SpiSlaveSimplex.test.sv"

module main;
    logic test_clk;

    localparam TEST_LENGTH = 5*1000*1000; // 5ms @ 1ns

    Test_SpiSlaveSimplex test1(.test_spi_clk(test_clk));

    initial begin
        $dumpfile("build/results.vcd");
        $dumpvars(0, main);
    end

    initial begin
        test_clk = 0;
        for (int i = 0; i < TEST_LENGTH; i++) begin
            #1000; // 1MHz
            test_clk <= !test_clk;
        end
        $finish;
    end 
endmodule
