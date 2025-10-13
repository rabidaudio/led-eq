import random

from framework import simulator_test


@simulator_test("SpiSlaveSimplex.sv", clk_pin="sck")
async def test_sss_data(dut):
    dut.set(cs=0)  # active low

    for i in range(100):
        r = random.randint(0, 255)
        bits = list(bin(r)[2:].zfill(8))
        for b in range(8):
            dut.set(mosi=bits[b])
            await dut.step()
            if b < 7:
                assert dut.data_avail.value == 0

        assert dut.data_avail.value == 1
        assert dut.data.value == r
        assert dut.byte_index.value == i


@simulator_test("SpiSlaveSimplex.sv", clk_pin="sck")
async def test_sss_cs(dut):
    dut.set(cs=1, mosi=1)

    # make sure no bytes come out when cs is high
    await dut.step(n=16)

    assert dut.data_avail.value == 0
    assert dut.byte_index.value == 0

    # spit out 4 whole bytes
    dut.set(cs=0)
    await dut.step(n=4 * 8)

    assert dut.data_avail.value == 1
    assert dut.byte_index.value == 4

    # start sending out another byte, but disable cs mid-byte
    bits = random.randint(1, 6)
    await dut.step(n=bits)

    # should hold last value
    assert dut.data_avail.value == 0
    assert dut.byte_index.value == 4

    dut.set(cs=1)
    await dut.step()

    # should still be reset
    assert dut.data_avail.value == 0
    assert dut.byte_index.value == 0

    # should stay that way
    for _ in range(100):
        await dut.step()

        assert dut.data_avail.value == 0
        assert dut.byte_index.value == 0
