`ifndef __CLOCK_DIVIDER__
`define __CLOCK_DIVIDER__

module ClockDivider #(
    parameter DIVIDER = (1 << 4)
) (
    input clk,
    input reset,
    output logic slow_clk,
    output logic next_rise,
    output logic next_fall
);
    localparam WIDTH = $clog2(DIVIDER);
    logic [WIDTH-1:0] counter;

    always_ff @(posedge clk) begin
        if (counter >= (DIVIDER/2)) begin
            slow_clk <= 1;
            counter <= counter - 1;
        end else begin
            if (counter == 0) counter <= DIVIDER-1;
            else counter <= counter - 1;
            slow_clk <= 0;
        end

        next_rise <= counter == 0;
        next_fall <= counter == DIVIDER/2;

        if (reset) begin
            slow_clk <= 0;
            counter <= 0;
            next_rise <= 0;
            next_fall <= 0;
        end
    end
endmodule

`endif // __CLOCK_DIVIDER__
