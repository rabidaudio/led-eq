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
    // GoldenMonitor #(.DELAY(9), .WIDTH(8)) gm_data (
    //     .clk(t_clk),
    //     .enable('1),
    //     .golden(g_data),
    //     .signal(t_data)
    // );

    logic [7:0] data_in [$];
    logic [7:0] data_out [$];

    initial begin
        logic [2:0] b;

        t_mosi <= 0;
        b <= 7;
        g_data <= 0;
        t_cs <= 1;
        @(negedge t_reset);
        repeat (10) @(posedge t_clk);

        for (int j = 0; j < 256; j++) begin
            
            for (int i = 0; i < 8; i++) begin
                @(posedge t_clk);
                t_cs <= 0;
                t_mosi <= g_data[b];
                b <= b - 1;
            end
            data_in.push_back(g_data);
            g_data <= g_data + 1;

            if (g_data % 16 == 0) begin
                @(posedge t_clk);
                t_cs <= 1;
                repeat (31) @(posedge t_clk);
            end
        end

        while (data_in.size() && data_out.size()) begin
            logic [7:0] expected;
            logic [7:0] actual;

            expected = data_in.pop_front();
            actual = data_out.pop_front();
            if (expected != actual)
                $error("test failed: [%d] expected %d but was %d", 256-data_in.size(), expected, actual);
        end
    end

    always_ff @(posedge t_clk) begin
        if (t_data_avail) data_out.push_back(t_data);
    end
endmodule
