`ifndef __FLAG_CANVAS__
`define __FLAG_CANVAS__

module FlagCanvas #(
    parameter SCAN_RATE = 8
) (
    input clk,
    input reset,

    // PixelBus
    input req_read,
    input [2:0] bitplane_addr,
    input [2:0] row_addr,
    input [4:0] pixel_addr,
    output logic [1:0] r_data,
    output logic [1:0] g_data,
    output logic [1:0] b_data
);
    localparam HEIGHT = 16;
    localparam WIDTH = 32;
    localparam COLOR_WIDTH = HEIGHT/SCAN_RATE;

    // (* ram_style = "block" *)
    logic [23:0] display [HEIGHT] [WIDTH];
    initial $readmemh("build/flag.h.mem", display);

    always_ff @(posedge clk) begin
        if (req_read) begin
            for (int i = 0; i < COLOR_WIDTH; i++) begin
                r_data[i] = display[row_addr + (i*SCAN_RATE)][pixel_addr][bitplane_addr + 16];
                g_data[i] = display[row_addr + (i*SCAN_RATE)][pixel_addr][bitplane_addr + 8];
                b_data[i] = display[row_addr + (i*SCAN_RATE)][pixel_addr][bitplane_addr + 0];
            end
        end

        if (reset) begin
            r_data <= 0;
            g_data <= 0;
            b_data <= 0;
        end
    end
endmodule

`endif // __FLAG_CANVAS__
