package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	// Use the new module path
	"github.com/aggnr/finch/components"
	"github.com/aggnr/finch/examples/todo_app/todo"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// Game implements the ebiten.Game interface
type Game struct {
	rootContainer *components.FlexContainer
	todoList      *todo.TodoList
	inputField    *components.TextArea
	addButton     *components.Button
	clearButton   *components.Button
	statusLabel   *components.Label
}

// NewGame creates a new game
func NewGame() *Game {
	game := &Game{}
	
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
	root.SetFlexDirection(components.FlexColumn)
	g.rootContainer = root
	
	// Create a title
	title := components.NewLabel("title", "Todo List Demo", 24, color.RGBA{50, 50, 50, 255})
	title.SetBounds(components.Rect{X: 0, Y: 20, Width: ScreenWidth, Height: 40})
	title.SetTextAlignment(components.TextAlignCenter)
	root.AddChild(title)
	
	// Create input container
	inputContainer := components.NewFlexContainer("input_container")
	inputContainer.SetBounds(components.Rect{X: 50, Y: 80, Width: ScreenWidth - 100, Height: 40})
	inputContainer.SetFlexDirection(components.FlexRow)
	root.AddChild(inputContainer)
	
	// Create text input field
	g.inputField = components.NewTextArea("todo_input")
	g.inputField.SetBounds(components.Rect{X: 0, Y: 0, Width: ScreenWidth - 200, Height: 40})
	g.inputField.SetText("")
	g.inputField.SetPlaceholder("Enter a new todo...")
	inputContainer.AddChild(g.inputField)
	
	// Create add button
	g.addButton = components.NewButton("add_button", "Add")
	g.addButton.SetBounds(components.Rect{X: ScreenWidth - 190, Y: 0, Width: 80, Height: 40})
	g.addButton.SetOnClick(func() {
		g.addTodo()
	})
	inputContainer.AddChild(g.addButton)
	
	// Create todo list container
	listContainer := components.NewFlexContainer("list_container")
	listContainer.SetBounds(components.Rect{X: 50, Y: 140, Width: ScreenWidth - 100, Height: ScreenHeight - 240})
	listContainer.SetBackgroundColor(color.RGBA{255, 255, 255, 255})
	root.AddChild(listContainer)
	
	// Create the todo list
	g.todoList = todo.NewTodoList("todo_list")
	g.todoList.SetBounds(components.Rect{X: 0, Y: 0, Width: ScreenWidth - 100, Height: ScreenHeight - 240})
	g.todoList.SetOnItemChanged(func(item todo.TodoItem) {
		g.updateStatus()
	})
	g.todoList.SetOnItemDeleted(func(id string) {
		g.updateStatus()
	})
	listContainer.AddChild(g.todoList)
	
	// Create bottom container
	bottomContainer := components.NewFlexContainer("bottom_container")
	bottomContainer.SetBounds(components.Rect{X: 50, Y: ScreenHeight - 80, Width: ScreenWidth - 100, Height: 40})
	bottomContainer.SetFlexDirection(components.FlexRow)
	root.AddChild(bottomContainer)
	
	// Create status label
	g.statusLabel = components.NewLabel("status", "0 items", 14, color.RGBA{100, 100, 100, 255})
	g.statusLabel.SetBounds(components.Rect{X: 0, Y: 0, Width: 150, Height: 40})
	bottomContainer.AddChild(g.statusLabel)
	
	// Create clear completed button
	g.clearButton = components.NewButton("clear_button", "Clear Completed")
	g.clearButton.SetBounds(components.Rect{X: ScreenWidth - 270, Y: 0, Width: 150, Height: 40})
	g.clearButton.SetOnClick(func() {
		g.todoList.ClearCompleted()
		g.updateStatus()
	})
	bottomContainer.AddChild(g.clearButton)
	
	// Add some sample todos
	g.todoList.AddTodo("Buy groceries")
	g.todoList.AddTodo("Finish project")
	g.todoList.AddTodo("Call John")
	
	// Update status
	g.updateStatus()
}

// addTodo adds a new todo from the input field
func (g *Game) addTodo() {
	text := g.inputField.GetText()
	if text != "" {
		g.todoList.AddTodo(text)
		g.inputField.SetText("")
		g.updateStatus()
	}
}

// updateStatus updates the status label
func (g *Game) updateStatus() {
	todos := g.todoList.GetTodos()
	
	// Count completed items
	completed := 0
	for _, todo := range todos {
		if todo.Done {
			completed++
		}
	}
	
	// Update status text
	if len(todos) == 0 {
		g.statusLabel.SetText("No items")
	} else if len(todos) == 1 {
		g.statusLabel.SetText("1 item")
	} else {
		g.statusLabel.SetText(fmt.Sprintf("%d items, %d completed", len(todos), completed))
	}
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
	
	// Handle keyboard events for adding todos with Enter key
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.addTodo()
	}
}

func main() {
	// Create the game
	game := NewGame()
	
	// Run the game
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Finch UI Todo List Demo")
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
} 