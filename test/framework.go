package test

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"

	"github.com/aggnr/finch/components"
)

// TestEvent represents a test event
type TestEvent struct {
	Type      string
	ElementID string
	X, Y      int
	Key       components.Key
	Char      rune
}

// TestResult represents the result of a test
type TestResult struct {
	Event     TestEvent
	Handled   bool
	Timestamp time.Time
}

// UITest represents an interactive UI test
type UITest struct {
	rootElement components.Element
	events      []TestEvent
	results     []TestResult
	surface     components.DrawSurface
}

// NewUITest creates a new UI test
func NewUITest(root components.Element) *UITest {
	return &UITest{
		rootElement: root,
		events:      make([]TestEvent, 0),
		results:     make([]TestResult, 0),
		surface:     NewMemorySurface(components.ScreenWidth, components.ScreenHeight),
	}
}

// AddClickEvent adds a mouse click event to the test sequence
func (t *UITest) AddClickEvent(elementID string, x, y int) {
	t.events = append(t.events, TestEvent{
		Type:      "click",
		ElementID: elementID,
		X:         x,
		Y:         y,
	})
}

// AddKeyEvent adds a keyboard event to the test sequence
func (t *UITest) AddKeyEvent(key components.Key) {
	t.events = append(t.events, TestEvent{
		Type: "key",
		Key:  key,
	})
}

// Run executes the test sequence
func (t *UITest) Run() {
	fmt.Println("Running UI test...")
	
	// Render the initial UI state
	t.rootElement.Draw(t.surface)
	t.SaveScreenshot("test_initial.png")
	
	// Process each event in sequence
	for i, event := range t.events {
		fmt.Printf("Processing event %d: %s\n", i+1, event.Type)
		
		var handled bool
		
		// Handle the event based on type
		switch event.Type {
		case "click":
			// Simulate mouse down
			handled = t.rootElement.HandleMouseDown(event.X, event.Y)
			
			// Render the UI after mouse down
			t.rootElement.Draw(t.surface)
			t.SaveScreenshot(fmt.Sprintf("test_event_%d_mousedown.png", i+1))
			
			// Small delay to simulate real interaction
			time.Sleep(100 * time.Millisecond)
			
			// Simulate mouse up
			handled = t.rootElement.HandleMouseUp(event.X, event.Y)
			
		case "key":
			// Keyboard events would be handled here
			handled = false
		}
		
		// Record the result
		t.results = append(t.results, TestResult{
			Event:     event,
			Handled:   handled,
			Timestamp: time.Now(),
		})
		
		// Render the UI after the event
		t.rootElement.Draw(t.surface)
		t.SaveScreenshot(fmt.Sprintf("test_event_%d.png", i+1))
		
		// Small delay to make interactive viewing possible
		time.Sleep(500 * time.Millisecond)
	}
	
	fmt.Println("Test completed.")
	t.PrintResults()
}

// Interactive runs the test with user input
func (t *UITest) Interactive() {
	fmt.Println("Starting interactive UI test...")
	fmt.Println("Press Enter to process each event...")
	
	// Render the initial UI state
	t.rootElement.Draw(t.surface)
	t.SaveScreenshot("test_interactive_initial.png")
	fmt.Println("Initial UI state saved to test_interactive_initial.png")
	
	waitForEnter()
	
	// Process each event in sequence
	for i, event := range t.events {
		fmt.Printf("Processing event %d: %s\n", i+1, event.Type)
		
		var handled bool
		
		// Handle the event based on type
		switch event.Type {
		case "click":
			// Simulate mouse down
			handled = t.rootElement.HandleMouseDown(event.X, event.Y)
			
			// Render the UI after mouse down
			t.rootElement.Draw(t.surface)
			t.SaveScreenshot(fmt.Sprintf("test_interactive_%d_mousedown.png", i+1))
			
			// Small delay to simulate real interaction
			time.Sleep(100 * time.Millisecond)
			
			// Simulate mouse up
			handled = t.rootElement.HandleMouseUp(event.X, event.Y)
			
		case "key":
			// Keyboard events would be handled here
			handled = false
		}
		
		// Record the result
		t.results = append(t.results, TestResult{
			Event:     event,
			Handled:   handled,
			Timestamp: time.Now(),
		})
		
		// Render the UI after the event
		t.rootElement.Draw(t.surface)
		t.SaveScreenshot(fmt.Sprintf("test_interactive_%d.png", i+1))
		fmt.Printf("UI state after event %d saved to test_interactive_%d.png\n", i+1, i+1)
		
		if i < len(t.events)-1 {
			waitForEnter()
		}
	}
	
	fmt.Println("Interactive test completed.")
	t.PrintResults()
}

// SaveScreenshot saves the current UI state as an image
func (t *UITest) SaveScreenshot(filename string) {
	image := t.surface.(*MemorySurface).Image()
	
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()
	
	if err := png.Encode(f, image); err != nil {
		fmt.Println("Error encoding image:", err)
	}
}

