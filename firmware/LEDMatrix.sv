`ifndef __LEDDMATRIX__
`define __LEDDMATRIX__

`include "ClockDivider.sv"
`include "PixelBus.sv"

`define MAX(a, b) (((``a) > (``b)) ? (``a) : (``b))

/**
 * 32x16 HUB75 display with 2 bits for each color,
 * a 3 bit row selector, and 1:8 scan rate.
 * The display is broken in to upper and lower halves.
 * Color bits r0,g0,b0 control the upper half and
 * r1,g1,b1 control the lower half. Rows 0+i and 8+i
 * are controlled by row_select, for i in [0,8).
 * https://news.sparkfun.com/2650
 */
module LEDMatrix_32x16_1to8 #(
    parameter BIT_DEPTH = 8
) (
        input clk,
        input reset,
        // input [7:0] brightness,
        output logic [15:0] hub_75
    );
    LEDMatrix #(
        .WIDTH(32),
        .HEIGHT(16),
        .SCAN_RATE(8),
        .COLOR_WIDTH(2),
        .BIT_DEPTH(BIT_DEPTH),
        .DWELL_CYCLES(16),
        .CLOCK_DIVIDER(6)
    ) matrix (
        .clk(clk),
        .reset(reset),

        // .brightness(brightness),

        .red({ hub_75[4], hub_75[0] }),
        .green({ hub_75[5], hub_75[1] }),
        .blue({ hub_75[6], hub_75[2] }),
        .row_select(hub_75[10:8]),
        .out_clk(hub_75[12]),
        .lat(hub_75[13]),
        .oe_n(hub_75[14])
    );

    // assign instead of always_comb to suppress warning
    assign hub_75[3] = 0;
    assign hub_75[7] = 0;
    assign hub_75[11] = 0;
    assign hub_75[15] = 0;
endmodule

/**
 * LEDMatrix drives HUB75-style led matrix displays. These displays
 * update multiple scan lines in parallel using a shift register.
 */
module LEDMatrix #(
    parameter WIDTH = 32, // in pixels
    parameter HEIGHT = 32, // in pixels
    // the ratio of pixels on at a given time. Alternatively,
    // how many shifts does it take to cycle through a whole frame
    parameter SCAN_RATE = 16,
    // how many rows are shifted out in parallel
    parameter COLOR_WIDTH = (HEIGHT/SCAN_RATE),
    // bit depth of our input colors, ie. how many bitplanes
    parameter BIT_DEPTH = 16,
    // The amount of time for the LEDs be on for the LSB bitplane.
    // This parameter controls the max perceived brightness, the max framerate,
    // and the resolution of `brightness` in some complex ways. See:
    // https://docs.google.com/spreadsheets/d/1fURkK-2R26DpsshLflX2rYrjb4WjTQQUoPvTYIlk8Hc/edit?gid=1681767793#gid=1681767793
    parameter DWELL_CYCLES = 32,
    // A divider from the system clock to run the display with
    parameter CLOCK_DIVIDER = 2
) (
    input clk,
    input reset,
    // number of cycles of DWELL_CYCLES the LEDs will actually be on, from 0 (0%)
    // to DWELL_CYCLES (100%), inclusive
    input [$clog2(DWELL_CYCLES)-1:0] brightness,

    output logic [COLOR_WIDTH-1:0] red,
    output logic [COLOR_WIDTH-1:0] green,
    output logic [COLOR_WIDTH-1:0] blue,
    output logic [$clog2(SCAN_RATE-1)-1:0] row_select,
    output logic out_clk,
    output logic lat,
    output logic oe_n
);
    // how many (LED clock) cycles it takes to shift out one row
    localparam SHIFT_CYCLES = (WIDTH*HEIGHT)/SCAN_RATE/COLOR_WIDTH;
    // how many (LED clock) cycles it takes to latch out the shifted data.
    // LEDs must be off (oe_n high) while latch is taking place.
    localparam LATCH_CYCLES = 2;

    // inputs: fixed dwell time
    // where "dwell" = period of time on (*2^bit depth)
    // period max(shift+latch time, dwell + latch time)
    // first iteration:
    //      shift, latch, dwell
    // second iteration:
    //      dwell, latch (with shift just in time)
    enum logic [2:0] {
        SHIFT = 1, // LEDs off, shift out the next row of data
        DWELL = 2, // LEDs on, wait for DWELL_CYCLES*2^n where n is the current bit plane
        LATCH = 3  // LEDs off, latch the row and increment the row counter
    } state;
    localparam MAX_DWELL_CYCLES = (1 << (BIT_DEPTH-1)) * DWELL_CYCLES;
    localparam COUNTER_SIZE = `MAX(SHIFT_CYCLES, `MAX(LATCH_CYCLES, MAX_DWELL_CYCLES));
    // core counter until the next state transition
    logic [$clog2(COUNTER_SIZE)-1:0] counter;

    // which pixel are we currently shifting out
    logic [$clog2(SHIFT_CYCLES-1)-1:0] pixel_index;
    // which row(s) are we currently shifting out
    logic [$clog2(SCAN_RATE-1)-1:0] row_index;

    logic [$clog2(BIT_DEPTH-1)-1:0] bitplane;

    // whether  the LEDs be on
    logic enable;
    // the system clock divided by CLOCK_DIVIDER
    logic low_clk;
    // pulled high one system clock before low_clk will fall.
    // this is where state changes should be made so that they
    // will be picked up on the next rising edge of low_clk
    logic about_to_fall;

    ClockDivider #(.DIVIDER(CLOCK_DIVIDER)) clk_div (
        .clk(clk),
        .reset(reset),
        .slow_clk(low_clk),
        .next_fall(about_to_fall)
    );

    // out_clk should only tick while shifting
    always_comb out_clk = (low_clk & state == SHIFT);
    
    always_comb oe_n = !(enable /*& b_pwm*/); // TODO: brightness
 
    always_ff @(posedge clk) begin
        if (about_to_fall) begin // trigger logic on falling edge of led_clk

            case (state)
                SHIFT: begin
                    enable <= 0;
                end
                DWELL: enable <= 1;
                LATCH: begin
                    enable <= 0;
                    lat <= ~lat;
                    if (counter == 1) begin
                        row_select <= row_select == SCAN_RATE-1 ? 0 : row_select + 1;
                    end
                end
            endcase

            // if state is SHIFT or is about to be SHIFT (end of LATCH)
            if (state == SHIFT || (state == LATCH && counter == 0)) begin
                // shift out pixels
                red[0] <= (row_index == 0 && pixel_index == 0);
                blue[1] <= (bitplane <= pixel_index);

                pixel_index <= pixel_index + 1;
            end

            // state transition
            counter <= counter - 1;
            if (counter == 0) begin
                if (state == SHIFT) begin
                    state <= DWELL;
                    counter <= (DWELL_CYCLES * (1 << bitplane))-1;

                    // primitively prepare to start shifting next
                    if (row_index == SCAN_RATE-1) begin
                        row_index <= 0;

                        // bitplane complete, start shifting out next bitplane
                        if (bitplane == BIT_DEPTH-1) begin
                            // end of frame. TODO: signal frame complete?
                            bitplane <= 0;
                        end else bitplane <= bitplane + 1;
                    end else row_index <= row_index + 1; // increment row

                    pixel_index <= 0; // reset pixel index
                    
                end else if (state == DWELL) begin
                    state <= LATCH;
                    counter <= LATCH_CYCLES-1;
                end else if (state == LATCH) begin
                    state <= SHIFT;
                    counter <= SHIFT_CYCLES-1;
                end
            end
        end

        if (reset) begin
            red <= 0;
            green <= 0;
            blue <= 0;
            row_select <= SCAN_RATE-2;
            lat <= 0;

            pixel_index <= 0;
            row_index <= 0;
            state <= DWELL;
            counter <= 0;
            bitplane <= 0;
            enable <= 0;
        end
    end

endmodule

`endif // __LEDDMATRIX__
