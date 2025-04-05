package todo

import (
	"fmt"
	"image/color"

	"github.com/vijay/finch-ui/components"
)

// TodoItem represents a single todo item with its data
type TodoItem struct {
	ID      string
	Text    string
	Done    bool
}

// Todo represents a single todo item component in the UI
type Todo struct {
	*components.Node
	item           TodoItem
	checkbox       *components.Checkbox
	textLabel      *components.Label
	deleteButton   *components.Button
	onDelete       func(id string)
	onToggle       func(id string, done bool)
	backgroundColor color.RGBA
	hoverColor     color.RGBA
	hovered        bool
}

// NewTodo creates a new todo item component
func NewTodo(id string, item TodoItem, onDelete func(id string), onToggle func(id string, done bool)) *Todo {
	todo := &Todo{
		Node:           components.NewNode(id),
		item:           item,
		onDelete:       onDelete,
		onToggle:       onToggle,
		backgroundColor: color.RGBA{255, 255, 255, 255},
		hoverColor:     color.RGBA{240, 240, 240, 255},
	}
	
	// Create checkbox for completion status
	todo.checkbox = components.NewCheckbox(id + "_checkbox")
	todo.checkbox.SetChecked(item.Done)
	todo.checkbox.SetCheckedChanged(func(checked bool) {
		todo.item.Done = checked
		if todo.onToggle != nil {
			todo.onToggle(todo.item.ID, checked)
		}
	})
	
	// Create label for todo text
	todo.textLabel = components.NewLabel(id+"_text", item.Text, 14, color.RGBA{0, 0, 0, 255})
	
	// Create delete button
	todo.deleteButton = components.NewButton(id+"_delete", "×")
	todo.deleteButton.SetFontSize(16)
	todo.deleteButton.SetOnClick(func() {
		if todo.onDelete != nil {
			todo.onDelete(todo.item.ID)
		}
	})
	
	// Add components to the todo item
	todo.AddChild(todo.checkbox)
	todo.AddChild(todo.textLabel)
	todo.AddChild(todo.deleteButton)
	
	return todo
}

// Layout positions the todo item's child components
func (t *Todo) Layout() {
	bounds := t.Bounds()
	
	// Position checkbox on the left
	t.checkbox.SetBounds(components.Rect{
		X:      5,
		Y:      (bounds.Height - 20) / 2,
		Width:  20,
		Height: 20,
	})
	
	// Position delete button on the right
	t.deleteButton.SetBounds(components.Rect{
		X:      bounds.Width - 40,
		Y:      (bounds.Height - 30) / 2,
		Width:  30,
		Height: 30,
	})
	
	// Position text label in the middle
	t.textLabel.SetBounds(components.Rect{
		X:      35,
		Y:      (bounds.Height - 20) / 2,
		Width:  bounds.Width - 80,
		Height: 20,
	})
	
	// Apply strikethrough style for completed todos
	if t.item.Done {
		t.textLabel.SetTextColor(color.RGBA{150, 150, 150, 255})
	} else {
		t.textLabel.SetTextColor(color.RGBA{0, 0, 0, 255})
	}
}

// Draw draws the todo item and its children
func (t *Todo) Draw(surface components.DrawSurface) {
	if !t.IsVisible() {
		return
	}
	
	// Update layout before drawing
	t.Layout()
	
	bounds := t.ComputedBounds()
	
	// Draw background
	bgColor := t.backgroundColor
	if t.hovered {
		bgColor = t.hoverColor
	}
	surface.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, bgColor)
	
	// Draw light border at the bottom
	surface.DrawLine(
		bounds.X,
		bounds.Y + bounds.Height - 1,
		bounds.X + bounds.Width,
		bounds.Y + bounds.Height - 1,
		color.RGBA{220, 220, 220, 255},
	)
	
	// Draw children
	for _, child := range t.Children() {
		child.Draw(surface)
	}
}

// HandleMouseMove handles mouse move events
func (t *Todo) HandleMouseMove(x, y int) bool {
	prevHovered := t.hovered
	bounds := t.ComputedBounds()
	t.hovered = components.PointInRect(components.Point{x, y}, bounds)
	
	// Check if any children handle the event
	for i := len(t.Children()) - 1; i >= 0; i-- {
		child := t.Children()[i]
		if child.HandleMouseMove(x, y) {
			return true
		}
	}
	
	// Return true if hover state changed
	return t.hovered || prevHovered != t.hovered
}

// HandleMouseDown handles mouse down events
func (t *Todo) HandleMouseDown(x, y int) bool {
	bounds := t.ComputedBounds()
	if components.PointInRect(components.Point{x, y}, bounds) {
		// Check if any children handle the event
		for i := len(t.Children()) - 1; i >= 0; i-- {
			child := t.Children()[i]
			if child.HandleMouseDown(x, y) {
				return true
			}
		}
		return true
	}
	return false
}

// HandleMouseUp handles mouse up events
func (t *Todo) HandleMouseUp(x, y int) bool {
	// Check if any children handle the event
	for i := len(t.Children()) - 1; i >= 0; i-- {
		child := t.Children()[i]
		if child.HandleMouseUp(x, y) {
			return true
		}
	}
	return false
}

