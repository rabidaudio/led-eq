`ifndef __LEDDMATRIX__
`define __LEDDMATRIX__

/**
 * 32x16 HUB75 display with 2 bits for each color,
 * a 3 bit row selector, and 1:8 scan rate.
 * The display is broken in to upper and lower halves.
 * Color bits r0,g0,b0 control the upper half and
 * r1,g1,b1 control the lower half. Rows 0+i and 8+i
 * are controlled by row_select, for i in [0,8).
 * https://news.sparkfun.com/2650
 */
module LEDMatrix_32x16_1to8 (
        input clk,
        input reset,
        output [15:0] hub_75
    );
    LEDMatrix #(
        .WIDTH(16),
        .HEIGHT(32),
        .SCAN_RATE(8),
        .COLOR_WIDTH(2)
    ) matrix (
        .clk(clk),
        .reset(reset),

        .red({ hub_75[0], hub_75[4] }),
        .green({ hub_75[1], hub_75[5] }),
        .blue({ hub_75[2], hub_75[6] }),
        .row_select(hub_75[10:8]),
        // .out_clk(hub_75[12]),
        .lat(hub_75[13]),
        .oe_n(hub_75[14])
    );
endmodule

/**
 * LEDMatrix drives HUB75-style led matrix displays. These displays
 * update multiple scan lines in parallel using a shift register.
 */
module LEDMatrix #(
    // parameter conf = known_configs[m32x16_1to8]
    parameter WIDTH = 32,
    parameter HEIGHT = 32,
    parameter SCAN_RATE = 16,
    parameter COLOR_WIDTH = (HEIGHT/SCAN_RATE)
) (
    input clk,
    input reset,

    output logic [COLOR_WIDTH-1:0] red,
    output logic [COLOR_WIDTH-1:0] green,
    output logic [COLOR_WIDTH-1:0] blue,
    output logic [$clog2(SCAN_RATE-1)-1:0] row_select,
    // output logic out_clk,
    output logic lat,
    output logic oe_n
);
    localparam PIXELS_PER_SCAN = (WIDTH*HEIGHT)/SCAN_RATE/COLOR_WIDTH;

    // STOPSHIP
    // logic display [16] [32];
    // initial begin
    //     $readmemb("hello.b.mem", display);
    // end

    logic [16:0] pixel_index;

    always_ff @(posedge clk) begin
        if (pixel_index == 0) pixel_index <= PIXELS_PER_SCAN;
        else pixel_index <= pixel_index - 1;

        // red[0] <= display[row_select][pixel_index];
        // red[1] <= display[row_select+SCAN_RATE][pixel_index];

        lat <= pixel_index == 0;

        if (pixel_index == 0) row_select <= row_select + 1;

        if (reset) begin
            red <= 'b01;
            green <= 'b10;
            blue <= 0;
            pixel_index <= PIXELS_PER_SCAN;
            row_select <= 0;

            oe_n <= 0; // turn on TODO pwm
            // out_clk <= 0;
        end
    end

endmodule

`endif // __LEDDMATRIX__
