package components

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// EbitenRenderer implements the DrawSurface interface using Ebiten
type EbitenRenderer struct {
	target   *ebiten.Image
	font     font.Face
	clipRect Rect
}

// NewEbitenRenderer creates a new Ebiten-based renderer
func NewEbitenRenderer(target *ebiten.Image) *EbitenRenderer {
	return &EbitenRenderer{
		target:   target,
		font:     basicfont.Face7x13,
		clipRect: Rect{X: 0, Y: 0, Width: ScreenWidth, Height: ScreenHeight},
	}
}

// Clear clears the screen with the specified color
func (r *EbitenRenderer) Clear(color color.RGBA) {
	r.target.Fill(color)
}

// DrawText draws text at the specified position
func (r *EbitenRenderer) DrawText(txt string, x, y int, clr color.RGBA, fontSize int) {
	text.Draw(r.target, txt, r.font, x, y+13, clr) // +13 for font baseline
}

// DrawRect draws a rectangle with the specified position and dimensions
func (r *EbitenRenderer) DrawRect(x, y, width, height int, clr color.RGBA) {
	// Apply clip rect
	if !r.isVisibleInClipRect(x, y, width, height) {
		return
	}

	// Draw rectangle
	x1, y1 := float32(x), float32(y)
	x2, y2 := float32(x+width), float32(y+height)

	vector.StrokeLine(r.target, x1, y1, x2, y1, 1, clr, false)
	vector.StrokeLine(r.target, x2, y1, x2, y2, 1, clr, false)
	vector.StrokeLine(r.target, x2, y2, x1, y2, 1, clr, false)
	vector.StrokeLine(r.target, x1, y2, x1, y1, 1, clr, false)
}

// FillRect fills a rectangle with the specified position, dimensions, and color
func (r *EbitenRenderer) FillRect(x, y, width, height int, clr color.RGBA) {
	// Apply clip rect
	if !r.isVisibleInClipRect(x, y, width, height) {
		return
	}

	// Fill rectangle
	vector.DrawFilledRect(r.target, float32(x), float32(y), float32(width), float32(height), clr, false)
}

// DrawLine draws a line from (x1, y1) to (x2, y2)
func (r *EbitenRenderer) DrawLine(x1, y1, x2, y2 int, clr color.RGBA) {
	// Draw line
	vector.StrokeLine(r.target, float32(x1), float32(y1), float32(x2), float32(y2), 1, clr, false)
}

// FillCircle fills a circle with the specified center, radius, and color
func (r *EbitenRenderer) FillCircle(x, y, radius int, clr color.RGBA) {
	vector.DrawFilledCircle(r.target, float32(x), float32(y), float32(radius), clr, false)
}

// DrawCircle draws a circle outline with the specified center, radius, and color
func (r *EbitenRenderer) DrawCircle(x, y, radius int, clr color.RGBA) {
	const segments = 36
	r.drawCircleSegments(x, y, radius, segments, clr)
}

// drawCircleSegments draws a circle using line segments
func (r *EbitenRenderer) drawCircleSegments(x, y, radius, segments int, clr color.RGBA) {
	for i := 0; i < segments; i++ {
		angle1 := 2 * math.Pi * float64(i) / float64(segments)
		angle2 := 2 * math.Pi * float64(i+1) / float64(segments)

		x1 := x + int(math.Cos(angle1)*float64(radius))
		y1 := y + int(math.Sin(angle1)*float64(radius))
		x2 := x + int(math.Cos(angle2)*float64(radius))
		y2 := y + int(math.Sin(angle2)*float64(radius))

		r.DrawLine(x1, y1, x2, y2, clr)
	}
}

// isVisibleInClipRect checks if a rectangle is visible within the clip rect
func (r *EbitenRenderer) isVisibleInClipRect(x, y, width, height int) bool {
	if x+width < r.clipRect.X || x > r.clipRect.X+r.clipRect.Width ||
		y+height < r.clipRect.Y || y > r.clipRect.Y+r.clipRect.Height {
		return false
	}
	return true
}

// SetClipRect sets the clipping rectangle
func (r *EbitenRenderer) SetClipRect(x, y, width, height int) {
	r.clipRect = Rect{X: x, Y: y, Width: width, Height: height}
}

// ResetClipRect resets the clipping rectangle to the full screen
func (r *EbitenRenderer) ResetClipRect() {
	r.clipRect = Rect{X: 0, Y: 0, Width: ScreenWidth, Height: ScreenHeight}
}

// DrawImage draws an image with the specified fit method
func (r *EbitenRenderer) DrawImage(img image.Image, x, y, width, height int, fitMethod ImageFitMethod) {
    // Implementation needed for EbitenRenderer
    // For now, just draw a placeholder
    r.FillRect(x, y, width, height, color.RGBA{200, 200, 200, 255})
    r.DrawRect(x, y, width, height, color.RGBA{150, 150, 150, 255})
}

// EbitenDrawSurface implements DrawSurface using Ebiten
type EbitenDrawSurface struct {
	target *ebiten.Image
	font   font.Face
}

// NewEbitenDrawSurface creates a new Ebiten-based draw surface
func NewEbitenDrawSurface(target *ebiten.Image) *EbitenDrawSurface {
	return &EbitenDrawSurface{
		target: target,
		font:   basicfont.Face7x13, // Default font
	}
}

// Clear clears the screen with the specified color
func (e *EbitenDrawSurface) Clear(color color.RGBA) {
	e.target.Fill(color)
}

