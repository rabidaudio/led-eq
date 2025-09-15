`timescale 1ns/10ps
`include "SpiSlaveSimplex.sv"

module Test_SpiSlaveSimplex(
    input test_spi_clk
);
    logic test_cs;
    logic test_mosi;
    logic [7:0] test_data_out;
    logic test_data_avail;
    logic [5:0] test_byte_index;

    SpiSlaveSimplex #(.BYTE_INDEX_SIZE(5)) test_spi_slave_simplex(
        .sck(test_spi_clk),
        .cs(test_cs),
        .mosi(test_mosi),
        .data(test_data_out),
        .data_avail(test_data_avail),
        .byte_index(test_byte_index)
    );

    logic [7:0] test_data_in;

    initial begin : test_spi_slave_msb_active_low
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
                    test_data_in <= (i << 2) + j;
                    test_mosi <= test_data_in[k];
                end
            end

            // stop transaction
            @(posedge test_spi_clk);
            test_cs <= 1;

            #1000; // wait a bit
        end
    end



    // TODO: not supported by iverilog:
    // perhaps we could make a macro for this?
    /*
    property assert_data_counts_up;
    @(posedge test_spi_clk) test_data_in |-> ##[9] test_data_out;
    endproperty
    assert property (assert_data_counts_up) else $error("invalid data");
    */
    initial begin : assert_data_counts_up
        localparam DELAY = 9;
        localparam WIDTH = DELAY*8-1;
        logic [WIDTH:0] delay_line;
        
        // populate queue
        for (int i = 0; i < DELAY; i++) begin
            @(posedge test_spi_clk);
            delay_line <= { delay_line[WIDTH-8:0], test_data_in };
        end

        while (!test_cs) begin // TODO: only testing first transaction
            @(posedge test_spi_clk);
            if (test_data_out != delay_line[WIDTH:WIDTH-7]) begin
                $error("data output invalid. expected %h but was %h", delay_line[WIDTH:WIDTH-7], test_data_out);
            end
            delay_line <= { delay_line[(DELAY-1)*8-1:0], test_data_in };
        end
    end

    initial begin : assert_data_retained_post_transaction
        logic [7:0] prev_value;

        @(posedge test_cs); // ignore initial state
        while (1) begin
            @(posedge test_cs);
            prev_value <= test_data_out;
            @(posedge test_spi_clk);
            if (test_data_out != prev_value) $error("data not retained after cs low");
        end
    end

endmodule
