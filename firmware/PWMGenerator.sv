`ifndef _PWMGenerator_
`define _PWMGenerator_

/**
 * PWMGenerator generates a pwm output signal. The period and duty cycle of the pwm signal
 * is be determined by the input parameters `pwm_period` and `pwm_duty_cycle`, both measured
 * in clock cycles.
 */
module PWMGenerator #(
    parameter INITIAL_PERIOD=256,
    parameter INITIAL_DUTY=128,
    parameter WIDTH=$clog2(INITIAL_PERIOD)
)  (
    input clk,
    input reset,

    // when pulled high, `pwm_period` and `pwm_duty_cycle` inputs are updated, to be used
    // the following period
    // TODO: changing on the last cycle of the period causes it to take 2 full cycles to update
    input update_parameters,

    // clock cycles for one period
    input [WIDTH-1:0] pwm_period,

    // clock cycles for high value, permitted to be 0% to 100% inclusive, i.e. 0 to `pwm_period`+1
    input [WIDTH:0] pwm_duty_cycle,

    // pulled high on the last cycle of a period
    output logic period_end,

    // output signal
    // TODO: use PRS generator instead:
    // https://www.isotel.eu/mixedsim/intro/prssine.html
    output logic pwm
);
    // the period to use for the next cycle
    logic [WIDTH-1:0] next_period;

    // the duty cycle for the next cycle
    logic [WIDTH:0] next_duty_cycle;

    // counts down to zero to track period
    logic [WIDTH-1:0] period_counter;
    // counts down to zero to track high portion of period
    logic [WIDTH:0] high_counter;

    always_ff @(posedge clk) begin
        if (reset) begin
            next_period <= INITIAL_PERIOD;
            next_duty_cycle <= INITIAL_DUTY;
            period_counter <= INITIAL_PERIOD-1;
            high_counter <= INITIAL_DUTY;
            pwm <= 0;
        end else begin
            if (update_parameters) begin
                // TODO: allow changing period
                next_duty_cycle <= pwm_duty_cycle;
            end

            period_counter <= period_counter == 0 ? (next_period-1) : period_counter - 1;

            period_end <= (period_counter == 0);

            if (period_counter == 0)
                high_counter <= (update_parameters ? pwm_duty_cycle : next_duty_cycle);
            else if (high_counter > 0) high_counter <= high_counter - 1;
            else high_counter <= 0;

            pwm <= (high_counter > 0);
        end
    end
endmodule

`endif // _PWMGenerator_
