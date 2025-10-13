import numpy as np
from framework import simulator_test


@simulator_test("HelloCanvas.sv")
async def test_hello_canvas(dut):
    await dut.step_reset()

    assert dut.r_data.value == 0, "should initialize to zero after reset"
    assert dut.g_data.value == 0
    assert dut.b_data.value == 0

    dut.set(req_read=1)

    results = np.empty((16, 32))

    for row in reversed(range(8)):
        for pixel in reversed(range(32)):
            dut.set(row_addr=row, pixel_addr=pixel)
            await dut.step()
            results[(row, pixel)] = dut.r_data.value[0]
            results[(row + 8, pixel)] = dut.r_data.value[1]

    expected = np.empty((16, 32))
    with open("../hello.b.mem", "r") as f:
        for r, line in enumerate(f):
            for p, v in enumerate(line.split(" ")):
                expected[(r, p)] = v

    np.testing.assert_equal(results, expected, "should return the same data as input")
