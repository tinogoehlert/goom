package main

import "github.com/go-gl/glfw/v3.2/glfw"

type Input struct {
	move chan [3]float32
	turn chan float32
}

func (i *Input) ProcessInput(w *glfw.Window) {
	var move [3]float32
	if w.GetKey(glfw.KeyUp) == glfw.Press {
		move[2] = -10
	}
	if w.GetKey(glfw.KeyDown) == glfw.Press {
		move[2] = 10
	}
	if w.GetKey(glfw.KeyW) == glfw.Press {
		move[1] = 5
	}
	if w.GetKey(glfw.KeyS) == glfw.Press {
		move[1] = -5
	}
}
