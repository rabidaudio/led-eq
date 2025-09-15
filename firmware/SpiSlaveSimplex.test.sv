`timescale 1ns/10ps
`include "SpiSlaveSimplex.sv"

module Test_SpiSlaveSimplex(
    input test_spi_clk
);
    logic test_cs;
    logic test_mosi;
    logic [7:0] test_data_out;
    logic test_data_ready;
    logic [5:0] test_byte_index;

    SpiSlaveSimplex #(.BYTE_INDEX_SIZE(5)) test_spi_slave_simplex(
        .sck(test_spi_clk),
        .cs(test_cs),
        .mosi(test_mosi),
        .data(test_data_out),
        .data_avail(test_data_ready),
        .byte_index(test_byte_index)
    );

    logic [7:0] test_data_in;

    initial begin
        #100;
        @(posedge test_spi_clk);
        test_cs <= 1;
        test_mosi <= 0;

        test_data_in <= 0;

        // send 32 byte chunks
        for (int i = 0; i < 8; i++) begin // each chunk
            // start transaction
            for (int j = 0; j < 32; j++) begin // each byte
                for (int k = 8; k-- > 0;) begin // each bit
                    @(posedge test_spi_clk);
                    test_cs <= 0;
                    test_data_in <= (i << 3) + j;
                    test_mosi <= test_data_in[k];
                end
            end

            // stop transaction
            @(posedge test_spi_clk);
            test_cs <= 1;

            #1000; // wait a bit
        end
        $finish;
    end

endmodule

