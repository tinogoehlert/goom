package opengl

import (
	"fmt"
	"io/ioutil"

	"github.com/go-gl/gl/v2.1/gl"
)

// ShaderProgram represents a shader program
type ShaderProgram struct {
	id       uint32
	uniforms map[string]int32
}

// NewShaderProgram creates a new shader program
func NewShaderProgram() *ShaderProgram {
	return &ShaderProgram{
		id:       gl.CreateProgram(),
		uniforms: make(map[string]int32),
	}
}

func (sp *ShaderProgram) Use() {
	gl.UseProgram(sp.id)
}

// AddFragmentShader add fragment shader from file to program
func (sp *ShaderProgram) AddFragmentShader(file string) error {
	id, err := compileShader(file, gl.FRAGMENT_SHADER)
	if err != nil {
		return fmt.Errorf("FRAG: %s", err.Error())
	}
	gl.AttachShader(sp.id, id)
	return nil
}

// Uniform1f set float32 uniform variable
func (sp *ShaderProgram) Uniform1f(name string, val float32) {
	if _, ok := sp.uniforms[name]; !ok {
		sp.uniforms[name] = gl.GetUniformLocation(sp.id, gl.Str(name+"\x00"))
	}
	gl.Uniform1f(sp.uniforms[name], val)
}

// Uniform1i set int uniform variable
func (sp *ShaderProgram) Uniform1i(name string, val int) {
	if _, ok := sp.uniforms[name]; !ok {
		sp.uniforms[name] = gl.GetUniformLocation(sp.id, gl.Str(name+"\x00"))
	}
	gl.Uniform1i(sp.uniforms[name], int32(val))
}

// Uniform1i set int uniform variable
func (sp *ShaderProgram) Uniform3f(name string, vec3 [3]float32) {
	if _, ok := sp.uniforms[name]; !ok {
		sp.uniforms[name] = gl.GetUniformLocation(sp.id, gl.Str(name+"\x00"))
	}
	gl.Uniform3fv(sp.uniforms[name], 1, &vec3[0])
}

// Uniform1i set int uniform variable
func (sp *ShaderProgram) Uniform2f(name string, vec2 [2]float32) {
	if _, ok := sp.uniforms[name]; !ok {
		sp.uniforms[name] = gl.GetUniformLocation(sp.id, gl.Str(name+"\x00"))
	}
	gl.Uniform2fv(sp.uniforms[name], 1, &vec2[0])
}

// UniformMatrix4fv set 4x4 matrix float32 uniform variable
func (sp *ShaderProgram) UniformMatrix4fv(name string, mat4 [16]float32) {
	if _, ok := sp.uniforms[name]; !ok {
		sp.uniforms[name] = gl.GetUniformLocation(sp.id, gl.Str(name+"\x00"))
	}
	gl.UniformMatrix4fv(sp.uniforms[name], 1, false, &mat4[0])
}

// AddVertexShader add fragment shader from file to program
func (sp *ShaderProgram) AddVertexShader(file string) error {
	id, err := compileShader(file, gl.VERTEX_SHADER)
	if err != nil {
		return fmt.Errorf("VERT: %s", err.Error())
	}
	gl.AttachShader(sp.id, id)
	return nil
}

// Link links all shaders to program
func (sp *ShaderProgram) Link() error {
	gl.LinkProgram(sp.id)
	var (
		errLog = [512]uint8{}
		length int32
	)

	gl.GetProgramInfoLog(sp.id, 512, &length, &errLog[0])
	if length > 0 {
		return fmt.Errorf("%s", errLog[:])
	}

	return nil
}

func compileShader(file string, shaderType uint32) (uint32, error) {
	var (
		errLog = [512]uint8{}
		length int32
	)

	src, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}
	glSrc, freeFn := gl.Strs(string(src) + "\x00")
	defer freeFn()

	id := gl.CreateShader(shaderType)
	gl.ShaderSource(id, 1, glSrc, nil)
	gl.CompileShader(id)
	gl.GetShaderInfoLog(id, 512, &length, &errLog[0])
	if length > 0 {
		return 0, fmt.Errorf(string(errLog[:]))
	}
	return id, nil
}
