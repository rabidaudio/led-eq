#!/bin/sh
yosys -s synth.yosys
nextpnr-ice40 \
    --up5k --package sg48 \
    --randomize-seed \
    --pcf pins.pcf --json build/proj.json \
    --opt-timing --tmg-ripup \
    --asc build/proj.asc
icepack build/proj.asc build/proj.bin
