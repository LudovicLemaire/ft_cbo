package main

import (
	"os"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

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
	if glfw.GetCurrentContext().GetKey(glfw.KeyR) == 1 {
		move = move.Add(mgl32.Vec3{0, -0.10 * multiplier, 0})
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyF) == 1 {
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

func EventsKeyboard(cbo *Cbo, colorTest *ColorTest) {
	var multiplier float32 = 1

	if glfw.GetCurrentContext().GetKey(glfw.KeyEscape) == 1 {
		os.Exit(1)
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyLeftShift) == 1 {
		multiplier *= 20.0
	}
	if glfw.GetCurrentContext().GetKey(glfw.KeyLeftControl) == 1 {
		multiplier *= 10.0
	}

	if glfw.GetCurrentContext().GetKey(glfw.KeySpace) == 1 {
		colorTest.r = 0.0
	} else {
		colorTest.r = 1.0
	}

	_keyboardTranslate(cbo, multiplier)
	_keyboardRotation(cbo)
}
