// A set of test bench utility modules
`ifndef _TestUtils_
`define _TestUtils_

/**
* GoldenMonitor is a test utility for comparing
* a generated signal to a golden signal, with
* an optional delay (in clock cycles).
*/
module GoldenMonitor #(
    // the bit width of both inputs
    parameter WIDTH=1,
    // how many clock cycles to delay golden to match signal
    parameter DELAY=0
) (
    input clk,
    input enable,
    input [(WIDTH-1):0] golden,
    input [(WIDTH-1):0] signal
);

    logic [(WIDTH-1):0] delay_line [$];

    initial begin
        for (int i = 0; i < DELAY; i++) begin
            @(posedge clk);
            delay_line.push_front(golden);
        end
        
        while (1) begin
            if (!DELAY) begin
                @(posedge clk);
                if (enable && signal != golden)
                    $error("monitor failed. expected %h but was %h", golden, signal);
            end else begin
                logic [(WIDTH-1):0] compare;
                
                @(posedge clk);
                delay_line.push_front(golden);
                compare = delay_line.pop_back();
                if (enable && signal != compare)
                    $error("monitor failed. expected %h but was %h", compare, signal);
            end
        end
    end
endmodule

`endif // _TestUtils_

