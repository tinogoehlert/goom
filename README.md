# GOOM
DOOM engine written in Go

currently WIP

# Development Setup

## Linux

First install the following system packages:

- libxcursor-dev
- libxrandr-dev
- libxinerama-dev
- libxi-dev

These are the package names on Debian/Ubuntu. On other distributions use corresponding packages.

Then run `go get github.com/go-gl/glfw/v3.2/glfw`.
This may result in errors that can be fixed as follows.

Until `go-gl` adds `glfw/v3.3` support, you now need to hack your `linux_joystick.c`. It is located on your Go `src` dir at `github.com/go-gl/glfw/v3.2/glfw/glfw/src/linux_joystick.c`.
Change any `path[20]` statement to `path[512]` in the code directly.
Then run `go get github.com/go-gl/glfw/v3.2/glfw` again.

Finally install the remaining Go dependecies:

- github.com/go-gl/mathgl/mgl32
- github.com/go-gl/gl/v2.1/gl

The run GOOM: `go run cmd/doom/main.go`
