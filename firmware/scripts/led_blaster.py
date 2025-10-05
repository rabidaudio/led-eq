#!/bin/python3

from PIL import Image
import numpy as np

import argparse
import os
import serial
import sys

import code

def netpbm_to_planed_rgb(fobj, color_depth=8, width_out=32, height_out=16):
    """
    Read a NetPBM file from an opened file, crop/pad, and convert to RGB array that's been scrambled
    according to the image format.
    """

    if (color_depth > 8):
        raise IndexError(f"color_depth should be no more than 8, but it's {color_depth}")

    img = Image.open(fobj).convert("RGB")
    w, h = img.size

    # crop / pad to target width and height
    start_x = max(0, (w - width_out) // 2)
    start_y = max(0, (h - height_out) // 2)
    end_x = min(start_x + width_out, w)
    end_y = min(start_y + height_out, h)
    cropped = img.crop((start_x, start_y, end_x, end_y))

    out = Image.new("RGB", (width_out, height_out), (0, 0, 0))
    paste_x = (width_out - cropped.width) // 2
    paste_y = (height_out - cropped.height) // 2
    out.paste(cropped, (paste_x, paste_y))

    top_half = np.array(out)[:height_out//2][:][:]  >> (8 - color_depth)
    bottom_half = np.array(out)[height_out//2:][:][:]  >> (8 - color_depth)

    rgb_halves = np.array([top_half[..., 2].flatten(), top_half[..., 1].flatten(), top_half[..., 0].flatten(),
                           bottom_half[..., 2].flatten(), bottom_half[..., 1].flatten(), bottom_half[..., 0].flatten()])

    packed = bytearray()

    for bit_index in range(color_depth):
        for i in range(len(rgb_halves[0])):
            w = (((rgb_halves[5][i] & 1) << 6) | ((rgb_halves[4][i] & 1) << 5) | ((rgb_halves[3][i] & 1) << 4) |
                 ((rgb_halves[2][i] & 1) << 2) | ((rgb_halves[1][i] & 1) << 1) | ((rgb_halves[0][i] & 1) << 0))
            packed.append(w)
        rgb_halves >>= 1

    #print(packed)

    return bytes(packed)

def main():
    parser = argparse.ArgumentParser(description="Send data over a serial port.")
    parser.add_argument("--port", required=True, help="Serial port (e.g., /dev/ttyUSB0 or COM3)")
    parser.add_argument("--baud", type=int, default=2000000, help="Baud rate (default: 2000000)")
    parser.add_argument("--src", type=str, default=None, help="")
    args = parser.parse_args()

    ser = serial.Serial(port=args.port, baudrate=args.baud, timeout=1)

    # list of bytearrays that can be transmitted to the screen
    frames = []
    # Read images in order from the source directory
    for fname in sorted(os.listdir(args.src)):
        path = os.path.join(args.src, fname)
        if not os.path.isfile(path):
            continue
        with open(path, "rb") as f:
            netpbm_to_planed_rgb(f)

    j = 0
    framecnt = 0
    while (True):
        # Wait for a byte
        received = ser.read(1)
        ser.reset_input_buffer()
        if (received):
            print(f"Got a vsync byte. sending frame {framecnt:7d} / {len(frames):7d}.")
            ser.write(frames[framecnt])
            framecnt = (framecnt + 1) % len(frames)
            print(f"done transmitting frame {framecnt:7d}")

if __name__ == "__main__":
    main()
