SystemVerilog code for controlling HUB75 LED matrix displays via SPI. Designed for Lattice ice40 chipset.

Project management and testing in Python (using [poetry](https://python-poetry.org/docs/) and [cocotb](https://docs.cocotb.org)).

## Commands

```bash
poetry fix # format code
poetry test # run test suite
poetry synth # run synthesizer
```

## Setup

```bash
python3 -m venv .venv
pipx install poetry
poetry install
```
