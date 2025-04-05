package components

import (
	"image/color"
	"image"
)

// ScreenWidth and ScreenHeight define the default screen dimensions
const (
	ScreenWidth  = 1200
	ScreenHeight = 800
)

// InputType defines the type of input event
type InputType int

const (
	InputTypeMouseDown InputType = iota
	InputTypeMouseUp
	InputTypeMouseMove
	InputTypeKeyDown
	InputTypeKeyUp
	InputTypeChar
)

// Key represents keyboard keys
type Key int

const (
	KeyUnknown Key = iota
	KeyEscape
	KeyEnter
	KeyBackspace
	KeyTab
	KeySpace
	// Add more keys as needed
)

// InputEvent represents an input event
type InputEvent struct {
	Type      InputType
	X         int
	Y         int
	Key       Key
	Char      rune
	ShiftDown bool
	CtrlDown  bool
	AltDown   bool
}

// Element is the interface for all UI elements
type Element interface {
	// Common methods for all UI elements
	ID() string
	SetID(id string)
	Bounds() Rect
	SetBounds(bounds Rect)
	Parent() Element
	SetParent(parent Element)
	Children() []Element
	AddChild(child Element)
	RemoveChild(child Element)
	
	// Input handling
	HandleMouseDown(x, y int) bool
	HandleMouseUp(x, y int) bool
	HandleMouseMove(x, y int) bool
	
	// Rendering
	Draw(surface DrawSurface)
	
	// State updates
	Update()
}

// Rect represents a rectangle with position and dimensions
type Rect struct {
	X, Y, Width, Height int
}

// Point represents a coordinate point
type Point struct {
	X, Y int
}

// DrawSurface is the interface for drawing to the screen
type DrawSurface interface {
	Clear(color color.RGBA)
	DrawText(text string, x, y int, color color.RGBA, fontSize int)
	DrawRect(x, y, width, height int, color color.RGBA)
	FillRect(x, y, width, height int, color color.RGBA)
	DrawLine(x1, y1, x2, y2 int, color color.RGBA)
	FillCircle(x, y, radius int, color color.RGBA)
	DrawCircle(x, y, radius int, color color.RGBA)
	SetClipRect(x, y, width, height int)
	ResetClipRect()
	DrawImage(img image.Image, x, y, width, height int, fitMethod ImageFitMethod)
}

// PointInRect checks if a point is inside a rectangle
func PointInRect(p Point, r Rect) bool {
	return p.X >= r.X && p.X < r.X+r.Width && p.Y >= r.Y && p.Y < r.Y+r.Height
} 