package opengl

import (
	"fmt"
	"runtime"
	"time"

	"github.com/tinogoehlert/goom"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//GLRenderer openGL renderer
type GLRenderer struct {
	window        *glfw.Window
	currentLevel  *level
	shaders       map[string]*ShaderProgram
	fbWidth       int
	fbHeight      int
	camera        *Camera
	modelMatrix   mgl32.Mat4
	sprites       spriteList
	inputCallback func(*glfw.Window)
}

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err.Error())
	}
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

// NewRenderer initialize the renderer
func NewRenderer() (*GLRenderer, error) {
	return &GLRenderer{
		shaders:     make(map[string]*ShaderProgram),
		camera:      NewCamera(),
		modelMatrix: mgl32.Ident4(),
	}, nil
}

// LoadShaderProgram loads a shader program
func (gr *GLRenderer) LoadShaderProgram(name, vertFile, fragFile string) error {
	var shader = NewShaderProgram()
	if err := shader.AddVertexShader(vertFile); err != nil {
		return err
	}

	if err := shader.AddFragmentShader(fragFile); err != nil {
		return err
	}

	if err := shader.Link(); err != nil {
		return err
	}
	gr.shaders["main"] = shader
	return nil
}

// BuildLevel builds the level
func (gr *GLRenderer) BuildLevel(m *goom.Map, gfx *goom.Graphics) {
	gr.currentLevel = RegisterMap(m, gfx)
}

// BuildLevel builds the level
func (gr *GLRenderer) BuildSprites(gfx *goom.Graphics) {
	gr.sprites = BuildSpritesFromGfx(gfx)
}

func (gr *GLRenderer) Camera() *Camera {
	return gr.camera
}

func (gr *GLRenderer) SetInputLoop(fn func(*glfw.Window)) {
	gr.inputCallback = fn
}

// CreateWindow creates the GL window
func (gr *GLRenderer) CreateWindow(width, height int, title string) (err error) {
	gr.window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return fmt.Errorf("could not create Window: %s", err.Error())
	}

	gr.window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		return fmt.Errorf("could not init GL: %s", err.Error())
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	return nil
}

// KeyPressed set callback for keypress
func (gr *GLRenderer) KeyPressed(keyPressed func(key int)) {
	gr.window.SetKeyCallback(func(
		w *glfw.Window,
		key glfw.Key,
		scancode int,
		action glfw.Action,
		mods glfw.ModifierKey) {
		keyPressed(int(key))
	})
}

func (gr *GLRenderer) setLight(light float32) {
	gr.shaders["main"].Uniform1f("sectorLight", light)
}

func (gr *GLRenderer) setProjection() {
	gr.shaders["main"].UniformMatrix4fv("projection",
		mgl32.Perspective(64.0, float32(gr.fbWidth)/float32(gr.fbHeight), 1.0, 10000.0),
	)
}

func (gr *GLRenderer) setView() {
	gr.shaders["main"].UniformMatrix4fv("view", gr.camera.ViewMat4())
}

func (gr *GLRenderer) setModel() {
	gr.shaders["main"].UniformMatrix4fv("model", gr.modelMatrix)
}

func (gr *GLRenderer) DrawSubSector(idx int) {
	var s = gr.currentLevel.subSectors[idx]
	gr.setLight(s.sector.LightLevel())
	s.Draw()
}

func (gr *GLRenderer) GetSectorForSSect(ssect *goom.SubSector) goom.Sector {
	var (
		fseg   = ssect.Segments()[0]
		line   = gr.currentLevel.mapRef.LinesDefs[fseg.GetLineDef()]
		side   = gr.currentLevel.mapRef.SideDefs[line.Right]
		sector = gr.currentLevel.mapRef.Sectors[side.Sector]
	)
	return sector
}

func (gr *GLRenderer) DrawThingsInBBox(sector goom.Sector, bbox goom.BBox) {
	gr.shaders["main"].Uniform1i("draw_phase", 1)
	gr.shaders["main"].Uniform2f("billboard_size", mgl32.Vec2{0.57, 0.57})

	for _, t := range gr.currentLevel.mapRef.Things {
		x, y := float32(t.X), float32(t.Y)
		if bbox.PosInBox(x, y) {
			if t.Type == 3004 {
				spr := gr.sprites["POSS"]
				gr.shaders["main"].Uniform3f("billboard_pos", mgl32.Vec3{-x, sector.FloorHeight() + spr.median + 4, y})
				gr.sprites["POSS"].mesh.DrawMesh(gl.TRIANGLES)
			}
			if t.Type == 3001 {
				spr := gr.sprites["TROO"]
				gr.shaders["main"].Uniform3f("billboard_pos", mgl32.Vec3{-x, sector.FloorHeight() + spr.median + 4, y})
				gr.sprites["TROO"].mesh.DrawMesh(gl.TRIANGLES)
			}
			if t.Type == 48 {
				spr := gr.sprites["ELEC"]
				gr.shaders["main"].Uniform3f("billboard_pos", mgl32.Vec3{-x, sector.FloorHeight() + spr.median + 6, y})
				gr.sprites["ELEC"].mesh.DrawMesh(gl.TRIANGLES)
			}
			if t.Type == 2018 {
				spr := gr.sprites["ARM1"]
				gr.shaders["main"].Uniform3f("billboard_pos", mgl32.Vec3{-x, sector.FloorHeight() + spr.median + 4, y})
				gr.sprites["ARM1"].mesh.DrawMesh(gl.TRIANGLES)
			}
		}
	}
	gr.shaders["main"].Uniform1i("draw_phase", 0)

}

func (gr *GLRenderer) DrawHUD() {
	Ortho := mgl32.Ortho2D(0, float32(gr.fbWidth), -float32(gr.fbHeight), 0)

	gl.Disable(gl.DEPTH_TEST) // Disable the Depth-testing

	gr.shaders["main"].UniformMatrix4fv("ortho", Ortho)
	gr.shaders["main"].Uniform1i("draw_phase", 2)
	gr.shaders["main"].Uniform2f("billboard_size", mgl32.Vec2{1.4, 1.4})
	gr.shaders["main"].Uniform3f("billboard_pos", mgl32.Vec3{float32(gr.fbWidth) / 2, -float32(gr.fbHeight), 0})
	gr.sprites["CHGG"].mesh.DrawMesh(gl.TRIANGLES)
	gr.shaders["main"].Uniform1i("draw_phase", 0)
}

// Loop starts the render loop
func (gr *GLRenderer) Loop(fps int, gameCB func(win *glfw.Window)) {
	for !gr.window.ShouldClose() {
		gr.fbWidth, gr.fbHeight = gr.window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(gr.fbWidth), int32(gr.fbHeight))
		// Do OpenGL stuff.
		gl.Enable(gl.DEPTH_TEST)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gr.shaders["main"].Use()
		gr.setProjection()
		gr.setView()
		gr.setModel()

		gr.currentLevel.mapRef.WalkBsp(goom.GLNodesName, func(i int, n *goom.Node, b goom.BBox) {
			rS := gr.currentLevel.subSectors[i]
			gr.DrawSubSector(i)
			gr.DrawThingsInBBox(rS.sector, b)

		})

		gr.DrawHUD()
		gr.window.SwapBuffers()
		glfw.PollEvents()
		gameCB(gr.window)
		time.Sleep(20 * time.Millisecond)
	}
}

/*
ticker := time.NewTicker(int64(second) / 60) // max 60 fps
*/