// PrintResults prints the test results
func (t *UITest) PrintResults() {
	fmt.Println("Test Results:")
	for i, result := range t.results {
		fmt.Printf("  Event %d (%s): Handled=%v\n", i+1, result.Event.Type, result.Handled)
	}
}

// Helper functions
func waitForEnter() {
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()
}

// MemorySurface is an in-memory implementation of the DrawSurface interface
type MemorySurface struct {
	img *image.RGBA
}

// NewMemorySurface creates a new memory surface
func NewMemorySurface(width, height int) *MemorySurface {
	return &MemorySurface{
		img: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

// Clear clears the surface
func (s *MemorySurface) Clear(color color.RGBA) {
	bounds := s.img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			s.img.SetRGBA(x, y, color)
		}
	}
}

// DrawText draws text on the surface
func (s *MemorySurface) DrawText(text string, x, y int, color color.RGBA, fontSize int) {
	// In a real implementation, this would use font rendering
	// For this test framework, we just draw a rectangle representing text
	width := len(text) * fontSize / 2
	height := fontSize
	
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if i == 0 || j == 0 || i == width-1 || j == height-1 {
				s.img.SetRGBA(x+i, y+j, color)
			}
		}
	}
}

// DrawRect draws a rectangle outline
func (s *MemorySurface) DrawRect(x, y, width, height int, color color.RGBA) {
	// Draw top and bottom borders
	for i := 0; i < width; i++ {
		s.img.SetRGBA(x+i, y, color)
		s.img.SetRGBA(x+i, y+height-1, color)
	}
	
	// Draw left and right borders
	for i := 0; i < height; i++ {
		s.img.SetRGBA(x, y+i, color)
		s.img.SetRGBA(x+width-1, y+i, color)
	}
}

// FillRect fills a rectangle
func (s *MemorySurface) FillRect(x, y, width, height int, color color.RGBA) {
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			s.img.SetRGBA(x+i, y+j, color)
		}
	}
}

// DrawLine draws a line
func (s *MemorySurface) DrawLine(x1, y1, x2, y2 int, color color.RGBA) {
	// Simple line drawing algorithm
	dx := x2 - x1
	dy := y2 - y1
	
	steps := abs(dx)
	if abs(dy) > steps {
		steps = abs(dy)
	}
	
	xInc := float64(dx) / float64(steps)
	yInc := float64(dy) / float64(steps)
	
	x := float64(x1)
	y := float64(y1)
	
	for i := 0; i <= steps; i++ {
		s.img.SetRGBA(int(x), int(y), color)
		x += xInc
		y += yInc
	}
}

// FillCircle fills a circle with the given color
func (s *MemorySurface) FillCircle(x, y, radius int, color color.RGBA) {
	// Naive implementation - scan through a bounding square
	for px := x - radius; px <= x + radius; px++ {
		for py := y - radius; py <= y + radius; py++ {
			// Only draw if within bounds
			if px >= 0 && px < s.img.Bounds().Dx() && py >= 0 && py < s.img.Bounds().Dy() {
				dx := px - x
				dy := py - y
				distSquared := dx*dx + dy*dy
				if distSquared <= radius*radius {
					s.img.SetRGBA(px, py, color)
				}
			}
		}
	}
}

// DrawCircle draws a circle outline with the given color
func (s *MemorySurface) DrawCircle(x, y, radius int, color color.RGBA) {
	// Naive implementation - scan through a bounding square
	for px := x - radius; px <= x + radius; px++ {
		for py := y - radius; py <= y + radius; py++ {
			// Only draw if within bounds
			if px >= 0 && px < s.img.Bounds().Dx() && py >= 0 && py < s.img.Bounds().Dy() {
				dx := px - x
				dy := py - y
				distSquared := dx*dx + dy*dy
				
				// Draw pixels that are close to the radius
				radius2 := radius*radius
				if distSquared <= radius2 && distSquared >= (radius-1)*(radius-1) {
					s.img.SetRGBA(px, py, color)
				}
			}
		}
	}
}

// SetClipRect sets the clipping rectangle
func (s *MemorySurface) SetClipRect(x, y, width, height int) {
	// Not implemented for this simple test framework
}

// ResetClipRect resets the clipping rectangle
func (s *MemorySurface) ResetClipRect() {
	// Not implemented for this simple test framework
}

// Image returns the underlying image
func (s *MemorySurface) Image() *image.RGBA {
	return s.img
}

// DrawImage draws an image with the specified fit method
func (s *MemorySurface) DrawImage(img image.Image, x, y, width, height int, fitMethod components.ImageFitMethod) {
	// Simple implementation - just draw a placeholder rectangle
	s.FillRect(x, y, width, height, color.RGBA{200, 200, 200, 255})
	s.DrawRect(x, y, width, height, color.RGBA{150, 150, 150, 255})
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
} 