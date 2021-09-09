package main

import (
	"os"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Keys struct {
	escape string
	t      string
	v      string
	shift  string
	kp1    string
	kp2    string
	kp3    string
	kp4    string
	kp5    string
	kp6    string
}

func initKeys(keys *Keys) {
	keys.escape = "null"
	keys.t = "null"
	keys.v = "null"
	keys.shift = "null"
	keys.kp1 = "null"
	keys.kp2 = "null"
	keys.kp3 = "null"
	keys.kp4 = "null"
	keys.kp5 = "null"
	keys.kp6 = "null"
}

func initGameValues(gameValues *GameValues) {
	gameValues.speed = 1
}

func _keyboardRotation(cbo *Cbo) {
	if glfw.GetCurrentContext().GetKey(glfw.KeyLeft) == 1 {
		cbo.rot = cbo.rot.Sub(mgl32.Vec3{0, 0.05, 0})
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyRight) == 1 {
		cbo.rot = cbo.rot.Add(mgl32.Vec3{0, 0.05, 0})
	}

	if glfw.GetCurrentContext().GetKey(glfw.KeyUp) == 1 {
		cbo.rot = cbo.rot.Sub(mgl32.Vec3{0.05, 0, 0})
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyDown) == 1 {
		cbo.rot = cbo.rot.Add(mgl32.Vec3{0.05, 0, 0})
	}
}

func _keyboardTranslate(cbo *Cbo, multiplier float32) {
	move := mgl32.Vec3{0, 0, 0}

	if glfw.GetCurrentContext().GetKey(glfw.KeyW) == 1 {
		move = move.Add(mgl32.Vec3{0, 0, 0.10 * multiplier})
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyS) == 1 {
		move = move.Add(mgl32.Vec3{0, 0, -0.10 * multiplier})
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyA) == 1 {
		move = move.Add(mgl32.Vec3{0.10 * multiplier, 0, 0})

	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyD) == 1 {
		move = move.Add(mgl32.Vec3{-0.10 * multiplier, 0, 0})
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeySpace) == 1 {
		move = move.Add(mgl32.Vec3{0, -0.10 * multiplier, 0})
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyLeftControl) == 1 {
		move = move.Add(mgl32.Vec3{0, 0.10 * multiplier, 0})
	}

	if move[0] != 0 || move[1] != 0 || move[2] != 0 {
		rotationMat := mgl32.HomogRotate3D(cbo.rot.Y(), mgl32.Vec3{0, -1, 0})
		rotationMat = rotationMat.Mul4((mgl32.HomogRotate3D(cbo.rot.X(), mgl32.Vec3{-1, 0, 0})))
		cbo.pos = cbo.pos.Add(rotationMat.Mul4x1(move.Vec4(1)).Vec3())
	}
}

func EventsMouse(cbo *Cbo) {
	posX, posY := glfw.GetCurrentContext().GetCursorPos()

	if _oldMousePosX == 0 {
		_oldMousePosX = posX
	}
	if _oldMousePosY == 0 {
		_oldMousePosY = posY
	}

	cbo.rot = cbo.rot.Add(mgl32.Vec3{
		-float32((_oldMousePosY - posY) * 0.001),
		-float32((_oldMousePosX - posX) * 0.001),
		0,
	})
	_oldMousePosX = posX
	_oldMousePosY = posY
}

func getKeyStatus(key glfw.Key, status string) string {
	if glfw.GetCurrentContext().GetKey(key) == 1 && !(status == "active" || status == "hold") {
		status = "active"
	} else if glfw.GetCurrentContext().GetKey(key) == 1 && (status == "active" || status == "hold") {
		status = "hold"
	} else if glfw.GetCurrentContext().GetKey(key) == 0 && (status == "active" || status == "hold") {
		status = "released"
	} else if glfw.GetCurrentContext().GetKey(key) == 0 && (status == "released" || status == "null") {
		status = "null"
	}
	return status
}

func EventsKeyboard(cbo *Cbo, colorTest *ColorRGB, k *Keys, gameValues *GameValues) {
	k.escape = getKeyStatus(glfw.KeyEscape, k.escape)
	k.t = getKeyStatus(glfw.KeyT, k.t)
	k.v = getKeyStatus(glfw.KeyV, k.v)
	k.shift = getKeyStatus(glfw.KeyLeftShift, k.shift)
	k.kp1 = getKeyStatus(glfw.KeyKP1, k.kp1)
	k.kp2 = getKeyStatus(glfw.KeyKP2, k.kp2)
	k.kp3 = getKeyStatus(glfw.KeyKP3, k.kp3)
	k.kp4 = getKeyStatus(glfw.KeyKP4, k.kp4)
	k.kp5 = getKeyStatus(glfw.KeyKP5, k.kp5)
	k.kp6 = getKeyStatus(glfw.KeyKP6, k.kp6)

	if k.escape == "active" {
		os.Exit(1)
	}
	if k.kp1 == "active" {
		gameValues.speed = 1.0
	}
	if k.kp2 == "active" {
		gameValues.speed = 7.5
	}
	if k.kp3 == "active" {
		gameValues.speed = 20.0
	}
	if k.kp4 == "active" {
		gameValues.speed = 100.0
	}
	if k.kp5 == "active" {
		gameValues.speed = 250.0
	}
	if k.kp6 == "active" {
		gameValues.speed = 1000.0
	}
	if k.shift == "active" {
		gameValues.speed *= 2.5
	} else if k.shift == "released" {
		gameValues.speed /= 2.5
	}

	if k.t == "active" {
		if colorTest.r == 1.0 {
			colorTest.r = 0.0
		} else {
			colorTest.r = 1.0
		}
	}

	_keyboardTranslate(cbo, gameValues.speed)
	_keyboardRotation(cbo)
}
