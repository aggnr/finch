package finch

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vijay/finch-ui/components"
)

// UI is the main entry point for the Finch UI framework
type UI struct {
	rootContainer *components.FlexContainer
	width         int
	height        int
	title         string
	currentParent components.Element
}

// PageConfig represents configuration for the page
type PageConfig struct {
	Title  string
	Layout string // "column" or "row"
	Width  int
	Height int
}

// New creates a new Finch UI instance
func New() *UI {
	root := components.NewFlexContainer("root")
	ui := &UI{
		rootContainer: root,
		currentParent: root,
		width:         800,
		height:        600,
		title:         "Finch UI App",
	}
	
	// Set default properties
	root.SetBounds(components.Rect{X: 0, Y: 0, Width: ui.width, Height: ui.height})
	root.SetBackgroundColor(color.RGBA{240, 240, 240, 255})
	root.SetFlexDirection(components.FlexColumn)
	
	return ui
}

// SetPageConfig configures the UI page
func (ui *UI) SetPageConfig(title string, layout string) *UI {
	ui.title = title
	
	if layout == "row" {
		ui.rootContainer.SetFlexDirection(components.FlexRow)
	} else {
		ui.rootContainer.SetFlexDirection(components.FlexColumn)
	}
	
	return ui
}

// Title adds a title to the UI
func (ui *UI) Title(text string) *Text {
	title := components.NewLabel("title_"+randomID(), text, 24, color.RGBA{50, 50, 50, 255})
	title.SetBounds(components.Rect{X: 0, Y: 20, Width: ui.width, Height: 40})
	
	ui.currentParent.AddChild(title)
	
	return &Text{
		label: title,
		ui:    ui,
	}
}

// Text adds a text element to the UI
func (ui *UI) Text(text string) *Text {
	label := components.NewLabel("text_"+randomID(), text, 16, color.RGBA{0, 0, 0, 255})
	label.SetBounds(components.Rect{X: 0, Y: 0, Width: ui.width, Height: 20})
	
	ui.currentParent.AddChild(label)
	
	return &Text{
		label: label,
		ui:    ui,
	}
}

// Container creates a container for organizing UI elements
func (ui *UI) Container() *Container {
	container := components.NewFlexContainer("container_" + randomID())
	container.SetBounds(components.Rect{X: 0, Y: 0, Width: ui.width, Height: 100})
	container.SetFlexDirection(components.FlexColumn)
	
	ui.currentParent.AddChild(container)
	
	return &Container{
		container: container,
		ui:        ui,
	}
}

// Button adds a button to the UI
func (ui *UI) Button(label string) *Button {
	button := components.NewButton("button_"+randomID(), label)
	button.SetBounds(components.Rect{X: 0, Y: 0, Width: 120, Height: 40})
	
	ui.currentParent.AddChild(button)
	
	return &Button{
		button: button,
		ui:     ui,
	}
}

// TextInput adds a text input field to the UI
func (ui *UI) TextInput(placeholder string) *TextInput {
	input := components.NewTextArea("input_" + randomID())
	input.SetBounds(components.Rect{X: 0, Y: 0, Width: ui.width - 150, Height: 40})
	input.SetPlaceholder(placeholder)
	
	ui.currentParent.AddChild(input)
	
	return &TextInput{
		input: input,
		ui:    ui,
	}
}

// Checkbox adds a checkbox to the UI
func (ui *UI) Checkbox(label string) *Checkbox {
	// Create a container for the checkbox and label
	container := components.NewFlexContainer("checkbox_container_" + randomID())
	container.SetBounds(components.Rect{X: 0, Y: 0, Width: ui.width, Height: 30})
	container.SetFlexDirection(components.FlexRow)
	
	// Create the checkbox
	checkbox := components.NewCheckbox("checkbox_" + randomID())
	checkbox.SetBounds(components.Rect{X: 0, Y: 5, Width: 20, Height: 20})
	
	// Create the label
	textLabel := components.NewLabel("checkbox_label_"+randomID(), label, 16, color.RGBA{0, 0, 0, 255})
	textLabel.SetBounds(components.Rect{X: 30, Y: 5, Width: ui.width - 50, Height: 20})
	
	// Add to container
	container.AddChild(checkbox)
	container.AddChild(textLabel)
	
	// Add container to parent
	ui.currentParent.AddChild(container)
	
	return &Checkbox{
		checkbox: checkbox,
		label:    textLabel,
		ui:       ui,
	}
}

// Columns creates a set of columns
func (ui *UI) Columns(count int, builder func([]*Column)) *UI {
	columnsContainer := components.NewFlexContainer("columns_" + randomID())
	columnsContainer.SetBounds(components.Rect{X: 0, Y: 0, Width: ui.width, Height: 100})
	columnsContainer.SetFlexDirection(components.FlexRow)
	
	columns := make([]*Column, count)
	columnWidth := ui.width / count
	
	for i := 0; i < count; i++ {
		colContainer := components.NewFlexContainer(fmt.Sprintf("column_%d_%s", i, randomID()))
		colContainer.SetBounds(components.Rect{X: i * columnWidth, Y: 0, Width: columnWidth, Height: 100})
		colContainer.SetFlexDirection(components.FlexColumn)
		
		columns[i] = &Column{
			container: colContainer,
			ui:        ui,
		}
		
		columnsContainer.AddChild(colContainer)
	}
	
	ui.currentParent.AddChild(columnsContainer)
	
	// Save the original parent
	originalParent := ui.currentParent
	
	// Call the builder function with our columns
	if builder != nil {
		builder(columns)
	}
	
	// Restore the original parent
	ui.currentParent = originalParent
	
	return ui
}

