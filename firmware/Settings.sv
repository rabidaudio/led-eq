`ifndef __SETTINGS__
`define __SETTINGS__

module Settings (
    /*
    ADDR    FIELD           DIR     DESC
    0x00    version[7:0]    ro      A constant indicating the firmware version. [1]
    0x01    width[7:0]      rw      How many pixels wide the display is. Used for
                                    determining the number of pixels expected per
                                    frame over SPI.
    0x02    height[7:0]     rw      How many pixels tall the display is. Used for
                                    determining the number of pixels expected per
                                    frame over SPI.
    0x03    c_depth[2:0]    rw      The color bit depth.
                                    0b000: 1 bit/color (3 bit color)
                                    0b001: 2 bits/color (6 bit color)
                                    ...
                                    0b111: 8 bits/color (24 bit color) [default]
            frame_int_en[3] rw      Frame complete gpio interrupt enabled
                                    [default=0 (disabled)]
            s_depth[7:4]    rw      The bit depth of the shift register (i.e.
                                    the number of shifts to complete one row).
                                    The number of shifts is (the value-1)*16:
                                    0b0000: 16-bit shift register
                                    0b0001: 32-bit shift register [default]
                                    ...
                                    0b1111: 256-bit shift register
    0x04    a_width[1:0]    rw      How many pixels per color are shifted out
                                    simultaneously for Port A.
                                    0b00: port a disabled
                                    0b01: 1/color (r1,g1,b1)
                                    0b10: 2/color (r1,g1,b1,r2,g2,b2) [default]
                                    0b11: 3/color NOTE: not standard HUB75
            a_rs[3:2]       rw      How many row select pins are used on Port A.
                                    Additional pins will be grounded.
                                    0b00: 2 pins (A, B)
                                    0b01: 3 pins (A, B, C)
                                    0b10: 4 pins (A, B, C, D)
                                    0b11: 5 pins NOTE: not standard HUB75
            b_width[5:4]    rw      How many pixels per color are shifted out
                                    simultaneously for Port B. [default=0b00 disabled]
            b_rs[7:6]       rw      How many row select pins are used on Port B.
                                    Additional pins will be grounded.
    0x05    dwell[7:0]      rw      The amount of time spent displaying a single
                                    pixel for the lowest bit depth color. Higher
                                    values have higher max brightness and more
                                    brightness resolution, but lower frame rates.
                                    0x00: 1 6MHz clock cycle (166.66ns)
                                    0x01: 2 6MHz clock cycles (333.33ns)
                                    ...
                                    0x0F: 16 clock cycles (2.66us) [default]
                                    0xFF: 256 clock cycles (42.66us)
    0x06    frame_dup[7:0]  rw      The number of times to repeat a frame.
                                    This reduces the effective framerate in
                                    exchange for visual quality.
                                    [default=0 (no repeats)]
    0x07    brightness[7:0] rw      The overall brightness of the display.
                                    0x00: display off
                                    0xFF: full brightness
                                    Actual resolution is limited by `dwell`,
                                    from 0 to `dwell`, inclusive. The provided
                                    value will be scaled automatically.
    0x08    cc_red[7:0]     rw      Color correction value for the red channel.
                                    Signed 8-bit integer. TODO: should this be
                                    added/subtracted? A fixed point scalar value?
    0x09    cc_green[7:0]   rw      Color correction value for the green channel.
    0x0A    cc_blue[7:0]    rw      Color correction value for the blue channel.
    0x0B    fps[7:0]        ro      Approximate current effective frame rate/4.
                                    Measured after applying `frame_dup` (e.g.
                                    300 refreshes/second with a `frame_dup` of 3
                                    will report 100 fps (0x19)).
                                    0x00: 4 frames/second or less
                                    0x01: 8 frames/second
                                    ...
                                    0xFF: 1024 frames/second or beyond
    0x0C    f_flag[7]       ro      Frame complete flag. Set to `1` when the most
                                    resent frame has completed rendering. Reset
                                    to `0` when any frame data is written.
            frame_cnt[6:0]  ro      Current frame counter. Rolls over every 128
                                    frames.
    0x0D    [reserved]
    0x0E    [reserved]
    0x0F    [reserved]
    */
    input [7:0] mem [16],
);

endmodule

`endif // __SETTINGS__