// FillRect fills a rectangle with the specified color
func (e *EbitenDrawSurface) FillRect(x, y, width, height int, color color.RGBA) {
	vector.DrawFilledRect(e.target, float32(x), float32(y), float32(width), float32(height), color, false)
}

// DrawRect draws a rectangle outline with the specified color
func (e *EbitenDrawSurface) DrawRect(x, y, width, height int, color color.RGBA) {
	// Top line
	vector.StrokeLine(e.target, float32(x), float32(y), float32(x+width), float32(y), 1, color, false)
	// Right line
	vector.StrokeLine(e.target, float32(x+width), float32(y), float32(x+width), float32(y+height), 1, color, false)
	// Bottom line
	vector.StrokeLine(e.target, float32(x+width), float32(y+height), float32(x), float32(y+height), 1, color, false)
	// Left line
	vector.StrokeLine(e.target, float32(x), float32(y+height), float32(x), float32(y), 1, color, false)
}

// DrawLine draws a line between two points
func (e *EbitenDrawSurface) DrawLine(x1, y1, x2, y2 int, color color.RGBA) {
	vector.StrokeLine(e.target, float32(x1), float32(y1), float32(x2), float32(y2), 1, color, false)
}

// DrawText draws text at the specified position
func (e *EbitenDrawSurface) DrawText(txt string, x, y int, color color.RGBA, fontSize int) {
	// In a real implementation, you'd use font caching and handle size changes
	// For this demo, just use the basic font
	text.Draw(e.target, txt, e.font, x, y+13, color) // +13 for font baseline
}

// FillCircle fills a circle with the specified center, radius, and color
func (e *EbitenDrawSurface) FillCircle(x, y, radius int, clr color.RGBA) {
	vector.DrawFilledCircle(e.target, float32(x), float32(y), float32(radius), clr, false)
}

// DrawCircle draws a circle outline with the specified center, radius, and color
func (e *EbitenDrawSurface) DrawCircle(x, y, radius int, clr color.RGBA) {
	const segments = 36
	e.drawCircleSegments(x, y, radius, segments, clr)
}

// drawCircleSegments draws a circle using line segments
func (e *EbitenDrawSurface) drawCircleSegments(x, y, radius, segments int, clr color.RGBA) {
	for i := 0; i < segments; i++ {
		angle1 := 2 * math.Pi * float64(i) / float64(segments)
		angle2 := 2 * math.Pi * float64(i+1) / float64(segments)

		x1 := x + int(math.Cos(angle1)*float64(radius))
		y1 := y + int(math.Sin(angle1)*float64(radius))
		x2 := x + int(math.Cos(angle2)*float64(radius))
		y2 := y + int(math.Sin(angle2)*float64(radius))

		e.DrawLine(x1, y1, x2, y2, clr)
	}
}

// SetClipRect sets the clipping rectangle
func (e *EbitenDrawSurface) SetClipRect(x, y, width, height int) {
	// Implement clipping if needed
}

// ResetClipRect resets the clipping rectangle to the full screen
func (e *EbitenDrawSurface) ResetClipRect() {
	// Implement clipping if needed
}

// DrawImage draws an image with the specified fit method
func (e *EbitenDrawSurface) DrawImage(img image.Image, x, y, width, height int, fitMethod ImageFitMethod) {
	if img == nil {
		// Draw placeholder if image is nil
		e.FillRect(x, y, width, height, color.RGBA{200, 200, 200, 255})
		e.DrawRect(x, y, width, height, color.RGBA{150, 150, 150, 255})
		return
	}
	
	// Convert to Ebiten image if needed
	var eImg *ebiten.Image
	if ebi, ok := img.(*ebiten.Image); ok {
		eImg = ebi
	} else {
		// In a real app, you'd convert the image.Image to an ebiten.Image here
		// For now, just draw a placeholder
		e.FillRect(x, y, width, height, color.RGBA{200, 200, 200, 255})
		e.DrawRect(x, y, width, height, color.RGBA{150, 150, 150, 255})
		return
	}
	
	// Get image dimensions
	imgWidth, imgHeight := eImg.Size()
	
	// Calculate scaling based on fit method
	var scale float64
	var offsetX, offsetY int
	
	switch fitMethod {
	case ImageFitContain:
		// Scale to fit within bounds while maintaining aspect ratio
		scaleX := float64(width) / float64(imgWidth)
		scaleY := float64(height) / float64(imgHeight)
		scale = math.Min(scaleX, scaleY)
		
		// Center the image
		offsetX = x + (width-int(float64(imgWidth)*scale))/2
		offsetY = y + (height-int(float64(imgHeight)*scale))/2
		
	case ImageFitCover:
		// Scale to cover entire bounds while maintaining aspect ratio
		scaleX := float64(width) / float64(imgWidth)
		scaleY := float64(height) / float64(imgHeight)
		scale = math.Max(scaleX, scaleY)
		
		// Center the image
		offsetX = x + (width-int(float64(imgWidth)*scale))/2
		offsetY = y + (height-int(float64(imgHeight)*scale))/2
		
	case ImageFitFill:
		// Scale to fill bounds, potentially distorting the image
		scaleX := float64(width) / float64(imgWidth)
		scaleY := float64(height) / float64(imgHeight)
		
		// Draw the image with different scales for X and Y
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(float64(x), float64(y))
		e.target.DrawImage(eImg, op)
		return // Early return because we handled this case differently
	}
	
	// Draw the image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(offsetX), float64(offsetY))
	e.target.DrawImage(eImg, op)
} 