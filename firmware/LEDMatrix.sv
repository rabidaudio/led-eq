`ifndef __LEDDMATRIX__
`define __LEDDMATRIX__

module LEDMatrix #(
    parameter WIDTH = 32,
    parameter HEIGHT = 16,
    parameter BITS_PER_COLOR = 2,
    parameter SCAN_RATE = 8,
    parameter CLOCK_DIV_WIDTH = 21
) (
    input clk,
    input reset,

    output logic [BITS_PER_COLOR-1:0] red,
    output logic [BITS_PER_COLOR-1:0] green,
    output logic [BITS_PER_COLOR-1:0] blue,
    output logic [$clog2(SCAN_RATE-1)-1:0] row_select,
    output logic out_clk,
    output logic lat,
    output logic oe
);
    localparam PIXELS_PER_SCAN = (WIDTH*HEIGHT)/SCAN_RATE;

    logic [CLOCK_DIV_WIDTH-1:0] clock_div;
    logic [16:0] pixel_index;

    always_comb out_clk = clk && clock_div == 0;
    // always_ff @(posedge clk) out_clk <= clk && clock_div == 0;

    always_ff @(posedge clk) begin
        clock_div <= clock_div - 1;
        oe <= ~oe;

        // if (clock_div == 0) out_clk <= ~out_clk;

        if (out_clk) begin
            if (pixel_index == 0) begin
                row_select <= row_select + 1;
                lat <= 1;
                pixel_index <= PIXELS_PER_SCAN;
            end else begin
                lat <= 0;
                pixel_index <= pixel_index - 1;
            end
        end

        if (reset) begin
            red <= 'b01;
            green <= 'b10;
            blue <= 0;
            pixel_index <= PIXELS_PER_SCAN;
            row_select <= 0;

            oe <= 1; // turn on TODO pwm
            out_clk <= 0;

            clock_div <= 0;
        end
    end

endmodule

`endif // __LEDDMATRIX__
