package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	// Use the module path instead of a relative path
	"github.com/vijay/finch-ui/components"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// Game implements the ebiten.Game interface
type Game struct {
	rootContainer *components.FlexContainer
	label         *components.Label
	counter       int
}

// NewGame creates a new game
func NewGame() *Game {
	game := &Game{
		counter: 0,
	}
	
	// Initialize UI
	game.initUI()
	
	return game
}

// initUI initializes the UI
func (g *Game) initUI() {
	// Create root container
	root := components.NewFlexContainer("root")
	root.SetBounds(components.Rect{X: 0, Y: 0, Width: ScreenWidth, Height: ScreenHeight})
	root.SetBackgroundColor(color.RGBA{240, 240, 240, 255})
	g.rootContainer = root
	
	// Create a title
	title := components.NewLabel("title", "Finch UI Simple Demo", 24, color.RGBA{0, 0, 0, 255})
	title.SetBounds(components.Rect{X: 20, Y: 20, Width: 300, Height: 40})
	root.AddChild(title)
	
	// Create a counter label
	g.label = components.NewLabel("counter", "Button clicked: 0 times", 16, color.RGBA{50, 50, 50, 255})
	g.label.SetBounds(components.Rect{X: 20, Y: 70, Width: 300, Height: 30})
	root.AddChild(g.label)
	
	// Create a button
	button := components.NewButton("button", "Click Me")
	button.SetBounds(components.Rect{X: 20, Y: 110, Width: 120, Height: 40})
	button.SetOnClick(func() {
		g.counter++
		g.label.SetText(fmt.Sprintf("Button clicked: %d times", g.counter))
	})
	root.AddChild(button)
	
	// Create a checkbox
	checkbox := components.NewCheckbox("checkbox")
	checkbox.SetBounds(components.Rect{X: 20, Y: 160, Width: 20, Height: 20})
	root.AddChild(checkbox)
	
	// Add a label for the checkbox
	checkLabel := components.NewLabel("check_label", "Enable feature", 16, color.RGBA{50, 50, 50, 255})
	checkLabel.SetBounds(components.Rect{X: 50, Y: 160, Width: 150, Height: 20})
	root.AddChild(checkLabel)
	
	// Create a text input field
	textArea := components.NewTextArea("text_input")
	textArea.SetBounds(components.Rect{X: 20, Y: 200, Width: 300, Height: 100})
	textArea.SetText("Enter text here...")
	root.AddChild(textArea)
	
	// Create a dropdown
	dropdown := components.NewSelect("dropdown", []string{"Option 1", "Option 2", "Option 3"})
	dropdown.SetBounds(components.Rect{X: 20, Y: 320, Width: 200, Height: 30})
	root.AddChild(dropdown)
}

// Update handles game logic updates
func (g *Game) Update() error {
	// Handle input
	g.handleInput()
	
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Create a draw surface
	surface := components.NewEbitenDrawSurface(screen)
	
	// Draw the UI
	g.rootContainer.Draw(surface)
}

// Layout implements the ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// handleInput handles input events
func (g *Game) handleInput() {
	// Get mouse position
	x, y := ebiten.CursorPosition()
	
	// Handle mouse events
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.rootContainer.HandleMouseDown(x, y)
	}
	
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.rootContainer.HandleMouseUp(x, y)
	}
	
	g.rootContainer.HandleMouseMove(x, y)
}

func main() {
	// Create the game
	game := NewGame()
	
	// Run the game
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Finch UI Simple Demo")
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
} 