package opengl

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/tinogoehlert/goom/game"
	"github.com/tinogoehlert/goom/goom"
	"github.com/tinogoehlert/goom/graphics"
	"github.com/tinogoehlert/goom/level"
)

//GLRenderer openGL renderer
type GLRenderer struct {
	currentLevel  *doomLevel
	shaders       map[string]*ShaderProgram
	fbWidth       int
	fbHeight      int
	camera        *Camera
	modelMatrix   mgl32.Mat4
	textures      glTextureStore
	spriter       *glSpriter
	lastTick      time.Time
	currentShader string
}

// Init initialize glfw
func Init() error {
	if err := gl.Init(); err != nil {
		return fmt.Errorf("could not init GL: %s", err.Error())
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	return nil
}

// NewRenderer initialize the renderer
func NewRenderer(gd *goom.GameData) (*GLRenderer, error) {
	gr := &GLRenderer{
		shaders:       make(map[string]*ShaderProgram),
		camera:        NewCamera(),
		modelMatrix:   mgl32.Ident4(),
		currentShader: "main",
		textures:      newGLTextureStore(),
	}

	for k, v := range gd.Textures {
		gr.textures.initTexture(k, 1)
		gr.textures.addTexture(k, 0, v)
	}

	for k, v := range gd.Flats {
		gr.textures.initTexture(k, 1)
		gr.textures.addTexture(k, 0, v[0])
	}

	for k, v := range gd.Fonts.GetAllGraphics() {
		gr.textures.initTexture(k, 1)
		gr.textures.addTexture(k, 0, v)
	}

	for _, v := range gd.Sprites {
		v.Frames(func(f *graphics.SpriteFrame) {
			gr.textures.initTexture(f.Name(), len(f.Angles()))
			for i, img := range f.Angles() {
				gr.textures.addTexture(f.Name(), i, img)
			}
		})
	}
	gr.spriter = NewSpriter()
	return gr, nil
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
func (gr *GLRenderer) LoadLevel(m *level.Level, gd *goom.GameData) {
	gr.currentLevel = RegisterMap(m, gd, gr.textures)
}

func (gr *GLRenderer) Camera() *Camera {
	return gr.camera
}

func (gr *GLRenderer) SetLight(light float32) {
	gr.shaders[gr.currentShader].Uniform1f("sectorLight", light)
}

func (gr *GLRenderer) setProjection() {
	gr.shaders[gr.currentShader].UniformMatrix4fv("projection",
		mgl32.Perspective(64.2, float32(gr.fbWidth)/float32(gr.fbHeight), 1.0, 8000.0),
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
	gr.SetLight(s.sector.LightLevel())
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

func (gr *GLRenderer) DrawThings(things []game.Thingable) {
	gr.shaders[gr.currentShader].Uniform1i("draw_phase", 1)

	for _, t := range things {
		if !t.IsShown() {
			continue
		}
		f := t.NextFrame()
		a, flipped := t.CalcAngle(gr.camera.position)
		img := gr.textures.Get(t.SpriteName()+string(f), a)
		gr.SetLight(t.GetSector().LightLevel())
		gr.shaders[gr.currentShader].Uniform3f("billboard_pos", mgl32.Vec3{
			-t.Position()[0],
			t.Height() + (float32(img.image.Height()) / 2) + 4,
			t.Position()[1],
		})

		gr.shaders[gr.currentShader].Uniform1i("billboard_flipped", flipped)
		gr.shaders[gr.currentShader].Uniform2f("billboard_size", mgl32.Vec2{
			float32(img.image.Width()) / 120,
			float32(img.image.Height()) / 100,
		})
		gr.spriter.Draw(gl.TRIANGLES, img)
	}
	gr.shaders[gr.currentShader].Uniform1i("draw_phase", 0)
}

var (
	currX = float32(1)
	currY = float32(1)
	angle float64
)

func (gr *GLRenderer) setUpHudShader(aspect float32) {

	ortho := mgl32.Ortho2D(640*aspect, 0, 0, 640)

	gl.Disable(gl.DEPTH_TEST) // Disable the Depth-testing
	gr.shaders[gr.currentShader].UniformMatrix4fv("ortho", ortho)
	gr.shaders[gr.currentShader].Uniform1i("draw_phase", 2)
}

func (gr *GLRenderer) drawHudImage(sprite string, pos mgl32.Vec3, offsetX, offsetY float32, scaleFactor float32) {
	img := gr.textures.Get(sprite, 0)
	gr.shaders[gr.currentShader].Uniform2f(
		"billboard_size",
		mgl32.Vec2{
			float32(img.image.Width()) / 40 * scaleFactor,
			float32(img.image.Height()) / 40 * scaleFactor,
		},
	)

	pos[1] += float32(-img.image.Top()) + offsetY
	pos[0] += +offsetX
	gr.shaders[gr.currentShader].Uniform3f("billboard_pos", pos)
	gr.spriter.Draw(gl.TRIANGLES, img)
}

// DrawHUD draws the game hud
func (gr *GLRenderer) DrawHUD(player *game.Player) {
	aspect := float32(gr.fbWidth) / float32(gr.fbHeight)

	gr.setUpHudShader(aspect)

	w := player.Weapon()
	frame, fire := w.NextFrames(float32((1000 / 30)))

	pos := mgl32.Vec3{(640 * aspect) / 2, 0, 0}
	scaleFactor := float32(1)

	gr.drawHudImage(w.Sprite+string(frame), pos, +w.Offset()[0], -20-w.Offset()[1], scaleFactor)

	if fire != 255 {
		gr.drawHudImage(w.FireSprite+string(fire), pos, w.FireOffset.X+w.Offset()[0], w.FireOffset.Y+w.Offset()[1], scaleFactor)
	}
	gr.shaders[gr.currentShader].Uniform1i("draw_phase", 0)
}

// DrawHUdElement draws a single element to the HUD
func (gr *GLRenderer) DrawHUdElement(name string, xpos, ypos float32, scaleFactor float32) {
	aspect := float32(gr.fbWidth) / float32(gr.fbHeight)

	gr.setUpHudShader(aspect)

	gr.drawHudImage(name, mgl32.Vec3{xpos, ypos, 0}, 0, 0, scaleFactor)
	gr.shaders[gr.currentShader].Uniform1i("draw_phase", 0)
}

func (gr *GLRenderer) SetViewPort(fbWidth, fbHeight int) {
	gr.fbWidth = fbWidth
	gr.fbHeight = fbHeight
	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
}

func (gr *GLRenderer) RenderNewFrame(frameTime float32) {
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gr.shaders[gr.currentShader].Use()
	gr.setProjection()
	gr.setView()
	gr.setModel()
}