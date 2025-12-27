`ifndef __ResetGenerator__
`define __ResetGenerator__

/**
 * This utility module will emit a reset signal shortly after power-up.
 * It takes advantage of the iCE40's initial / reset capabilities on
 * power-up.
 */
module ResetGenerator #(
    parameter AFTER = 'hFFFF
) (
    input clk,
    output logic reset
);
    logic [15:0] reset_counter;
    initial reset_counter = AFTER;
    initial reset = 1;

    always_ff @(posedge clk) begin
        if (reset_counter == 0) begin
            reset <= 0;
            reset_counter <= 0;
        end else reset_counter <= reset_counter - 1;
    end
endmodule

`endif // __ResetGenerator__
