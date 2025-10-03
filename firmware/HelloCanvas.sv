`ifndef __HELLO_CANVAS__
`define __HELLO_CANVAS__

/**
 * This is a simple implementation of the `PixelBus` interface
 * that the `LEDMatrix` will use to read data from.
 * It loads "hello.b.mem" (a 1-bit 16x32 image) into memory
 * and writes it to the red channel on request.
 */
module HelloCanvas #(
    parameter SCAN_RATE = 8
) (
    input clk,
    input reset,

    // PixelBus
    input req_read,
    // NOTE: we only support 1 bit plane, parameter is ignored
    // input [4:0] bitplane_addr,
    input [2:0] row_addr,
    input [4:0] pixel_addr,
    output logic [1:0] r_data,
    output logic [1:0] g_data,
    output logic [1:0] b_data
);
    localparam HEIGHT = 16;
    localparam WIDTH = 32;
    localparam COLOR_WIDTH = HEIGHT/SCAN_RATE;

    (* ram_style = "block" *) logic display [HEIGHT] [WIDTH];
    initial $readmemb("hello.b.mem", display);

    always_ff @(posedge clk) begin
        if (req_read) begin
            for (int i = 0; i < COLOR_WIDTH; i++) begin
                r_data[i] <= display[row_addr + (i * SCAN_RATE)][pixel_addr];
            end
        end

        if (reset) begin
            r_data <= 0;
            g_data <= 0;
            b_data <= 0;
        end
    end
endmodule

`endif // __HELLO_CANVAS__
