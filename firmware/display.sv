`default_nettype none
`timescale 1ns/10ps

typedef enum { MSB, LSB } bit_order_e;
typedef enum { ACTIVE_LOW = 0, ACTIVE_HIGH = 1 } pin_direction_e;

/**
  * The SpiSlaveSimplex module is an SPI slave
  * with only read capabilities. Every time a full
  * byte is loaded into data, the data_ready pin
  * goes high for one cycle.
  * The clock is expected to be an external input
  * from the master device, so the system clock
  * will need to be greater than master_clk / 8.
  * byte_index will indicate the current byte
  * number since the start of the transaction.
  */
module SpiSlaveSimplex #(
    parameter ORDER = MSB,
    parameter CS_DIRECTION = ACTIVE_LOW,
    parameter BYTE_INDEX_SIZE = 11 // roll over at 2048 bytes
)(
    input master_clk,
    input cs,
    input mosi,
    output logic [7:0] data,
    output logic data_ready,
    output logic [BYTE_INDEX_SIZE:0] byte_index
);

    logic [2:0] bit_index;

    always_ff @(posedge master_clk)
        if (cs == CS_DIRECTION) begin
            // else load bit into data
            if (ORDER == MSB) data <= { data[6:0], mosi };
            else data <= { mosi, data[7:1] };

            if (bit_index == 7) begin
                data_ready <= 1;
                byte_index <= byte_index + 1;
            end else begin
                data_ready <= 0;
                byte_index <= byte_index;
            end
            
            bit_index <= bit_index + 1;
        end else begin
            // if cs disabled, reset
            data <= 0;
            data_ready <= 0;
            bit_index <= 0;
            byte_index <= 0;
        end
endmodule

module main;

    Test_SpiSlaveSimplex test1();

    initial begin
        $dumpfile("results.vcd");
        $dumpvars(0, main);
    end
endmodule

module Test_SpiSlaveSimplex;
    logic test_spi_clk;

    initial begin
        test_spi_clk = 0;
        while (1) begin
            #1000; // 1MHz
            test_spi_clk <= !test_spi_clk;
        end
    end

    logic test_cs;
    logic test_mosi;
    logic [7:0] test_data_out;
    logic test_data_ready;
    logic [5:0] test_byte_index;

    SpiSlaveSimplex #(.BYTE_INDEX_SIZE(5)) test_spi_slave_simplex(
        .master_clk(test_spi_clk),
        .cs(test_cs),
        .mosi(test_mosi),
        .data(test_data_out),
        .data_ready(test_data_ready),
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
