`ifndef __SPI__
`define __SPI__

module Spi (
    /*
        Settings registers:
        COMMAND [DATA...]
            command[7]      0b0         indicates accessing settings registers
            command[6]      RW bit      1 = write mode, 0 = read mode
            command[5:4]                [reserved]
            command[3:0]    addr        The address of the setting to
                                        begin reading/writing.

        Data registers:
        COMMAND ADDR1 ADDR2 [DATA...]
        command[7]          0b1         indicates accessing data registers
        command[6]          RW bit      1 = write mode, 0 = read mode
        command[5:4]                    [reserved]
        command[0]          addr[17]    MSB of 17 bit address to begin
        addr1[7:0]          addr[16:8]  Next most significant bits of address
        addr2[7:0]          addr[7:0]   Least significant bits of address

        Data addresses are 17 bits. Data is written per pixel in row-column order,
        with color data in RGB order. Color data is bit-packed subject to the `bit_pack`
        setting. Thus the length of one frame is dependent on the settings `width`,
        `height`, and `bit_pack`. 17 bit addresses set a maximum of the combination of
        these parameters, e.g. 256x256x2, 64x256x8, 128x256x4, etc.

        TODO: it seems like all the HUB75 panels use 2bit parallel colors and therefore max
        out at 32px high. In order to get 256 pixels high, we'd either need 8 ports, non-
        standard connectors, or a way to handle arbitrary vertical/horizontal tiling through
        chaining (increasing shift register depth). Unsure the best approach. For now perhaps
        limit to 64x256?
    */
);

endmodule

`endif // __SPI__
