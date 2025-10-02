`ifndef __MOCKCANVAS__
`define __MOCKCANVAS__

`include "PixelBus.sv"

/**
 * This is a simple implementation of the `PixelBus` interface
 * that the `LEDMatrix` will use to read data from.
 * It loads "hello.b.mem" (a 1-bit 16x32 image) into memory
 * and writes it to the red channel on request.
 */
module MockCanvas (
    input reset,

    /*
    PixelBus #(
        .WIDTH(32), 
        .HEIGHT(16),
        .SCAN_RATE(8),
        // it's actually 1 but we'll accept higher bit depths and always return 1
        .BIT_DEPTH(8)
    ) bus*/
    PixelBus bus
);
    logic display [bus.HEIGHT] [bus.WIDTH];
    initial $readmemb("hello.b.mem", display);

    always_ff @(posedge bus.clk) begin
        if (bus.req_read) begin
            for (int i = 0; i < bus.COLOR_WIDTH; i++) begin
                bus.red[i] <= display[bus.row_addr + (i * bus.SCAN_RATE)][bus.pixel_addr];
            end
            bus.read_ready <= 1;
        end else bus.read_ready <= 0;

        if (reset) begin
            bus.read_ready <= 0;
            bus.red <= 0;
            bus.green <= 0;
            bus.blue <= 0;
        end
    end
endmodule

`endif // __MOCKCANVAS__
