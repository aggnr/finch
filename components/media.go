package components

import (
	"image"
	"image/color"
)

// Image represents an image element in the UI
type Image struct {
	*Node
	source    image.Image
	srcPath   string
	fitMethod ImageFitMethod
}

// ImageFitMethod defines how an image should be sized to fit its container
type ImageFitMethod int

const (
	ImageFitContain ImageFitMethod = iota // Maintain aspect ratio, fit within bounds
	ImageFitCover                         // Maintain aspect ratio, cover entire bounds (may crop)
	ImageFitFill                          // Stretch to fill bounds (may distort)
)

// NewImage creates a new image element
func NewImage(id string) *Image {
	return &Image{
		Node:      NewNode(id),
		source:    nil,
		srcPath:   "",
		fitMethod: ImageFitContain,
	}
}

// SetSource sets the image source
func (i *Image) SetSource(img image.Image) {
	i.source = img
}

// SetSourcePath sets the path to the image source
func (i *Image) SetSourcePath(path string) {
	i.srcPath = path
	// In a real implementation, this would load the image from the path
}

// SetFitMethod sets how the image should fit within its bounds
func (i *Image) SetFitMethod(method ImageFitMethod) {
	i.fitMethod = method
}

// Draw draws the image
func (i *Image) Draw(surface DrawSurface) {
	if !i.IsVisible() || i.source == nil {
		return
	}
	
	bounds := i.ComputedBounds()
	
	// Draw the image
	surface.DrawImage(i.source, bounds.X, bounds.Y, bounds.Width, bounds.Height, i.fitMethod)
	
	// Draw children (if any)
	for _, child := range i.Children() {
		child.Draw(surface)
	}
}

// HandleMouseDown handles mouse down events
func (i *Image) HandleMouseDown(x, y int) bool {
	// Image doesn't handle mouse events directly, but we check children
	for j := len(i.Children()) - 1; j >= 0; j-- {
		child := i.Children()[j]
		if child.HandleMouseDown(x, y) {
			return true
		}
	}
	return false
}

// Video represents a video element in the UI
type Video struct {
	*Node
	source     string
	isPlaying  bool
	onPlay     func()
	onPause    func()
	onComplete func()
	volume     float64
}

// NewVideo creates a new video element
func NewVideo(id string) *Video {
	return &Video{
		Node:      NewNode(id),
		source:    "",
		isPlaying: false,
		onPlay:    nil,
		onPause:   nil,
		volume:    1.0,
	}
}

// SetSource sets the video source
func (v *Video) SetSource(source string) {
	v.source = source
}

// Play starts playing the video
func (v *Video) Play() {
	if !v.isPlaying {
		v.isPlaying = true
		if v.onPlay != nil {
			v.onPlay()
		}
	}
}

// Pause pauses the video
func (v *Video) Pause() {
	if v.isPlaying {
		v.isPlaying = false
		if v.onPause != nil {
			v.onPause()
		}
	}
}

// IsPlaying returns whether the video is playing
func (v *Video) IsPlaying() bool {
	return v.isPlaying
}

// SetVolume sets the video volume (0.0 to 1.0)
func (v *Video) SetVolume(volume float64) {
	if volume < 0.0 {
		v.volume = 0.0
	} else if volume > 1.0 {
		v.volume = 1.0
	} else {
		v.volume = volume
	}
}

// SetOnPlay sets the handler for when the video starts playing
func (v *Video) SetOnPlay(handler func()) {
	v.onPlay = handler
}

// SetOnPause sets the handler for when the video is paused
func (v *Video) SetOnPause(handler func()) {
	v.onPause = handler
}

// SetOnComplete sets the handler for when the video finishes playing
func (v *Video) SetOnComplete(handler func()) {
	v.onComplete = handler
}

// Draw draws the video
func (v *Video) Draw(surface DrawSurface) {
	if !v.IsVisible() {
		return
	}
	
	bounds := v.ComputedBounds()
	
	// Draw video background (representing video content)
	surface.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{0, 0, 0, 255})
	
	// Draw play/pause indicator
	if v.isPlaying {
		// Draw pause icon
		pauseX1 := bounds.X + bounds.Width / 2 - 10
		pauseX2 := bounds.X + bounds.Width / 2 + 2
		pauseY := bounds.Y + bounds.Height / 2 - 10
		pauseHeight := 20
		
		surface.FillRect(pauseX1, pauseY, 6, pauseHeight, color.RGBA{255, 255, 255, 200})
		surface.FillRect(pauseX2, pauseY, 6, pauseHeight, color.RGBA{255, 255, 255, 200})
	} else {
		// Draw play icon (triangle)
		playX := bounds.X + bounds.Width / 2 - 5
		playY := bounds.Y + bounds.Height / 2 - 10
		
		surface.DrawLine(playX, playY, playX, playY + 20, color.RGBA{255, 255, 255, 200})
		surface.DrawLine(playX, playY, playX + 15, playY + 10, color.RGBA{255, 255, 255, 200})
		surface.DrawLine(playX, playY + 20, playX + 15, playY + 10, color.RGBA{255, 255, 255, 200})
	}
	
	// Draw children (if any)
	for _, child := range v.Children() {
		child.Draw(surface)
	}
}

// HandleMouseDown handles mouse down events
func (v *Video) HandleMouseDown(x, y int) bool {
	bounds := v.ComputedBounds()
	if PointInRect(Point{x, y}, bounds) {
		// Toggle play/pause on click
		if v.isPlaying {
			v.Pause()
		} else {
			v.Play()
		}
		return true
	}
	
	// Check children
	for i := len(v.Children()) - 1; i >= 0; i-- {
		child := v.Children()[i]
		if child.HandleMouseDown(x, y) {
			return true
		}
	}
	
	return false
}

