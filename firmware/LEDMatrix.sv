`ifndef __LEDDMATRIX__
`define __LEDDMATRIX__

`include "ClockDivider.sv"

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
        output logic [15:0] hub_75
    );
    LEDMatrix #(
        .WIDTH(32),
        .HEIGHT(16),
        .SCAN_RATE(8),
        .COLOR_WIDTH(2)
    ) matrix (
        .clk(clk),
        .reset(reset),

        .red({ hub_75[0], hub_75[4] }),
        .green({ hub_75[1], hub_75[5] }),
        .blue({ hub_75[2], hub_75[6] }),
        .row_select(hub_75[10:8]),
        .out_clk(hub_75[12]),
        .lat(hub_75[13]),
        .oe_n(hub_75[14])
    );

    always_comb begin : ground
        hub_75[3] = 0;
        hub_75[7] = 0;
        hub_75[11] = 0;
        hub_75[15] = 0;
    end
endmodule

`define MAX(a, b) (((``a) > (``b)) ? (``a) : (``b))

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
    output logic out_clk,
    output logic lat,
    output logic oe_n
);
    localparam SHIFT_CYCLES = (WIDTH*HEIGHT)/SCAN_RATE/COLOR_WIDTH;;
    localparam LATCH_CYCLES = 2;
    localparam DWELL_CYCLES = SHIFT_CYCLES-LATCH_CYCLES;

    // STOPSHIP
    logic display [HEIGHT] [WIDTH];
    initial begin
        $readmemb("hello.b.mem", display);
    end

    enum logic [2:0] { SHIFT = 1, LATCH = 2, DWELL = 3 } state;
    logic [$clog2(SHIFT_CYCLES-1)-1:0] counter;
    logic [$clog2(`MAX(SHIFT_CYCLES, `MAX(LATCH_CYCLES, DWELL_CYCLES)))-1:0] pixel_index;

    logic low_clk;
    logic about_to_rise;
    logic about_to_fall;

    ClockDivider #(.WIDTH(4)) clk_div (
        .clk(clk),
        .reset(reset),
        .slow_clk(low_clk),
        .next_rise(about_to_rise),
        .next_fall(about_to_fall)
    );

    always_comb out_clk = (low_clk & state == SHIFT);
 
    always_ff @(posedge clk) begin
        if (about_to_fall) begin // trigger logic on falling edge of led_clk
            case (state)
                SHIFT: begin
                    // shift out pixels
                    red[1] <= display[(row_select-1)][pixel_index];
                    red[0] <= display[(row_select-1)+SCAN_RATE][pixel_index];

                    pixel_index <= pixel_index + 1;
                end
                LATCH: begin
                    oe_n <= 1; // disable
                    lat <= ~lat;
                    if (counter == 1) row_select <= row_select + 1;
                end
                // DWELL: nothing to do
                DWELL: oe_n <= 0; // disable
            endcase

            // state transition
            counter <= counter - 1;
            if (counter == 0) begin
                if (state == SHIFT) begin
                    state <= LATCH;
                    counter <= LATCH_CYCLES-1;

                end else if (state == LATCH) begin
                    state <= DWELL;
                    counter <= DWELL_CYCLES-1;
                end else if (state == DWELL) begin
                    state <= SHIFT;
                    counter <= SHIFT_CYCLES-1;

                    pixel_index <= 0;
                end
            end
        end

        if (reset) begin
            red <= 0;
            green <= 0;
            blue <= 0;
            pixel_index <= 0;
            row_select <= '1;
            lat <= 0;
            state <= DWELL;
            counter <= 0;

            oe_n <= 1; // turn on TODO pwm
        end
    end

endmodule

`endif // __LEDDMATRIX__
