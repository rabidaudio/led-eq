`ifndef __UART__
`define __UART__

/**
 * Whenever a data byte is presented on "txdata_i" and "txdata_valid_i" is strobed high for a
 * single clock cycle, the byte on "txdata_i" is latched and then shifted out on "uart_tx_o".
 *
 * While a byte is busy being shifted out, "uart_busy_o" is held high.
 */
module uart_tx #(
    parameter CLK_DIVIDER = 104
)  (
    input clk_i,
    input reset_i,

    input txdata_valid_i,
    input [7:0] txdata_i,

    output logic uart_busy_o,
    output logic uart_tx_o
);
    logic [15:0] clk_count;
    logic [9:0] uart_tx_shreg;
    logic [3:0] bitidx;

    always_ff @(posedge clk_i) begin
        if (uart_busy_o) begin
            clk_count <= clk_count + 1;
            if (clk_count == (CLK_DIVIDER - 1)) begin
                clk_count <= 0;

                // shift the next bit out
                if (bitidx == 10) uart_tx_o <= 1;
                else uart_tx_o <= uart_tx_shreg[bitidx];

                // Decide what state to go to next.
                // Note that we don't leave until we've finished
                // fully transmitting the last bit.
                bitidx <= bitidx + 1;
                if (bitidx == 10) uart_busy_o <= '0;
            end
        end else begin
            uart_tx_o <= '1;

            if (txdata_valid_i) begin
                uart_tx_shreg <= {1'b1, txdata_i, 1'b0};
                clk_count <= (CLK_DIVIDER - 1);
                bitidx <= 0;
                uart_busy_o <= '1;
            end
        end

        if (reset_i) begin
            clk_count <= 0;
            uart_busy_o <= '0;
            uart_tx_shreg <= 'x;
        end
    end
endmodule


/**
 * This module listens for new bytes to arrive on "uart_rx_i". Whenever a new byte has been
 * received, it is presented on "rxdata_o" and "rxdata_valid_o" is strobed high for a single cycle.
 *
 * When "rxdata_valid_o" is strobed, the downstream hardware must accept the presented byte.
 */
module uart_rx #(
    // a clock divider of 104 gives us a 115,200 baud serial connection given a 12MHz clk_i.
    parameter CLK_DIVIDER = 104
)  (
    input clk_i,
    input reset_i,

    input uart_rx_i,

    output logic rxdata_valid_o,
    output logic [7:0] rxdata_o,
    output logic uart_busy_o
);
    // metastability resolution
    logic [1:0] _uart_rx;
    always_ff @(posedge clk_i) _uart_rx = {_uart_rx[0], uart_rx_i};

    logic [15:0] clk_count;
    logic [9:0] uart_rx_shreg;
    logic [3:0] bitidx;

    always_ff @(posedge clk_i) begin
        rxdata_valid_o <= 0;
        if (uart_busy_o) begin
            // we sample the signal halfway through the sampling period
            if (clk_count == (CLK_DIVIDER >> 1)) begin
                uart_rx_shreg[bitidx] <= _uart_rx[1];

                // we bail out early if we've sampled the final bit.
                if (bitidx == 9) begin
                    uart_busy_o <= 0;
                    rxdata_o <= uart_rx_shreg[8:1];
                    rxdata_valid_o <= 1;
                end
            end

            // If we're done with our sampling period, advance.
            clk_count <= clk_count + 1;
            if (clk_count == (CLK_DIVIDER - 1)) begin
                clk_count <= 0;
                bitidx <= bitidx + 1;
            end
        end else begin
            // if we're idle, as soon as we see a low signal, we start shifting in.
            if (!_uart_rx[1]) begin
                uart_busy_o <= 1;
                clk_count <= 0;
                bitidx <= 0;
            end
        end

        if (reset_i) begin
            rxdata_valid_o <= 0;
            uart_busy_o <= 0;
            clk_count <= 0;
            bitidx <= 0;
        end
    end
endmodule

`endif
