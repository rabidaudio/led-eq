`timescale 1ns/10ps
`include "SpiSlaveSimplex.sv"
`include "utils.test.sv"

module Test_SpiSlaveSimplex(
    input test_spi_clk
);
    logic test_cs;
    logic test_mosi;
    logic [7:0] test_data_out;
    logic test_data_avail;
    logic [5:0] test_byte_index;
    logic test_reset;

    SpiSlaveSimplex #(.BYTE_COUNTER_MAX(64)) test_spi_slave_simplex (
        .sck(test_spi_clk),
        .reset(test_reset),
        .cs(test_cs),
        .mosi(test_mosi),
        .data(test_data_out),
        .data_avail(test_data_avail),
        .byte_index(test_byte_index)
    );

    logic [7:0] test_data_in;
    GoldenMonitor #(.DELAY(10), .WIDTH(8)) gm_data (
        .clk(test_spi_clk),
        .enable('1),
        .golden(test_data_in),
        .signal(test_data_out)
    );

    initial begin
        logic [2:0] b;

        @(posedge test_spi_clk);
        test_reset <= 1;
        test_mosi <= 0;
        test_data_in <= 0;
        b <= 7;
        @(posedge test_spi_clk);
        test_reset <= 0;
        test_data_in <= 0;
        test_cs <= 1;

        while (1) begin
            for (int i = 0; i < 8; i++) begin
                @(posedge test_spi_clk);
                test_cs <= 0;
                test_mosi <= test_data_in[b];
                b <= b - 1;
            end
            test_data_in <= test_data_in + 1;
        end
    end

    // initial begin : test_spi_slave_msb_active_low
    //     repeat(4) @(posedge test_spi_clk); // wait a bit

    //     @(posedge test_spi_clk);
    //     test_cs <= 1;
    //     test_mosi <= 0;

    //     test_data_in <= 0;

    //     // send 32 byte chunks
    //     for (int i = 0; i < 8; i++) begin // each chunk
    //         // start transaction
    //         for (int j = 0; j < 32; j++) begin // each byte
    //             for (int k = 8; k-- > 0;) begin // each bit
    //                 @(posedge test_spi_clk);
    //                 test_cs <= 0;
    //                 test_data_in <= (i << 5) + j;
    //                 test_mosi <= test_data_in[k];
    //             end
    //         end

    //         // stop transaction
    //         @(posedge test_spi_clk);
    //         test_cs <= 1;

    //         repeat(10) @(posedge test_spi_clk); // wait a bit
    //     end
    // end

    // initial begin : assert_data_counts_up
    //     localparam DELAY = 9;
    //     logic [7:0] delay_line [$];
        
    //     // populate queue
    //     for (int i = 0; i < DELAY; i++) begin
    //         @(posedge test_spi_clk);
    //         delay_line.push_front(test_data_in);
    //     end

    //     while (!test_cs) begin // TODO: only testing first transaction
    //         logic [7:0] compare;

    //         @(posedge test_spi_clk);
    //         compare = delay_line.pop_back();
    //         if (test_data_out != compare) $error("data output invalid. expected %h but was %h", compare, test_data_out);
    //         delay_line.push_front(test_data_in);
    //     end
    // end

    // initial begin : assert_data_retained_post_transaction
    //     logic [7:0] prev_value;

    //     @(posedge test_cs); // ignore initial state
    //     while (1) begin
    //         @(posedge test_cs);
    //         prev_value <= test_data_out;
    //         @(posedge test_spi_clk);
    //         if (test_data_out != prev_value) $error("data not retained after cs low");
    //     end
    // end

endmodule
