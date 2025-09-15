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

