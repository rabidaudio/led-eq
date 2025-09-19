`timescale 1ns/10ps
`include "SpiSlaveSimplex.sv"
`include "utils.test.sv"

module SpiSlaveSimplex_Tests(
    input t_clk,
    input t_reset
);

    SpiSlaveSimplex_TestCounter sss_test_counter (t_clk, t_reset);

endmodule

module SpiSlaveSimplex_TestCounter (
    input t_clk,
    input t_reset
);

    logic t_cs;
    logic t_mosi;
    logic [7:0] t_data;
    logic t_data_avail;
    logic [5:0] t_byte_index;
    logic test_reset;

    SpiSlaveSimplex #(.BYTE_COUNTER_MAX(64)) test_spi_slave_simplex (
        .sck(t_clk),
        .reset(t_reset),
        .cs(t_cs),
        .mosi(t_mosi),
        .data(t_data),
        .data_avail(t_data_avail),
        .byte_index(t_byte_index)
    );

    logic [7:0] g_data;
    GoldenMonitor #(.DELAY(9), .WIDTH(8)) gm_data (
        .clk(t_clk),
        .enable('1),
        .golden(g_data),
        .signal(t_data)
    );

    initial begin
        logic [2:0] b;

        t_mosi <= 0;
        b <= 7;
        g_data <= 0;
        t_cs <= 1;
        @(negedge t_reset);
        repeat (10) @(posedge t_clk);

        while (1) begin
            logic [7:0] g_data_next;

            
            for (int i = 0; i < 8; i++) begin
                @(posedge t_clk);
                t_cs <= 0;
                t_mosi <= g_data[b];
                b <= b - 1;
            end

            //g_data_next = g_data + 1;
            //t_mosi <= g_data[b];
            g_data <= g_data + 1;

            // if (g_data % 8 == 0) begin
            //     @(posedge t_clk);
            //     t_cs <= 1; 
            //     repeat (31) @(posedge t_clk);
            // end

        end
    end
endmodule
