from pathlib import Path

from cocotb_tools.runner import get_runner

def test_runner():
    sim = "icarus"

    proj_path = Path(__file__).resolve().parent.parent

    sources = [proj_path / "HelloCanvas.sv"]

    runner = get_runner(sim)
    runner.build(
        sources=sources,
        hdl_toplevel="HelloCanvas",
        timescale=("10ns", "100ps"),
        build_dir=proj_path,
    )

    runner.test(hdl_toplevel="HelloCanvas", test_module="test_hello_canvas,", timescale=("10ns", "100ps"))


if __name__ == "__main__":
    test_runner()
