import contextlib
import functools
import logging
import os
from pathlib import Path

import cocotb
from cocotb.triggers import RisingEdge, Timer
from cocotb.types import Array
from cocotb_tools.runner import get_runner

# A global list to put runners from all tests
test_runners = {}


# A wrapper for cocotb's dut module that adds some helper functions
class DutWrapper(cocotb.handle.HierarchyObject):
    def __init__(
        self, dut, period=10, total_cycles=250_000, clk_pin="clk", reset_pin="reset"
    ):
        self.__class__ = type(
            dut.__class__.__name__, (self.__class__, dut.__class__), {}
        )
        self.__dict__ = dut.__dict__.copy()
        self._dut = dut
        self._clk_pin = clk_pin
        self._reset_pin = reset_pin
        self._period = period
        self._total_cycles = total_cycles
        self._logger = logging.getLogger(dut.__class__.__name__)
        self._logger.setLevel(logging.INFO)

    async def step(self, n=1, pin=None, direction=RisingEdge):
        pin = self._pin_by_name(self._clk_pin) if pin is None else pin
        for _ in range(n):
            await direction(pin)
        await self.tick()  # step one cycle more to let things settle

    async def tick(self, n=1):
        await Timer(n, unit="step")

    def _pin_by_name(self, name: str):
        return getattr(self._dut, name)

    def set_now(self, **kwargs):
        for key, val in kwargs.items():
            self._pin_by_name(key).value = val

    def set(self, **kwargs):
        async def _wait_set():
            await self.tick()
            self.set_now(**kwargs)

        self.start_parallel(_wait_set())

    async def step_reset(self):
        await self.step(n=2)

    def print(self, *args, **kwargs):
        self._logger.info(*args, **kwargs)

    @property
    def state(self):
        res = {}
        for key in self._dut._keys():
            if key == key.upper():
                continue  # ignore properties
            val = getattr(self._dut, key).value
            if type(val) is Array:
                if len(val) > 16:
                    "".join([str(v) for v in val[0:16]]) + f"...({len(val)})"
                else:
                    "".join([str(v) for v in val])
            else:
                res[key] = str(val)
        return res

    def start_parallel(self, task):
        return cocotb.start_soon(task)

    async def generate_clock(self):
        for _ in range(self._total_cycles):
            self._pin_by_name(self._clk_pin).value = 0
            await Timer(self._period / 2, unit="step")
            self._pin_by_name(self._clk_pin).value = 1
            await Timer(self._period / 2, unit="step")

    async def generate_reset(self):
        self._pin_by_name(self._reset_pin).value = 1
        await self.step(2)
        self._pin_by_name(self._reset_pin).value = 0


@contextlib.contextmanager
def simulator_module(dut, reset=True, **kwargs):
    wdut = DutWrapper(dut, **kwargs)
    cocotb.start_soon(wdut.generate_clock())
    if reset:
        cocotb.start_soon(wdut.generate_reset())

    yield wdut


def simulator_test(
    path: str,
    module_name: str = None,
    params: dict = {},
    simulator: str = "icarus",
    proj_dir=None,
    timescale=("10ns", "100ps"),
    reset=True,
    period=10,
    total_cycles=250_000,
    clk_pin="clk",
    reset_pin="reset",
    **kwargs,
):
    proj_dir = Path(__file__).resolve().parent.parent if proj_dir == None else proj_dir
    module_name = (
        os.path.basename(path).removesuffix(".sv")
        if module_name == None
        else module_name
    )

    def _decorator(func):
        test_name = func.__name__
        test_file_name = func.__module__

        @cocotb.test(name=test_name, **kwargs)
        @functools.wraps(func)
        async def _test_wrapper(*args, **w_kwargs):
            __tracebackhide__ = True  # Hide the traceback when using pytest
            with simulator_module(
                args[0],
                reset=reset,
                period=period,
                total_cycles=total_cycles,
                clk_pin=clk_pin,
                reset_pin=reset_pin,
            ) as wdut:
                if reset:
                    await wdut.step_reset()
                return await func(wdut, **w_kwargs)

        _test_wrapper.__name__ = test_name
        _test_wrapper.__tracebackhide__ = True

        if test_file_name not in test_runners:

            def _runner():
                sources = [proj_dir / path]

                build_dir = proj_dir / "build"

                runner = get_runner(simulator)
                runner.build(
                    sources=sources,
                    hdl_toplevel=module_name,
                    parameters=params,
                    build_args=["-DSIMULATION"],
                    timescale=timescale,
                    build_dir=build_dir,
                    cwd=proj_dir,
                    always=True,
                    waves=True,
                )
                runner.test(
                    hdl_toplevel=module_name,
                    test_module=f"{test_file_name},",
                    timescale=timescale,
                    waves=True,
                )

            test_runners[test_file_name] = _runner

        return _test_wrapper

    return _decorator
