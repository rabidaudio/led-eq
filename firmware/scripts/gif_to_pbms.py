#!/usr/bin/env python3
import sys, os
from PIL import Image

from PIL import Image
import numpy as np

def gif_to_ppm(gif_path, outdir, width_out=32, height_out=16):
    os.makedirs(outdir, exist_ok=True)
    with Image.open(gif_path) as im:
        frame = 0
        while True:
            rgb = im.convert("RGB")

            # Downsample so that smaller side = 32
            w, h = rgb.size
            scale = min(w / width_out, h / height_out)
            new_w, new_h = int(round(w * scale)), int(round(h * scale))
            resized = rgb.resize((new_w, new_h), Image.LANCZOS)

            # Center-crop to w x h
            left = (new_w - width_out) // 2
            top = (new_h - height_out) // 2
            cropped = resized.crop((left, top, left + width_out, top + height_out))

            # Save as PPM (Netpbm RGB)
            fname = os.path.join(outdir, f"frame{frame:03d}.ppm")
            cropped.save(fname, format="PPM")

            try:
                im.seek(im.tell() + 1)
                frame += 1
            except EOFError:
                break

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(f"Converts a gif to a directory full of pbm bitmap files.")
        print(f"You can run this script on a gif and then run led_blaster.py on the resulting dir.")
        print(f"Usage: {sys.argv[0]} input.gif output_dir")
        sys.exit(1)
    gif_to_ppm(sys.argv[1], sys.argv[2])