// SetBackgroundColor sets the todo item background color
func (t *Todo) SetBackgroundColor(color color.RGBA) {
	t.backgroundColor = color
}

// SetHoverColor sets the todo item hover color
func (t *Todo) SetHoverColor(color color.RGBA) {
	t.hoverColor = color
}

// GetItem returns the todo item data
func (t *Todo) GetItem() TodoItem {
	// Make sure we have the latest checkbox value
	t.item.Done = t.checkbox.IsChecked()
	return t.item
}

// SetDone sets the completion status of the todo
func (t *Todo) SetDone(done bool) {
	t.item.Done = done
	t.checkbox.SetChecked(done)
}

// SetText sets the todo text
func (t *Todo) SetText(text string) {
	t.item.Text = text
	t.textLabel.SetText(text)
}

// GetCheckbox returns the checkbox component of the todo
func (t *Todo) GetCheckbox() *components.Checkbox {
	return t.checkbox
}

// GetTextLabel returns the text label component of the todo
func (t *Todo) GetTextLabel() *components.Label {
	return t.textLabel
}

// GetDeleteButton returns the delete button component of the todo
func (t *Todo) GetDeleteButton() *components.Button {
	return t.deleteButton
}

// TodoList represents a list of todo items
type TodoList struct {
	*components.FlexContainer
	todos         map[string]*Todo
	nextID        int
	onItemChanged func(item TodoItem)
	onItemDeleted func(id string)
}

// NewTodoList creates a new todo list
func NewTodoList(id string) *TodoList {
	list := &TodoList{
		FlexContainer: components.NewFlexContainer(id),
		todos:         make(map[string]*Todo),
		nextID:        1,
	}
	
	// Set vertical layout
	list.SetFlexDirection(components.FlexColumn)
	
	return list
}

// AddTodo adds a new todo item to the list
func (tl *TodoList) AddTodo(text string) *Todo {
	// Create a unique ID for the new todo
	id := fmt.Sprintf("todo_%d", tl.nextID)
	tl.nextID++
	
	// Create a new todo item
	todoItem := TodoItem{
		ID:   id,
		Text: text,
		Done: false,
	}
	
	// Create a new todo component
	todo := NewTodo(id, todoItem, tl.handleDelete, tl.handleToggle)
	todo.SetBounds(components.Rect{
		X:      0,
		Y:      0,
		Width:  tl.Bounds().Width,
		Height: 40,
	})
	
	// Store the todo in our map
	tl.todos[id] = todo
	
	// Add the todo to the container
	tl.AddChild(todo)
	
	// Return the new todo
	return todo
}

// GetTodos returns all todo items
func (tl *TodoList) GetTodos() []TodoItem {
	result := make([]TodoItem, 0, len(tl.todos))
	for _, todo := range tl.todos {
		result = append(result, todo.GetItem())
	}
	return result
}

// RemoveTodo removes a todo item from the list
func (tl *TodoList) RemoveTodo(id string) {
	if todo, ok := tl.todos[id]; ok {
		// Remove from the container
		tl.RemoveChild(todo)
		
		// Remove from our map
		delete(tl.todos, id)
		
		// Notify if callback is set
		if tl.onItemDeleted != nil {
			tl.onItemDeleted(id)
		}
	}
}

// SetOnItemChanged sets the callback for when a todo item changes
func (tl *TodoList) SetOnItemChanged(callback func(item TodoItem)) {
	tl.onItemChanged = callback
}

// SetOnItemDeleted sets the callback for when a todo item is deleted
func (tl *TodoList) SetOnItemDeleted(callback func(id string)) {
	tl.onItemDeleted = callback
}

// handleDelete is the internal handler for when a todo's delete button is clicked
func (tl *TodoList) handleDelete(id string) {
	tl.RemoveTodo(id)
}

// handleToggle is the internal handler for when a todo's checkbox is toggled
func (tl *TodoList) handleToggle(id string, done bool) {
	if todo, ok := tl.todos[id]; ok {
		// Get the updated item
		updatedItem := todo.GetItem()
		
		// Notify if callback is set
		if tl.onItemChanged != nil {
			tl.onItemChanged(updatedItem)
		}
	}
}

// ClearCompleted removes all completed todos from the list
func (tl *TodoList) ClearCompleted() {
	// Make a list of IDs to remove
	toRemove := make([]string, 0)
	for id, todo := range tl.todos {
		if todo.GetItem().Done {
			toRemove = append(toRemove, id)
		}
	}
	
	// Remove each completed todo
	for _, id := range toRemove {
		tl.RemoveTodo(id)
	}
}

// UpdateLayout updates the layout of all todo items
func (tl *TodoList) UpdateLayout() {
	// Update the width of all todos to match the container width
	y := 0
	for _, todo := range tl.todos {
		bounds := todo.Bounds()
		todo.SetBounds(components.Rect{
			X:      0,
			Y:      y,
			Width:  tl.Bounds().Width,
			Height: bounds.Height,
		})
		y += bounds.Height
	}
}

// GetTodoByID returns a todo by its ID
func (tl *TodoList) GetTodoByID(id string) (*Todo, bool) {
	todo, ok := tl.todos[id]
	return todo, ok
} 