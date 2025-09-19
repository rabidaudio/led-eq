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

    logic prev_cs;
    // logic first_flag;
    logic [2:0] bit_index;
    logic [7:0] buffer;

    always_ff @(posedge sck) begin
        if (reset) begin
            data <= 0;
            data_avail <= 0;
            byte_index <= '1;
            buffer <= 0;
            bit_index <= '1;
            prev_cs <= cs;
        end else begin
            if (cs == CS_ACTIVE) begin
                if (bit_index == 0) begin
                    data_avail <= 1;
                    byte_index <=  byte_index + 1;                   
                end else begin
                    data_avail <= 0;
                    byte_index <= byte_index;
                end

                if (bit_index == 7) data <= buffer;
                else data <= data;

                // buffer[7-bit_index] <= mosi;
                buffer <= { buffer[6:0], mosi };
                bit_index <= bit_index - 1;
                // prev_cs <= cs;
            end else begin
                // TODO
            end
        end
    end

    // always_ff @(posedge sck) begin
    //     prev_cs <= cs;
    //     byte_index <= byte_index;
    //     data <= data;
    //     data_avail <= 0;

    //     if (prev_cs != CS_ACTIVE && cs == CS_ACTIVE) begin
    //         // start transaction
            
    //         buffer <= 0;
    //         first_flag <= 1;
            
    //         if (ORDER == MSB) buffer <= { buffer[6:0], mosi };
    //         else buffer <= { mosi, buffer[7:1] };
    //         bit_index <= 1;

    //     end else if (prev_cs == CS_ACTIVE && cs != CS_ACTIVE) begin
    //         // end transaction
                        

    //     end else if (cs == CS_ACTIVE) begin
    //         // continue transaction

    //         if (bit_index == 0) begin
    //             // send out byte
    //             if (first_flag) byte_index <= 0;
    //             else byte_index <= byte_index + 1;
    //             data_avail <= 1;
    //         end else begin
                
    //         end

    //         bit_index <= bit_index + 1;

    //     end else begin
    //         // idle
    //         buffer <= 0;
    //         bit_index <= 0;
    //     end

    //     if (cs == CS_ACTIVE) begin
    //         // else load bit into data
    //         if (ORDER == MSB) buffer <= { buffer[6:0], mosi };
    //         else buffer <= { mosi, buffer[7:1] };

    //         if (bit_index == 0) begin
    //             if (first_flag) begin
    //                 data_avail <= 1;
    //                 byte_index <= byte_index + 1;
    //                 data <= buffer;
    //             end else begin
    //                 first_flag <= 1;
    //                 byte_index <= '1;
    //             end
    //         end else begin
    //             data_avail <= 0;
    //             byte_index <= byte_index;
    //             data <= data;
    //         end
            
    //         bit_index <= bit_index + 1;
    //         first_flag <= 1;
    //     end else begin
    //         // if cs disabled, reset

    //         first_flag <= 0;
    //     end
    // end
endmodule

