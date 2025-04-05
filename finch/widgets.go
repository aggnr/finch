package finch

import (
	"fmt"
	"image/color"

	"github.com/vijay/finch-ui/components"
)

// Text represents a text element like a label
type Text struct {
	label *components.Label
	ui    *UI
}

// Centered centers the text
func (t *Text) Centered() *Text {
	t.label.SetTextAlignment(components.TextAlignCenter)
	return t
}

// Size sets the font size
func (t *Text) Size(size int) *Text {
	t.label.SetFontSize(size)
	return t
}

// Color sets the text color
func (t *Text) Color(hexColor string) *Text {
	// Parse hex color (simplified)
	var r, g, b uint8 = 0, 0, 0
	fmt.Sscanf(hexColor, "#%02x%02x%02x", &r, &g, &b)
	t.label.SetTextColor(color.RGBA{r, g, b, 255})
	return t
}

// SetText updates the text content
func (t *Text) SetText(text string) *Text {
	t.label.SetText(text)
	return t
}

// Container represents a container element for layout
type Container struct {
	container *components.FlexContainer
	ui        *UI
}

// Background sets the background color
func (c *Container) Background(hexColor string) *Container {
	// Parse hex color (simplified)
	var r, g, b uint8 = 255, 255, 255
	fmt.Sscanf(hexColor, "#%02x%02x%02x", &r, &g, &b)
	c.container.SetBackgroundColor(color.RGBA{r, g, b, 255})
	return c
}

// Padding sets the padding
func (c *Container) Padding(padding int) *Container {
	c.container.SetBoxModel(components.BoxModel{
		Padding: components.Spacing{padding, padding, padding, padding},
	})
	return c
}

// Margin sets the margin
func (c *Container) Margin(margin int) *Container {
	c.container.SetBoxModel(components.BoxModel{
		Margin: components.Spacing{margin, margin, margin, margin},
	})
	return c
}

// Width sets the width
func (c *Container) Width(width interface{}) *Container {
	switch w := width.(type) {
	case int:
		bounds := c.container.Bounds()
		bounds.Width = w
		c.container.SetBounds(bounds)
	case string:
		// For percentage strings like "80%"
		var percentage int
		fmt.Sscanf(w, "%d%%", &percentage)
		if percentage > 0 {
			bounds := c.container.Bounds()
			bounds.Width = c.ui.width * percentage / 100
			c.container.SetBounds(bounds)
		}
	}
	return c
}

// Height sets the height
func (c *Container) Height(height int) *Container {
	bounds := c.container.Bounds()
	bounds.Height = height
	c.container.SetBounds(bounds)
	return c
}

// Grow makes the container take available space
func (c *Container) Grow(factor int) *Container {
	// In a real implementation, this would need a proper layout system
	bounds := c.container.Bounds()
	bounds.Height = c.ui.height - bounds.Y - 50 // Simplified
	c.container.SetBounds(bounds)
	return c
}

// Layout sets the layout direction and runs the builder function
func (c *Container) Layout(direction string, builder func(*Container)) *Container {
	if direction == "row" {
		c.container.SetFlexDirection(components.FlexRow)
	} else {
		c.container.SetFlexDirection(components.FlexColumn)
	}
	
	// Save the current parent
	originalParent := c.ui.currentParent
	
	// Set this container as the current parent
	c.ui.currentParent = c.container
	
	// Call the builder function
	if builder != nil {
		builder(c)
	}
	
	// Restore the original parent
	c.ui.currentParent = originalParent
	
	return c
}

// Text adds a text element to the container
func (c *Container) Text(text string) *Text {
	// Save the current parent
	originalParent := c.ui.currentParent
	
	// Set this container as the current parent
	c.ui.currentParent = c.container
	
	// Add the text
	textElement := c.ui.Text(text)
	
	// Restore the original parent
	c.ui.currentParent = originalParent
	
	return textElement
}

// Button adds a button to the container
func (c *Container) Button(label string) *Button {
	// Save the current parent
	originalParent := c.ui.currentParent
	
	// Set this container as the current parent
	c.ui.currentParent = c.container
	
	// Add the button
	button := c.ui.Button(label)
	
	// Restore the original parent
	c.ui.currentParent = originalParent
	
	return button
}

// TextInput adds a text input to the container
func (c *Container) TextInput(placeholder string) *TextInput {
	// Save the current parent
	originalParent := c.ui.currentParent
	
	// Set this container as the current parent
	c.ui.currentParent = c.container
	
	// Add the text input
	input := c.ui.TextInput(placeholder)
	
	// Restore the original parent
	c.ui.currentParent = originalParent
	
	return input
}

// Checkbox adds a checkbox to the container
func (c *Container) Checkbox(label string) *Checkbox {
	// Save the current parent
	originalParent := c.ui.currentParent
	
	// Set this container as the current parent
	c.ui.currentParent = c.container
	
	// Add the checkbox
	checkbox := c.ui.Checkbox(label)
	
	// Restore the original parent
	c.ui.currentParent = originalParent
	
	return checkbox
}

// RemoveAllChildren removes all child elements from this container
func (c *Container) RemoveAllChildren() {
	c.container.RemoveAllChildren()
}

// TodoList adds a todo list to the container (sample implementation)
func (c *Container) TodoList() *TodoList {
	// Create a todo list
	todoList := &TodoList{
		container: c.container,
		ui:        c.ui,
	}
	
	return todoList
}

// Button represents a button element
type Button struct {
	button *components.Button
	ui     *UI
}

// OnClick sets the click handler
func (b *Button) OnClick(handler func()) *Button {
	b.button.SetOnClick(handler)
	return b
}

