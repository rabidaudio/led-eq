import optparse
import sys

import requests
from imageio.v3 import imread
from parsel import Selector

FLAG_URL = "https://r74n.com/pixelflags/"


def get_flag(country: str) -> str:
    home = requests.get(FLAG_URL)
    doc = Selector(home.text)
    img_src = doc.css(f"#{country} td img").attrib["src"]
    return f"{FLAG_URL}{img_src}"


def write_memfile(png_uri: str, output, width=32, height=16, x_offset=0, y_offset=1):
    img = imread(png_uri)
    for x in range(width):
        for y in range(height):
            x = x_offset + x
            y = y_offset + y
            hex = bytearray(img[y, x, :3]).hex().upper()
            output.write(f"{hex} // [{x}, {y}]\n")


def main():
    parser = optparse.OptionParser()
    parser.add_option("--width")
    parser.add_option("--height")
    parser.add_option("--x-offset")
    parser.add_option("--y-offset")
    parser.add_option("-o", "--output")

    opts, args = parser.parse_args()

    uri = get_flag(args[0])
    dopts = {k: v for k, v in vars(opts).items() if k != "output" and v != None}
    if opts.output:
        with open(opts.output, "w") as f:
            write_memfile(uri, f, **dopts)
    else:
        write_memfile(uri, sys.stdout, **dopts)


if __name__ == "__main__":
    main()
