`include "ResetGenerator.sv"
`include "LEDMatrix.sv"
`include "BrightnessRamp.sv"
`include "FlagCanvas.sv"

module top #(
    parameter RESET_AFTER = 'hFFFF,
    parameter BIT_DEPTH = 8
) (
    input clk,
    output logic [15:0] hub_75_o
);
    logic reset;

    ResetGenerator #(.AFTER(RESET_AFTER)) reset_gen (.clk(clk), .reset(reset));

    logic req_read;
    logic [$clog2(BIT_DEPTH-1)-1:0] bitplane_addr;
    logic [2:0] row_addr;
    logic [4:0] pixel_addr;
    logic read_ready;
    logic [1:0] r_data;
    logic [1:0] g_data;
    logic [1:0] b_data;
    
    FlagCanvas cv (
        .clk(clk),
        .reset(reset),
        
        .req_read(req_read),
        .bitplane_addr(bitplane_addr),
        .row_addr(row_addr),
        .pixel_addr(pixel_addr),
        .r_data(r_data),
        .g_data(g_data),
        .b_data(b_data)
    );

    logic frame_end;
    logic [4:0] brightness;

    BrightnessRamp ramp (
        .clk(clk),
        .reset(reset),
        .inc(frame_end),
        .brightness(brightness)
    );

    LEDMatrix_32x16_1to8 #(.BIT_DEPTH(8)) matrix (
        .clk(clk),
        .reset(reset),

        // pixel bus
        .req_read(req_read),
        .bitplane_addr(bitplane_addr),
        .row_addr(row_addr),
        .pixel_addr(pixel_addr),
        .r_data(r_data),
        .g_data(g_data),
        .b_data(b_data),

        .frame_complete(frame_end),
        // .brightness(brightness),
        .brightness(4'h03),

        .hub_75(hub_75_o)
    );
endmodule
