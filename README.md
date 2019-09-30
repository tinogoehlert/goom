# GOOM
DOOM engine written in Go

![GOOM](/misc/goom.png?raw=true "GOOM")

Place a `DOOM1.wad` in the root dir and type `make run` to run GOOM.

# Development Status

The project is an experiment and is still lacking a lot of
features, such as shooting, enemy behavior, sound, music, menus, and more.

![DEMO](/misc/goom-preview.gif?raw=true "DEMO")

# Development Setup

Running `make` will initialize go modules. Some of the used modules
use C-bindings and may show compile warnings that can be ignored.

## Linux

On Arch/Manjaro, just install `glbsp`, the rest should be present.

On Ubuntu, install the following system packages:

- libxcursor-dev
- libxrandr-dev
- libxinerama-dev
- libxi-dev
- glbsp

For other systems, please check the Debian/Ubuntu requirements and use your corresponding packages.

## OSX

TBD

## Windows

TBD