// Audio represents an audio element in the UI
type Audio struct {
	*Node
	source     string
	isPlaying  bool
	onPlay     func()
	onPause    func()
	onComplete func()
	volume     float64
}

// NewAudio creates a new audio element
func NewAudio(id string) *Audio {
	return &Audio{
		Node:      NewNode(id),
		source:    "",
		isPlaying: false,
		onPlay:    nil,
		onPause:   nil,
		volume:    1.0,
	}
}

// SetSource sets the audio source
func (a *Audio) SetSource(source string) {
	a.source = source
}

// Play starts playing the audio
func (a *Audio) Play() {
	if !a.isPlaying {
		a.isPlaying = true
		if a.onPlay != nil {
			a.onPlay()
		}
	}
}

// Pause pauses the audio
func (a *Audio) Pause() {
	if a.isPlaying {
		a.isPlaying = false
		if a.onPause != nil {
			a.onPause()
		}
	}
}

// IsPlaying returns whether the audio is playing
func (a *Audio) IsPlaying() bool {
	return a.isPlaying
}

// SetVolume sets the audio volume (0.0 to 1.0)
func (a *Audio) SetVolume(volume float64) {
	if volume < 0.0 {
		a.volume = 0.0
	} else if volume > 1.0 {
		a.volume = 1.0
	} else {
		a.volume = volume
	}
}

// SetOnPlay sets the handler for when the audio starts playing
func (a *Audio) SetOnPlay(handler func()) {
	a.onPlay = handler
}

// SetOnPause sets the handler for when the audio is paused
func (a *Audio) SetOnPause(handler func()) {
	a.onPause = handler
}

// SetOnComplete sets the handler for when the audio finishes playing
func (a *Audio) SetOnComplete(handler func()) {
	a.onComplete = handler
}

// Draw draws the audio control
func (a *Audio) Draw(surface DrawSurface) {
	if !a.IsVisible() {
		return
	}
	
	bounds := a.ComputedBounds()
	
	// Draw audio control background
	surface.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{80, 80, 80, 255})
	
	// Draw play/pause button
	buttonX := bounds.X + 5
	buttonY := bounds.Y + 5
	buttonSize := bounds.Height - 10
	
	surface.FillRect(buttonX, buttonY, buttonSize, buttonSize, color.RGBA{120, 120, 120, 255})
	
	if a.isPlaying {
		// Draw pause icon
		pauseX1 := buttonX + buttonSize / 3
		pauseX2 := buttonX + buttonSize * 2 / 3 - 2
		pauseY := buttonY + buttonSize / 4
		pauseHeight := buttonSize / 2
		
		surface.FillRect(pauseX1, pauseY, 2, pauseHeight, color.RGBA{255, 255, 255, 255})
		surface.FillRect(pauseX2, pauseY, 2, pauseHeight, color.RGBA{255, 255, 255, 255})
	} else {
		// Draw play icon (triangle)
		playX := buttonX + buttonSize / 3
		playY := buttonY + buttonSize / 4
		
		surface.DrawLine(playX, playY, playX, playY + buttonSize / 2, color.RGBA{255, 255, 255, 255})
		surface.DrawLine(playX, playY, playX + buttonSize / 3, playY + buttonSize / 4, color.RGBA{255, 255, 255, 255})
		surface.DrawLine(playX, playY + buttonSize / 2, playX + buttonSize / 3, playY + buttonSize / 4, color.RGBA{255, 255, 255, 255})
	}
	
	// Draw volume slider background
	sliderX := buttonX + buttonSize + 10
	sliderY := bounds.Y + bounds.Height / 2
	sliderWidth := bounds.Width - buttonSize - 20
	
	surface.DrawLine(sliderX, sliderY, sliderX + sliderWidth, sliderY, color.RGBA{150, 150, 150, 255})
	
	// Draw volume slider position
	volumePos := int(a.volume * float64(sliderWidth))
	surface.FillRect(sliderX, sliderY - 3, volumePos, 6, color.RGBA{200, 200, 200, 255})
	
	// Draw children (if any)
	for _, child := range a.Children() {
		child.Draw(surface)
	}
}

// HandleMouseDown handles mouse down events
func (a *Audio) HandleMouseDown(x, y int) bool {
	bounds := a.ComputedBounds()
	if !PointInRect(Point{x, y}, bounds) {
		return false
	}
	
	// Check if click is on play/pause button
	buttonX := bounds.X + 5
	buttonY := bounds.Y + 5
	buttonSize := bounds.Height - 10
	
	buttonRect := Rect{buttonX, buttonY, buttonSize, buttonSize}
	if PointInRect(Point{x, y}, buttonRect) {
		// Toggle play/pause
		if a.isPlaying {
			a.Pause()
		} else {
			a.Play()
		}
		return true
	}
	
	// Check if click is on volume slider
	sliderX := buttonX + buttonSize + 10
	sliderY := bounds.Y + bounds.Height / 2
	sliderWidth := bounds.Width - buttonSize - 20
	
	if x >= sliderX && x <= sliderX + sliderWidth && y >= sliderY - 5 && y <= sliderY + 5 {
		// Set volume based on click position
		a.SetVolume(float64(x - sliderX) / float64(sliderWidth))
		return true
	}
	
	// Check children
	for i := len(a.Children()) - 1; i >= 0; i-- {
		child := a.Children()[i]
		if child.HandleMouseDown(x, y) {
			return true
		}
	}
	
	return true
} 