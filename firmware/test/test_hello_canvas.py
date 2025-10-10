import pytest
import os
from enum import Enum
import contextlib

class FourState(Enum):
    ZERO = 0
    ONE = 1
    UNSPECIFIED = 'x'
    HIZ = 'z'

    def __bool__(self):
        return self == FourState.ONE

Logic = list[FourState]

class ModuleState():
    def __init__(self, time: float, pins: dict[str, Logic]):
        self.time = time
        self.pins = pins
        for (name, value) in pins.items():
            setattr(self, name, value)

class Pin():
    def __init__(self, name: str, width: int, output: bool):
        self.name = name
        self.width = width
        self.output = output

class Pins():
    def __init__(self, pin_names: list[str]):
        [setattr(self, p, Pin(p, 1, False)) for p in pin_names]

class SVModule():
    def __init__(self, generator, path: str, name: str = None, **kwargs):
        self.path = path
        self.params = kwargs
        self.generator = generator
        # assumes module name from filename if unspecified
        self.name = os.path.basename(path).removesuffix(".sv") if name == None else name
        self._assertions = []
        # TODO: parse module and extract pins
        self.t = 0

    def __enter__(self):
        self.setup = True

    def __exit__(self, exception_type, exception_value, exception_traceback):
        if (exception_value):
            raise exception_value

        self.setup = False
        # begin execution

    @property
    def pins(self) -> Pins:
        return Pins(["clk"])

    def step(self, pin: Pin):
        if self.setup:
            self.t += 1
            return asyncio.as_completed()

    @property
    def state(self) -> ModuleState:
        return ModuleState(self.t, {
            "clk": [FourState.ONE],
            "reset": [FourState.ZERO],
            "row_addr": [FourState.ZERO, FourState.ZERO],
            "pixel_addr": [],
            "r_data": [],
            "g_data": [],
            "b_data": [],
        })
    


class IGenerator():
    def setup(timestep):
        raise RuntimeError("Unimplemented")
    
    @contextlib.contextmanager
    def module(self, path: str, name: str = None, **kwargs):
        m = SVModule(generator=self, path=path, name=name, **kwargs)
        self.append_module(m)
        yield m

    def append_module(module: SVModule):
        raise RuntimeError("Unimplemented")

    def simulate(steps: list[ModuleState]) -> list[ModuleState]:
        raise RuntimeError("Unimplemented")

class IcarusVerilogGenerator(IGenerator):
    def setup(self, timestep):
        pass

    def append_module(self, module: SVModule):
        pass

    def simulate(self, steps: list[ModuleState]) -> list[ModuleState]:
        pass

gen = IcarusVerilogGenerator()

@pytest.mark.asyncio
async def test_hello_canvas():
    with gen.module("../HelloCanvas.sv", SCAN_RATE=8) as m:
        while m.state.reset:
            await m.step(m.pins.clk)
        
        assert m.state.r_data == 0
        assert m.state.g_data == 0
        assert m.state.b_data == 0

        m.set(row_addr=4, pixel_addr=19, req_read=1)
        await m.step(m.pins.clk)
        assert m.pins.r_data == 1
        
        m.set(row_addr=4, pixel_addr=20, req_read=1)
        await m.step(m.pins.clk)
        assert m.pins.r_data == 1

        m.set(req_read=0)
        await m.step(m.pins.clk)
        
        assert m.pins.r_data == 1 # keeps previous value
