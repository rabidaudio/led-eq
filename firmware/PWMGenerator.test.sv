`include "PWMGenerator.sv"
`include "utils.test.sv"

module PWMGenerator_Tests (
    input t_clk,
    input t_reset
);
    PWMGenerator_Test0 test_pwm_0(t_clk, t_reset);
    PWMGenerator_Test100 test_pwm_100(t_clk, t_reset);
    PWMGenerator_Test50 test_pwm_50(t_clk, t_reset);
    PWMGenerator_TestUpdateDuty test_pwm_update(t_clk, t_reset);
endmodule

module PWMGenerator_Test0 (
    input t_clk,
    input t_reset
);
    localparam PERIOD = 4;

    logic t_pwm;

    PWMGenerator #(
        .INITIAL_PERIOD(PERIOD),
        .INITIAL_DUTY(0)
    ) dut (
        .clk(t_clk),
        .reset(t_reset),
        .update_parameters('0),
        .pwm(t_pwm)
    );

    // assert that t_pwm is low all the time
    logic gm_pwm_enable;
    GoldenMonitor gm_pwm(
        .clk(t_clk),
        .enable(gm_pwm_enable),
        .golden('0),
        .signal(t_pwm)
    );
    initial begin
        gm_pwm_enable <= 0;
        repeat(3) @(posedge t_clk);
        gm_pwm_enable <= 1;
    end
endmodule

module PWMGenerator_Test100 (
    input t_clk,
    input t_reset
);
    localparam PERIOD = 4;

    logic t_pwm;

    PWMGenerator #(
        .INITIAL_PERIOD(PERIOD),
        .INITIAL_DUTY(PERIOD)
    ) dut (
        .clk(t_clk),
        .reset(t_reset),
        .update_parameters('0),
        .pwm(t_pwm)
    );

    // assert that t_pwm is high all the time
    logic gm_pwm_enable;
    GoldenMonitor gm_pwm(
        .clk(t_clk),
        .enable(gm_pwm_enable),
        .golden('1),
        .signal(t_pwm)
    );
    initial begin
        gm_pwm_enable <= 0;
        repeat(3) @(posedge t_clk);
        gm_pwm_enable <= 1;
    end
endmodule

module PWMGenerator_Test50 (
    input t_clk,
    input t_reset
);
    localparam PERIOD = 8;
    localparam logic [3:0] DUTY = 4;

    logic t_pwm;
    logic t_period_end;

    PWMGenerator #(
        .INITIAL_PERIOD(PERIOD),
        .INITIAL_DUTY(DUTY)
    ) dut (
        .clk(t_clk),
        .reset(t_reset),
        .update_parameters('0),
        .period_end(t_period_end),
        .pwm(t_pwm)
    );

    // assert duty correct
    logic [3:0] t_duty;
    DutyCounter #(.WINDOW(PERIOD), .ASSERT(1)) duty_counter (
        .t_clk(t_clk), .signal(t_pwm), .g_duty(DUTY), .duty(t_duty)
    );

    // assert that period_end goes high once at the end of each period
    logic g_period_end;
    ConstantPWM #(.PERIOD(PERIOD), .DUTY(1)) start_reference (
        .t_clk(t_clk), .reset(t_reset), .pwm(g_period_end)
    );
    GoldenMonitor #(.DELAY(PERIOD+1)) gm_p_end(
        .clk(t_clk),
        .enable('1),
        .golden(g_period_end),
        .signal(t_period_end)
    );
endmodule

module PWMGenerator_TestUpdateDuty (
    input t_clk,
    input t_reset
);
    localparam PERIOD = 8;

    logic t_pwm;
    logic [3:0] t_pwm_duty_cycle;
    logic t_update_parameters;
    logic t_period_end;

    PWMGenerator #(
        .INITIAL_PERIOD(PERIOD),
        .INITIAL_DUTY(0)
    ) dut (
        .clk(t_clk),
        .reset(t_reset),
        .pwm_duty_cycle(t_pwm_duty_cycle),
        .update_parameters(t_update_parameters),
        .period_end(t_period_end),
        .pwm(t_pwm)
    );

    logic [3:0] measured_duty;
    initial begin
        measured_duty <= 0;
        repeat(3) @(posedge t_clk);
        
        while (1) begin
            if (t_period_end) measured_duty <= 0;
            else measured_duty <= measured_duty + t_pwm;
            @(posedge t_clk);
        end
    end

    // when an end of period is reached, compare the measured period to
    // the one set the previous period
    GoldenMonitor #(.WIDTH(4), .DELAY(PERIOD)) gm_duty(
        .clk(t_clk),
        .enable(t_period_end),
        .golden(t_pwm_duty_cycle),
        .signal(measured_duty)
    );

    // change duty cycle at random points excluding the first period.
    // changes should take place the following period.
    initial begin
        int change_after;

        repeat(3) @(posedge t_clk);
        while (1) begin
            for (int duty = 0; duty < PERIOD; duty++) begin
                change_after = $urandom_range(PERIOD-2);
                for (int t = 0; t < PERIOD; t++) begin
                    t_pwm_duty_cycle <= duty;
                    t_update_parameters <= (t == change_after);
                    @(posedge t_clk);
                end
            end
        end
    end
endmodule
