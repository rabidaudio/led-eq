`ifndef __BRIGHTNESS_RAMP__
`define __BRIGHTNESS_RAMP__

module BrightnessRamp #(
    parameter COUNTER_WIDTH = 5,
    parameter BRIGHTNESS_WIDTH = 5
) (
    input clk,
    input reset,

    input inc,
    output logic [BRIGHTNESS_WIDTH-1:0] brightness
);
    localparam MAX_BRIGHTNESS = 1 << (BRIGHTNESS_WIDTH-1);
    logic [COUNTER_WIDTH-1:0] counter;

    always_ff @(posedge clk) begin
        if (inc) begin
            counter <= counter + 1;
            if (counter == 0) brightness <= brightness == MAX_BRIGHTNESS ? 0 : brightness + 1;
        end

        if (reset) begin
            brightness <= 0;
            counter <= 0;
        end
    end
endmodule

`endif // __BRIGHTNESS_RAMP__
