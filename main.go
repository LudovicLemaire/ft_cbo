package main

import (
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 500
)

var (
	/*
		points = []float32{
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 0.0, 0.0,
		}
	*/
	points                       = []float32{}
	_oldMousePosX, _oldMousePosY float64
)

type Cbo struct {
	pos mgl32.Vec3
	rot mgl32.Vec3
}
type ColorTest struct {
	r float32
	g float32
	b float32
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShaderSource := _getShaderSource("Shaders/qd.vertex.glsl")
	fragmentShaderSource := _getShaderSource("Shaders/qd.fragment.glsl")

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "woop woop ft_cbo", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func getDistance(x, y, z float64) float64 {
	return math.Sqrt(math.Pow(x-0.0, 2) + math.Pow(y-0, 2) + math.Pow(z-0, 2))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	NoiseInitPermtables(42)

	var angleIncrement float64 = 10.1664073846 // TAU * golden ratio
	var totalPoints int = 5000
	var s float32 = 50 // scale

	var seaLevel float64 = 0.3

	for i := 0; i <= totalPoints; i++ {
		var t float64 = float64(i) / float64(totalPoints)
		var angle1 float64 = math.Acos(1 - 2*float64(t))
		var angle2 float64 = angleIncrement * float64(i)

		var x float32 = float32(math.Sin(angle1) * math.Cos(angle2))
		var y float32 = float32(math.Sin(angle1) * math.Sin(angle2))
		var z float32 = float32(math.Cos(angle1))
		var r float32 = float32(math.Pow(rand.Float64()/2.5, 1.5))
		var g float32 = rand.Float32()
		var b float32 = rand.Float32()

		simplexNoiseBiome := float32(Noise3dSimplex(float64(x+250), float64(y+250), float64(z+250), 1))
		simplexNoise := Noise3dSimplex(float64(x+250), float64(y+250), float64(z+250), 0)
		if simplexNoise > seaLevel {
			x *= (1 + float32(simplexNoise))
			y *= (1 + float32(simplexNoise))
			z *= (1 + float32(simplexNoise))
		} else {
			x *= (1 + float32(0.3))
			y *= (1 + float32(0.3))
			z *= (1 + float32(0.3))
		}

		distanceFromCenter := getDistance(float64(x), float64(y), float64(z))
		if distanceFromCenter < seaLevel+1 {
			r = simplexNoiseBiome
			g = ((rand.Float32()*25)+62.0)/255.0 + simplexNoiseBiome
			b = ((rand.Float32()*25)+120.0)/255.0 + simplexNoiseBiome
		} else {
			r = simplexNoiseBiome
			g = ((rand.Float32() * 25) + 154.0) / 255.0
			b = ((rand.Float32() * 25) + 23.0) / 255.0
		}
		//r = 0.5
		//g = 0.5
		//b = 0.5

		//fmt.Println(mgl32.Vec3{float32(x), float32(y), float32(z)})
		points = append(points, (x*s)-2000, (y * s), (z * s), r, g, b)
	}
	points = append(points, (0*s)-2000, (0 * s), (0 * s), 1, 0, 0)

	/*
		points = append(points, (1*s)-2000, (1 * s), (1 * s), 0, 1, 0)
		points = append(points, (-1*s)-2000, (-1 * s), (-1 * s), 0, 1, 0)

		points = append(points, (1*s)-2000, (1 * s), (-1 * s), 0, 1, 0)
		points = append(points, (1*s)-2000, (-1 * s), (-1 * s), 0, 1, 0)
		points = append(points, (1*s)-2000, (-1 * s), (1 * s), 0, 1, 0)

		points = append(points, (-1*s)-2000, (-1 * s), (1 * s), 0, 1, 0)
		points = append(points, (-1*s)-2000, (1 * s), (-1 * s), 0, 1, 0)
		points = append(points, (-1*s)-2000, (1 * s), (1 * s), 0, 1, 0)
	*/

	var cbo Cbo
	var colorTest ColorTest
	colorTest.r, colorTest.g, colorTest.b = 1.0, 1.0, 1.0
	cbo.pos[0], cbo.pos[1], cbo.pos[2] = 2000.0, 0.0, -500.0

	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	redUniform := gl.GetUniformLocation(program, gl.Str("red\x00"))

	cameraCStr, free := gl.Strs("camera")
	defer free()
	cameraId := gl.GetUniformLocation(program, *cameraCStr)

	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	gl.PointSize(5)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.POINT)

	vao := makeVao(points)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Disable(gl.CULL_FACE)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	for !window.ShouldClose() {
		draw(vao, window, program, &cbo, &colorTest, cameraId)
		gl.Uniform1f(redUniform, colorTest.r)
	}
}

func setCamera(cameraId int32, cbo *Cbo) {
	camera := mgl32.HomogRotate3D(cbo.rot.X(), mgl32.Vec3{1, 0, 0})
	camera = camera.Mul4(mgl32.HomogRotate3D(cbo.rot.Y(), mgl32.Vec3{0, 1, 0}))
	camera = camera.Mul4(mgl32.Translate3D(cbo.pos.X(), cbo.pos.Y(), cbo.pos.Z()))

	projection := mgl32.Perspective(mgl32.DegToRad(80.0), float32(width)/float32(height), 0.1, 5000)
	view := projection.Mul4(camera)

	gl.UniformMatrix4fv(cameraId, 1, false, &view[0])
}

func draw(vao uint32, window *glfw.Window, program uint32, cbo *Cbo, colorTest *ColorTest, cameraId int32) {
	EventsKeyboard(cbo, colorTest)
	EventsMouse(cbo)
	setCamera(cameraId, cbo)
	//camera := GetCameraMatrix(cbo)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.POINTS, 0, int32(len(points)/3))

	glfw.PollEvents()
	window.SwapBuffers()
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return vao
}