// Tabs creates a set of tabs
func (ui *UI) Tabs(names []string, builder func([]*Tab)) *UI {
	tabsContainer := components.NewFlexContainer("tabs_container_" + randomID())
	tabsContainer.SetBounds(components.Rect{X: 0, Y: 0, Width: ui.width, Height: 300})
	tabsContainer.SetFlexDirection(components.FlexColumn)
	
	// Create tab headers container
	tabHeadersContainer := components.NewFlexContainer("tab_headers_" + randomID())
	tabHeadersContainer.SetBounds(components.Rect{X: 0, Y: 0, Width: ui.width, Height: 40})
	tabHeadersContainer.SetFlexDirection(components.FlexRow)
	tabHeadersContainer.SetBackgroundColor(color.RGBA{220, 220, 220, 255})
	
	// Create content container
	contentContainer := components.NewFlexContainer("tab_content_" + randomID())
	contentContainer.SetBounds(components.Rect{X: 0, Y: 40, Width: ui.width, Height: 260})
	contentContainer.SetFlexDirection(components.FlexColumn)
	
	// Add to tabs container
	tabsContainer.AddChild(tabHeadersContainer)
	tabsContainer.AddChild(contentContainer)
	
	// Create tabs
	tabs := make([]*Tab, len(names))
	tabWidth := ui.width / len(names)
	
	// Create content containers (only one is visible at a time)
	for i, name := range names {
		// Create tab header
		tabHeader := components.NewButton("tab_header_"+randomID(), name)
		tabHeader.SetBounds(components.Rect{X: i * tabWidth, Y: 0, Width: tabWidth, Height: 40})
		tabHeadersContainer.AddChild(tabHeader)
		
		// Create content container for this tab
		tabContent := components.NewFlexContainer("tab_content_"+randomID())
		tabContent.SetBounds(components.Rect{X: 0, Y: 0, Width: ui.width, Height: 260})
		tabContent.SetFlexDirection(components.FlexColumn)
		tabContent.SetVisible(i == 0) // Only first tab visible by default
		contentContainer.AddChild(tabContent)
		
		// Create tab object
		tabs[i] = &Tab{
			header:    tabHeader,
			container: tabContent,
			ui:        ui,
		}
		
		// Set up tab switching
		index := i // Capture index for closure
		tabHeader.SetOnClick(func() {
			// Hide all tab contents
			for j := 0; j < len(tabs); j++ {
				tabs[j].container.SetVisible(j == index)
			}
		})
	}
	
	ui.currentParent.AddChild(tabsContainer)
	
	// Save the original parent
	originalParent := ui.currentParent
	
	// Call the builder function with our tabs
	if builder != nil {
		builder(tabs)
	}
	
	// Restore the original parent
	ui.currentParent = originalParent
	
	return ui
}

// State creates a new reactive state value
func (ui *UI) State(initialValue interface{}) *State {
	return &State{
		value:    initialValue,
		watchers: make([]func(interface{}), 0),
	}
}

// Run starts the UI application
func (ui *UI) Run(width, height int) {
	ui.width = width
	ui.height = height
	ui.rootContainer.SetBounds(components.Rect{X: 0, Y: 0, Width: width, Height: height})
	
	// Create the game
	game := &Game{
		rootContainer: ui.rootContainer,
		width:         width,
		height:        height,
	}
	
	// Run the game
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle(ui.title)
	
	if err := ebiten.RunGame(game); err != nil {
		fmt.Printf("Error running game: %v\n", err)
	}
}

// Game implements the ebiten.Game interface
type Game struct {
	rootContainer *components.FlexContainer
	width         int
	height        int
}

// Update implements ebiten.Game's Update method
func (g *Game) Update() error {
	// Handle input in a simpler way
	x, y := ebiten.CursorPosition()
	
	// Mouse events
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.rootContainer.HandleMouseDown(x, y)
	} else {
		g.rootContainer.HandleMouseUp(x, y)
	}
	
	g.rootContainer.HandleMouseMove(x, y)
	
	return nil
}

// Draw implements ebiten.Game's Draw method
func (g *Game) Draw(screen *ebiten.Image) {
	// Create a draw surface
	surface := components.NewEbitenDrawSurface(screen)
	
	// Draw the UI
	g.rootContainer.Draw(surface)
}

// Layout implements ebiten.Game's Layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

// Helper function to generate random IDs
func randomID() string {
	return fmt.Sprintf("%d", ebiten.TPS())
} 