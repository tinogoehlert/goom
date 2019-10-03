package opengl

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/tinogoehlert/goom/cmd/doom/internal/game"
	"github.com/tinogoehlert/goom/graphics"

	"github.com/tinogoehlert/goom/level"

	"github.com/tinogoehlert/goom"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//GLRenderer openGL renderer
type GLRenderer struct {
	window        *glfw.Window
	currentLevel  *doomLevel
	shaders       map[string]*ShaderProgram
	fbWidth       int
	fbHeight      int
	camera        *Camera
	modelMatrix   mgl32.Mat4
	textures      glTextureStore
	spriter       *glSpriter
	currentShader string
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
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

// NewRenderer initialize the renderer
func NewRenderer() (*GLRenderer, error) {
	return &GLRenderer{
		shaders:       make(map[string]*ShaderProgram),
		camera:        NewCamera(),
		modelMatrix:   mgl32.Ident4(),
		currentShader: "main",
		textures:      newGLTextureStore(),
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
	gr.shaders[name] = shader
	return nil
}

// SetShaderProgram sets shader program
func (gr *GLRenderer) SetShaderProgram(name string) error {
	gr.currentShader = name
	return nil
}

// BuildLevel builds the level
func (gr *GLRenderer) BuildLevel(m *level.Level, gd *goom.GameData) {
	gr.currentLevel = RegisterMap(m, gd, gr.textures)
}

// BuildLevel builds the level
func (gr *GLRenderer) BuildGraphics(gd *goom.GameData) {
	for k, v := range gd.Textures {
		gr.textures.initTexture(k, 1)
		gr.textures.addTexture(k, 0, v)
	}

	for k, v := range gd.Flats {
		gr.textures.initTexture(k, 1)
		gr.textures.addTexture(k, 0, v[0])
	}

	for _, v := range gd.Sprites {
		v.Frames(func(f *graphics.SpriteFrame) {
			gr.textures.initTexture(f.Name(), len(f.Angles()))
			if strings.Contains(f.Name(), "TROO") {
				fmt.Println(f.Name())
			}
			for i, img := range f.Angles() {
				gr.textures.addTexture(f.Name(), i, img)
			}
		})
	}
	gr.spriter = NewSpriter()
	//gr.sprites = BuildSpritesFromGfx(gd.Sprites, gd.DefaultPalette())
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
	gr.shaders[gr.currentShader].Uniform1f("sectorLight", light)
}

func (gr *GLRenderer) setProjection() {
	gr.shaders[gr.currentShader].UniformMatrix4fv("projection",
		mgl32.Perspective(64.0, float32(gr.fbWidth)/float32(gr.fbHeight), 1.0, 10000.0),
	)
}

func (gr *GLRenderer) setView() {
	gr.shaders[gr.currentShader].UniformMatrix4fv("view", gr.camera.ViewMat4())
}

func (gr *GLRenderer) setModel() {
	gr.shaders[gr.currentShader].UniformMatrix4fv("model", gr.modelMatrix)
}

func (gr *GLRenderer) DrawSubSector(idx int) {
	var s = gr.currentLevel.subSectors[idx]
	gr.setLight(s.sector.LightLevel())
	s.Draw(gr.textures)
}

func (gr *GLRenderer) GetSectorForSSect(ssect *level.SubSector) level.Sector {
	var (
		fseg   = ssect.Segments()[0]
		line   = gr.currentLevel.mapRef.LinesDefs[fseg.LineDef()]
		side   = gr.currentLevel.mapRef.SideDefs[line.Right]
		sector = gr.currentLevel.mapRef.Sectors[side.Sector]
	)
	return sector
}

func (gr *GLRenderer) DrawThings(things []*game.DoomThing) {
	gr.shaders[gr.currentShader].Uniform1i("draw_phase", 1)

	for _, t := range things {
		f := t.NextFrame()
		a, flipped := t.CalcAngle(gr.camera.position)
		img := gr.textures.Get(t.SpriteName()+string(f), a)
		gr.shaders[gr.currentShader].Uniform3f("billboard_pos", mgl32.Vec3{
			-t.Position()[0],
			t.Height() + (float32(img.image.Height()) / 2) + 4,
			t.Position()[1],
		})

		gr.shaders[gr.currentShader].Uniform1i("billboard_flipped", flipped)
		gr.shaders[gr.currentShader].Uniform2f("billboard_size", mgl32.Vec2{
			float32(img.image.Width()) / 90,
			float32(img.image.Height()) / 100,
		})
		gr.spriter.Draw(gl.TRIANGLES, img)
	}
	gr.shaders[gr.currentShader].Uniform1i("draw_phase", 0)
}

// DrawHUD draws the game hud
func (gr *GLRenderer) DrawHUD() {
	Ortho := mgl32.Ortho2D(0, float32(gr.fbWidth), float32(-gr.fbHeight), 0)

	gl.Disable(gl.DEPTH_TEST) // Disable the Depth-testing
	gr.shaders[gr.currentShader].UniformMatrix4fv("ortho", Ortho)
	gr.shaders[gr.currentShader].Uniform1i("draw_phase", 2)
	gr.shaders[gr.currentShader].Uniform2f("billboard_size", mgl32.Vec2{float32(gr.fbHeight) / 420, float32(gr.fbHeight) / 420})
	gr.shaders[gr.currentShader].Uniform3f("billboard_pos", mgl32.Vec3{float32(gr.fbWidth) / 2, -float32(gr.fbHeight) + 170, 0})
	img := gr.textures.Get("PISGA", 0)
	gr.spriter.Draw(gl.TRIANGLES, img)
	gr.shaders[gr.currentShader].Uniform1i("draw_phase", 0)
}

// Loop starts the render loop
func (gr *GLRenderer) Loop(fps int, thingPass func(), inputCB func(win *glfw.Window)) {
	frameTime := time.Duration(1000 / fps)
	for !gr.window.ShouldClose() {
		t0 := time.Now()
		gr.fbWidth, gr.fbHeight = gr.window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(gr.fbWidth), int32(gr.fbHeight))
		// Do OpenGL stuff.
		gl.Enable(gl.DEPTH_TEST)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gr.shaders[gr.currentShader].Use()
		gr.setProjection()
		gr.setView()
		gr.setModel()

		gr.currentLevel.mapRef.WalkBsp(level.GLNodesName, func(i int, n *level.Node, b level.BBox) {
			gr.DrawSubSector(i)
		})

		thingPass()
		gr.DrawHUD()
		gr.window.SwapBuffers()
		glfw.PollEvents()
		inputCB(gr.window)
		remaining := frameTime - time.Now().Sub(t0)
		if remaining > 0 {
			time.Sleep(remaining)
		}
	}
}
