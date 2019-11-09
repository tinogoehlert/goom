# GOOM
DOOM engine written in Go

![GOOM](/resources/press/goom.png?raw=true "GOOM")

Make sure you have `DOOM1.WAD` in the root dir. Then type `make run` to run GOOM.

# Development Status

The project is an experiment and is still lacking a lot of
features, such as shooting, enemy behavior, sound, menus, and more.

![DEMO](/resources/press/goom-preview.gif?raw=true "DEMO")

# Development Setup

Running `make` will initialize go modules and run the tests.

For testing, it is also useful to run `make test-run`, which starts the game,
loads the WADs, runs the event loop once, and then exits automatically.

Note that some of the used modules use C-bindings and may show compile warnings.
Please ignore them.

## Linux

On Arch/Manjaro install `glbsp`
and setup [TiMidity](https://wiki.archlinux.org/index.php/Timidity#Installation).

On Ubuntu, install the following system packages:

- libxcursor-dev
- libxrandr-dev
- libxinerama-dev
- libxi-dev
- glbsp
- timidity
- libportmidi-dev
- librtmidi-dev

For other systems, please check the above requirements and use your corresponding packages.

## OSX

TBD

## Windows

TBD

## DOOM1.WAD

This project includes a copy of the [shareware version of DOOM](https://doomwiki.org/wiki/DOOM1.WAD) licensed under the [Original DOOM1 Shareware License](DOOM1.LICENSE).
