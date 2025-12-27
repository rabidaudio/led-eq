import argparse
import os
import shutil
import subprocess
import sys
from pathlib import Path


def synthesize(module_path, module_name, pins_path, build_dir=None):
    if build_dir is None:
        build_dir = Path(os.getcwd())
    if type(build_dir) == str:
        build_dir = Path(build_dir).resolve()

    build_name = os.path.basename(module_path).removesuffix(".sv")

    synth_path = build_dir / f"{build_name}.yosys"
    synth_out = build_dir / f"{build_name}.json"
    asc_out = build_dir / f"{build_name}.asc"
    bin_out = build_dir / f"{build_name}.bin"
    with open(synth_path, "w") as f:
        f.write(f"read_verilog -sv {module_path}\n")
        f.write(f"synth_ice40 -top {module_name} -json {synth_out}\n")
    subprocess.run(["yosys", "-s", synth_path], check=True)
    subprocess.run(
        [
            "nextpnr-ice40",
            "--up5k",
            "--package",
            "sg48",
            "--randomize-seed",
            "--pcf",
            pins_path,
            "--json",
            synth_out,
            "--opt-timing",
            "--tmg-ripup",
            "--asc",
            asc_out,
        ],
        check=True,
    )
    subprocess.run(["icepack", asc_out, bin_out], check=True)
    return bin_out


def load(bin):
    match sys.platform:
        case "darwin":
            shutil.copyfile(bin, "/Volumes/iCE40/")
        case _:
            raise NotImplementedError()


def main():
    parser = argparse.ArgumentParser(
        description="Synthesize a top SystemVerilog module for the iCE40"
    )
    parser.add_argument(
        "-m",
        "--module-name",
        dest="module_name",
        default="top",
        help="The name of the module",
    )
    parser.add_argument(
        "-b",
        "--build-dir",
        dest="build_dir",
        help="Path to directory to use for outputs",
    )
    parser.add_argument(
        "-p", "--pins", dest="pins_path", help="Path to pin definition file (.pcf)"
    )
    parser.add_argument(
        "-l",
        "--load",
        action="store_true",
        dest="load",
        help="Load the resulting synthesis to the iCE40 board",
    )
    parser.add_argument("module_path", help="Path to top level module (.sv)")

    args = parser.parse_args()

    bin = synthesize(
        args.module_path,
        pins_path=args.pins_path,
        build_dir=args.build_dir,
        module_name=args.module_name,
    )
    print(f"Synthesized {bin}")
    if args.load:
        load(bin)
        print("Loaded to board")


if __name__ == "__main__":
    main()
