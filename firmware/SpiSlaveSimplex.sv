typedef enum { MSB, LSB } bit_order_e;
typedef enum { ACTIVE_LOW = 0, ACTIVE_HIGH = 1 } pin_direction_e;

/**
  * The SpiSlaveSimplex module is an SPI slave
  * with only read capabilities. Every time a full
  * byte is loaded into data, the data_avail pin
  * goes high for one cycle.
  * byte_index will indicate the current byte
  * number since the start of the transaction.
  */
module SpiSlaveSimplex #(
    parameter ORDER = MSB,
    parameter CS_ACTIVE = ACTIVE_LOW,
    parameter BYTE_INDEX_SIZE = 11 // roll over at 2048 bytes
)(
    input sck,
    input cs,
    input mosi,
    output logic [7:0] data,
    output logic data_avail,
    output logic [BYTE_INDEX_SIZE:0] byte_index
);

    logic first_flag;
    logic [2:0] bit_index;
    logic [7:0] buffer;

    always_ff @(posedge sck)
        if (cs == CS_ACTIVE) begin
            // else load bit into data
            if (ORDER == MSB) buffer <= { buffer[6:0], mosi };
            else buffer <= { mosi, buffer[7:1] };

            if (bit_index == 0) begin
                if (first_flag) begin
                    data_avail <= 1;
                    byte_index <= byte_index + 1;
                    data <= buffer;
                end else begin
                    first_flag <= 1;
                    byte_index <= ~0;
                end
            end else begin
                data_avail <= 0;
                byte_index <= byte_index;
                data <= data;
            end
            
            bit_index <= bit_index + 1;
            first_flag <= 1;
        end else begin
            // if cs disabled, reset
            buffer <= 0;
            data_avail <= 0;
            bit_index <= 0;
            byte_index <= byte_index;
            data <= data;
            first_flag <= 0;
        end
endmodule

