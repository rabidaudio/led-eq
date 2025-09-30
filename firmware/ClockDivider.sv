`ifndef __CLOCK_DIVIDER__
`define __CLOCK_DIVIDER__

module ClockDivider #(
    parameter WIDTH = 2
) (
    input clk,
    input reset,
    output logic slow_clk,
    output logic next_rise,
    output logic next_fall
);

    logic [WIDTH-1:0] counter;

    always_ff @(posedge clk) begin
        if (counter[WIDTH-1]) begin
            slow_clk <= 1;
            counter <= counter - 1;
        end else begin
            counter <= counter - 1;
            slow_clk <= 0;
        end

        next_rise <= counter == 0;
        next_fall <= counter == (1 << (WIDTH - 1));

        if (reset) begin
            slow_clk <= 0;
            counter <= 0;
            next_rise <= 0;
            next_fall <= 0;
        end
    end
endmodule

`endif // __CLOCK_DIVIDER__
