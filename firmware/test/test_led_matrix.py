import numpy as np
from framework import simulator_test
from pathlib import Path
import math

from imageio.v3 import imwrite

@simulator_test(
    "LEDMatrix.sv",
    module_name="LEDMatrix_32x16_1to8",
    params={"BIT_DEPTH": 2, "CLOCK_DIVIDER": 2},
)
async def test_led_matrix(dut):
    dut.set(brightness=8)  # half brightness

    # fake PixelBus returns nothing for odd pixels.
    # for even pixels, returns red, green, blue, or white by row
    async def _pixel_generator():
        while True:
            if dut.req_read.value == 1:
                if dut.pixel_addr.value[0] == 1:
                    dut.set(r_data=0, g_data=0, b_data=0)
                else:
                    match int(dut.row_addr.value) % 4:
                        case 0:
                            dut.set(r_data=1, g_data=0, b_data=0)
                        case 1:
                            dut.set(r_data=0, g_data=1, b_data=0)
                        case 2:
                            dut.set(r_data=0, g_data=0, b_data=1)
                        case 3:
                            dut.set(r_data=1, g_data=1, b_data=1)
            await dut.step()

    dut.start_parallel(_pixel_generator())

    frame_data = np.zeros((16, 32, 3), dtype=np.uint8)

    frames = 0
    pixel_idx = 0
    row_idx = 0
    while frames < 4:
        await dut.step()

        if dut.matrix.out_clk.value:
            for y in [row_idx, row_idx + 8]:
                for c in range(3): # r, g, b
                    frame_data[y, pixel_idx, c] += int(dut.hub_75.value[c])
            # dut.print(f"f={frames} r={row_idx} p={pixel_idx} data={dut.hub_75.value[2:0]}")
            pixel_idx = 0 if pixel_idx >= 31 else pixel_idx + 1
            if dut.hub_75.value[13]:  # lat
                row_idx = 0 if row_idx >= 7 else row_idx + 1

        if dut.frame_complete.value == 1:
            frames += 1
    
    # normalize
    frame_data = frame_data*math.floor(256/frame_data.max())

    outpath = Path(__file__).resolve().parent.parent / "build/test_led_matrix.png"
    imwrite(outpath, frame_data)

    # for x in range(32):
    #     for y in range(16):
    #         if x % 2 == 0:
    #             if y % 4 == 0:
    #                 assert frame_data[y, x, 0] > 0
    #                 assert frame_data[y, x, 1] == 0
    #                 assert frame_data[y, x, 2] == 0
    #         else:
    #             assert frame_data[y,x] == np.vector([0,0,0])
            