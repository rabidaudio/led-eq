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
    parameter BIT_DEPTH = 8,
    parameter CLOCK_DIVIDER = 6
) (
        input clk,
        input reset,

        // PixelBus interface
        // output logic req_read,
        // output logic [$clog2(BIT_DEPTH-1)-1:0] bitplane,
        // output logic [2:0] row_addr,
        // output logic [4:0] pixel_addr,
        // input read_ready,
        // input [1:0] red,
        // input [1:0] green,
        // input [1:0] blue,

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
        .CLOCK_DIVIDER(CLOCK_DIVIDER)
    ) matrix (
        .clk(clk),
        .reset(reset),

        // pixel bus
        // .req_read(req_read),
        // .bitplane(bitplane),
        // .row_addr(row_addr),
        // .pixel_addr(pixel_addr),
        // .read_ready(read_ready),
        // .red(red),
        // .green(green),
        // .blue(blue),

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
 * Shift register docs: https://www.micros.com.pl/mediaserver/UIMBI5020gp_0001.pdf
 * Data is shifted on clock rising edge. Data is latched when latch goes low.
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
    parameter BIT_DEPTH = 3,
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

    // PixelBus interface
    // output logic req_read,
    // output logic [$clog2(BIT_DEPTH-1)-1:0] bitplane,
    // output logic [$clog2(SCAN_RATE-1)-1:0] row_addr,
    // output logic [$clog2(WIDTH-1)-1:0] pixel_addr,
    // input read_ready,
    // input [COLOR_WIDTH-1:0] red,
    // input [COLOR_WIDTH-1:0] green,
    // input [COLOR_WIDTH-1:0] blue,

    // number of cycles of DWELL_CYCLES the LEDs will actually be on, from 0 (0%)
    // to DWELL_CYCLES (100%), inclusive
    input [$clog2(DWELL_CYCLES)-1:0] brightness,

    // Raised for one system clock cycle when the last latch for a frame occurs
    output logic frame_complete,

    output logic [COLOR_WIDTH-1:0] red,
    output logic [COLOR_WIDTH-1:0] green,
    output logic [COLOR_WIDTH-1:0] blue,
    output logic [$clog2(SCAN_RATE-1)-1:0] row_select,
    output logic out_clk,
    output logic lat,
    output logic oe_n
);
    // how many (LED clock) cycles it takes to latch out the shifted data.
    // LEDs must be off (oe_n high) while latch is taking place.
    localparam LATCH_CYCLES = 2;
    localparam PIXELS_PER_ROW = WIDTH;

    // Shift state

    // While a shift is happening, we request data from the PixelBus (SHIFTING).
    // When `read_ready`, the data is shifted out and the next pixel is requested.
    // when all the pixels are shifted out, `SHIFT_COMPLETE` goes high.
    enum logic { SHIFTING, SHIFT_COMPLETE } shift_state;
    
    // which pixel are we currently shifting out
    logic [$clog2(PIXELS_PER_ROW-1)-1:0] pixel_index;
    // which row(s) are we currently shifting out
    logic [$clog2(SCAN_RATE-1)-1:0] row_index;
    // which bitplane are we currently on
    logic [$clog2(BIT_DEPTH)-1:0] bitplane;

    // LEDs go on for the current `row_select` for the length of
    // `DWELL_TIME*2^bitplane`. When this time is finished, LEDs
    // are turned off and, if not yet complete, we wait until
    // `SHIFT_COMPLETE`. Then the data is latched. After the latch
    // completes, the DWELL period for the just shifted data begins
    // and we start shifting out the new data.
    enum logic { DWELL, LATCH } dwell_state;

    localparam MAX_DWELL_CYCLES = (1 << (BIT_DEPTH-1)) * DWELL_CYCLES;
    localparam COUNTER_SIZE = MAX_DWELL_CYCLES + LATCH_CYCLES;
    // how much dwell time is left
    logic [$clog2(COUNTER_SIZE)-1:0] dwell_counter;

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
    logic trigger_shift;
    always_comb out_clk = (low_clk & trigger_shift);
    
    always_comb oe_n = !(enable /*& b_pwm*/); // TODO: brightness
 
    always_ff @(posedge clk) begin
        frame_complete <= 0;

        if (about_to_fall) begin // trigger logic on falling edge of led_clk

            // shift out data
            if (shift_state == SHIFTING || dwell_state == LATCH) begin
                // shift out pixels
                red[0] <= (pixel_index == 0 && row_index == 0);
                blue[1] <= (bitplane <= pixel_index);

                pixel_index <= pixel_index + 1;
                trigger_shift <= 1;
            end

            // shift
            if (shift_state == SHIFTING) begin
                if (pixel_index == PIXELS_PER_ROW-1) begin
                    // row complete
                    pixel_index <= 0;
                    shift_state <= SHIFT_COMPLETE;
                    trigger_shift <= 1;
                    
                    if (row_index == SCAN_RATE-1) begin
                        row_index <= 0;
                        bitplane <= (bitplane == BIT_DEPTH-1) ? 0 : bitplane + 1;
                    end else row_index <= row_index + 1;
                end else pixel_index <= pixel_index + 1;
            end else trigger_shift <= 0;

            case (dwell_state)
                DWELL: begin
                    lat <= 0;
                    if (dwell_counter == 0) begin
                        enable <= 0;
                        if (shift_state == SHIFT_COMPLETE) begin
                            dwell_state <= LATCH;
                            dwell_counter <= LATCH_CYCLES-1;
                            lat <= 1;
                        end // else nothing to do, keep waiting
                    end else begin
                        enable <= 1;
                        dwell_counter <= dwell_counter - 1;
                    end
                end
                LATCH: begin
                    enable <= 0;
                    lat <= 0;
                    row_select <= row_select == SCAN_RATE-1 ? 0 : row_select + 1;
                    // restart
                    dwell_state <= DWELL;
                    dwell_counter <= (DWELL_CYCLES-1) << bitplane;
                    shift_state <= SHIFTING;
                    trigger_shift <= 1;

                    if (row_index == 0 && bitplane == 0) frame_complete <= 1;
                end
            endcase
        end

        if (reset) begin
            red <= 0;
            green <= 0;
            blue <= 0;
            row_select <= SCAN_RATE-2;
            lat <= 0;

            pixel_index <= 0;
            row_index <= 0;
            bitplane <= 0;
            shift_state <= SHIFT_COMPLETE;
            dwell_state <= DWELL;
            dwell_counter <= 0;
            enable <= 0;
            frame_complete <= 0;
            trigger_shift <= 0;
        end
    end

endmodule

`endif // __LEDDMATRIX__
