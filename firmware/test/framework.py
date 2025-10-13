import os
import cocotb
from cocotb.triggers import RisingEdge, Timer
from cocotb.types import Array
import contextlib
import logging
import functools
from pathlib import Path
from cocotb_tools.runner import get_runner

# A global list to put runners from all tests
test_runners = []

# A wrapper for cocotb's dut module that adds some helper functions
class DutWrapper(cocotb.handle.HierarchyObject):
    def __init__(self, dut, period = 10, total_cycles = 250_000, clk_pin = "clk", reset_pin = "reset"):
        self.__class__ = type(dut.__class__.__name__,
                              (self.__class__, dut.__class__),
                              {})
        self.__dict__ = dut.__dict__.copy()
        self._dut = dut
        self._clk_pin = clk_pin
        self._reset_pin = reset_pin
        self._period = period
        self._total_cycles = total_cycles
        self._logger = logging.getLogger(dut.__class__.__name__)
        self._logger.setLevel(logging.INFO)

    async def step(self, n = 1, pin: str = None, direction = RisingEdge):
        pin = self._pin_by_name(pin or self._clk_pin)
        for _ in range(n):
            await direction(pin)
        await Timer(1, unit="step") # step one cycle more to let things settle

    def _pin_by_name(self, name: str):
        return getattr(self._dut, name)
    
    def set(self, **kwargs):
        for (key, val) in kwargs.items():
            self._pin_by_name(key).value = val

    async def step_reset(self):
        await self.step(n=2)

    def print(self, *args, **kwargs):
        self._logger.info(*args, **kwargs)

    @property
    def state(self):
        res = {}
        for key in self._dut._keys():
            if key == key.upper():
                continue # ignore properties
            val = getattr(self._dut, key).value
            if type(val) is Array:
                if len(val) > 16:
                    "".join([str(v) for v in val[0:16]]) + f"...({len(val)})"
                else:
                    "".join([str(v) for v in val])
            else:
                res[key] = str(val)
        return res

    async def generate_clock(self):
        for _ in range(self._total_cycles):
            self._pin_by_name(self._clk_pin).value = 0
            await Timer(self._period/2, unit="step")
            self._pin_by_name(self._clk_pin).value = 1
            await Timer(self._period/2, unit="step")

    async def generate_reset(self):
        self._pin_by_name(self._reset_pin).value = 1
        await self.step(2)
        self._pin_by_name(self._reset_pin).value = 0

@contextlib.contextmanager
def simulator_module(dut, reset = True, **kwargs):
    wdut = DutWrapper(dut, **kwargs)
    cocotb.start_soon(wdut.generate_clock())
    if reset:
        cocotb.start_soon(wdut.generate_reset())

    yield wdut

def simulator_test(
        path: str,
        module_name: str = None,
        simulator: str = "icarus",
        proj_path = None,
        timescale = ("10ns", "100ps"),
        **kwargs
):
    proj_path = Path(__file__).resolve().parent.parent if proj_path == None else proj_path
    module_name = os.path.basename(path).removesuffix(".sv") if module_name == None else module_name
    def _decorator(func):
        test_name = func.__name__

        @cocotb.test(name=test_name, **kwargs)
        @functools.wraps(func)
        async def _test_wrapper(*args, **w_kwargs):
            with simulator_module(args[0], **kwargs) as wdut:
                return await func(wdut, **w_kwargs)
        _test_wrapper.__name__ = test_name

        def _runner():
            sources = [proj_path / path]

            runner = get_runner(simulator)
            runner.build(
                sources=sources,
                hdl_toplevel=module_name,
                timescale=timescale,
                build_dir=proj_path,
                waves=True,
            )
            runner.test(
                hdl_toplevel=module_name,
                test_module=f"{test_name},",
                timescale=timescale,
                waves=True,
            )

        test_runners.append(_runner)

        return _test_wrapper
    return _decorator