// Width sets the button width
func (b *Button) Width(width int) *Button {
	bounds := b.button.Bounds()
	bounds.Width = width
	b.button.SetBounds(bounds)
	return b
}

// TextInput represents a text input field
type TextInput struct {
	input *components.TextArea
	ui    *UI
}

// Value gets the current text value
func (t *TextInput) Value() string {
	return t.input.GetText()
}

// SetValue sets the text value
func (t *TextInput) SetValue(value string) *TextInput {
	t.input.SetText(value)
	return t
}

// Clear clears the input
func (t *TextInput) Clear() *TextInput {
	t.input.SetText("")
	return t
}

// OnChange sets the change handler
func (t *TextInput) OnChange(handler func(string)) *TextInput {
	t.input.SetOnChange(handler)
	return t
}

// Checkbox represents a checkbox element
type Checkbox struct {
	checkbox *components.Checkbox
	label    *components.Label
	ui       *UI
}

// Value gets the current checked state
func (c *Checkbox) Value() bool {
	return c.checkbox.IsChecked()
}

// SetValue sets the checked state
func (c *Checkbox) SetValue(checked bool) *Checkbox {
	c.checkbox.SetChecked(checked)
	return c
}

// OnChange sets the change handler
func (c *Checkbox) OnChange(handler func(bool)) *Checkbox {
	c.checkbox.SetCheckedChanged(handler)
	return c
}

// BindValue binds a boolean pointer to the checkbox
func (c *Checkbox) BindValue(value *bool) *Checkbox {
	// Set initial value
	c.checkbox.SetChecked(*value)
	
	// Set up change handler
	c.checkbox.SetCheckedChanged(func(checked bool) {
		*value = checked
	})
	
	return c
}

// Column represents a column in a columns layout
type Column struct {
	container *components.FlexContainer
	ui        *UI
}

// Text adds a text element to the column
func (c *Column) Text(text string) *Text {
	// Save the current parent
	originalParent := c.ui.currentParent
	
	// Set this column as the current parent
	c.ui.currentParent = c.container
	
	// Add the text
	textElement := c.ui.Text(text)
	
	// Restore the original parent
	c.ui.currentParent = originalParent
	
	return textElement
}

// Button adds a button to the column
func (c *Column) Button(label string) *Button {
	// Save the current parent
	originalParent := c.ui.currentParent
	
	// Set this column as the current parent
	c.ui.currentParent = c.container
	
	// Add the button
	button := c.ui.Button(label)
	
	// Restore the original parent
	c.ui.currentParent = originalParent
	
	return button
}

// Tab represents a tab in a tabs layout
type Tab struct {
	header    *components.Button
	container *components.FlexContainer
	ui        *UI
}

// Text adds a text element to the tab
func (t *Tab) Text(text string) *Text {
	// Save the current parent
	originalParent := t.ui.currentParent
	
	// Set this tab as the current parent
	t.ui.currentParent = t.container
	
	// Add the text
	textElement := t.ui.Text(text)
	
	// Restore the original parent
	t.ui.currentParent = originalParent
	
	return textElement
}

// Button adds a button to the tab
func (t *Tab) Button(label string) *Button {
	// Save the current parent
	originalParent := t.ui.currentParent
	
	// Set this tab as the current parent
	t.ui.currentParent = t.container
	
	// Add the button
	button := t.ui.Button(label)
	
	// Restore the original parent
	t.ui.currentParent = originalParent
	
	return button
}

// TodoList adds a todo list to the tab
func (t *Tab) TodoList() *TodoList {
	// Create a todo list
	todoList := &TodoList{
		container: t.container,
		ui:        t.ui,
	}
	
	return todoList
}

// Checkbox adds a checkbox to the tab
func (t *Tab) Checkbox(label string) *Checkbox {
	// Save the current parent
	originalParent := t.ui.currentParent
	
	// Set this tab as the current parent
	t.ui.currentParent = t.container
	
	// Add the checkbox
	checkbox := t.ui.Checkbox(label)
	
	// Restore the original parent
	t.ui.currentParent = originalParent
	
	return checkbox
}

// State represents a reactive state value
type State struct {
	value    interface{}
	watchers []func(interface{})
}

// Update updates the state value using a transform function
func (s *State) Update(transform func(interface{}) interface{}) {
	newValue := transform(s.value)
	s.value = newValue
	
	// Notify watchers
	for _, watcher := range s.watchers {
		watcher(s.value)
	}
}

// Watch adds a watcher function that is called when the state changes
func (s *State) Watch(watcher func(interface{})) {
	s.watchers = append(s.watchers, watcher)
	
	// Call the watcher with the current value
	watcher(s.value)
}

// Value gets the current state value
func (s *State) Value() interface{} {
	return s.value
}

// TodoList represents a todo list (simplified example)
type TodoList struct {
	list      interface{} // This would be an actual TodoList component
	container *components.FlexContainer
	ui        *UI
	items     []interface{}
	onChange  func(interface{})
}

// BindItems binds a list of items to the todo list
func (t *TodoList) BindItems(state *State) *TodoList {
	// This is a simplified implementation
	t.items = state.Value().([]interface{})
	
	// Watch for changes
	state.Watch(func(value interface{}) {
		t.items = value.([]interface{})
		// Update the UI (simplified)
	})
	
	return t
}

// OnItemChange sets a handler for when an item changes
func (t *TodoList) OnItemChange(handler func(interface{})) *TodoList {
	t.onChange = handler
	return t
}

// FilterItems filters the displayed items
func (t *TodoList) FilterItems(filter func(interface{}) bool) *TodoList {
	// This would actually filter the displayed items
	return t
} 