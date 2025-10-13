`ifndef __SpiSlaveSimplex__
`define __SpiSlaveSimplex__

typedef enum { MSB = 0, LSB = 1 } bit_order_e;
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
    parameter BYTE_COUNTER_MAX = 2048
)(
    input sck,
    input reset,
    input cs,
    input mosi,
    output logic [7:0] data,
    output logic data_avail,
    output logic [$clog2(BYTE_COUNTER_MAX-1)-1:0] byte_index
);
    logic [2:0] bit_index;
    logic [7:0] buffer;

    always_ff @(posedge sck) begin
        logic [7:0] new_data;
        
        new_data = { buffer[6:0], mosi };

        if (cs == CS_ACTIVE) begin
            if (bit_index == 0) begin
                byte_index <= byte_index + 1;
                data <= new_data;
            end
            
            data_avail <= bit_index == 0;
            buffer <= new_data;
            bit_index <= bit_index - 1;
        end else begin
            data_avail <= 0;
            byte_index <= 0;
        end

        if (reset) begin
            data <= 0;
            data_avail <= 0;
            byte_index <= '1;
            buffer <= 0;
            bit_index <= '1;
        end 
    end
endmodule

`endif // __SpiSlaveSimplex__
