package components

import (
	"image/color"
)

// Button represents a button in the UI
type Button struct {
	*Node
	text           string
	onClick        func()
	backgroundColor color.RGBA
	textColor      color.RGBA
	hoverColor     color.RGBA
	pressedColor   color.RGBA
	fontSize       int
	hovered        bool
	pressed        bool
	disabled       bool
}

// NewButton creates a new button
func NewButton(id string, text string) *Button {
	return &Button{
		Node:           NewNode(id),
		text:           text,
		onClick:        nil,
		backgroundColor: color.RGBA{200, 200, 200, 255},
		textColor:      color.RGBA{0, 0, 0, 255},
		hoverColor:     color.RGBA{220, 220, 220, 255},
		pressedColor:   color.RGBA{180, 180, 180, 255},
		fontSize:       14,
		hovered:        false,
		pressed:        false,
		disabled:       false,
	}
}

// SetDisabled sets whether the button is disabled
func (b *Button) SetDisabled(disabled bool) {
	b.disabled = disabled
}

// IsDisabled returns whether the button is disabled
func (b *Button) IsDisabled() bool {
	return b.disabled
}

// Draw draws the button
func (b *Button) Draw(surface DrawSurface) {
	if !b.IsVisible() {
		return
	}
	
	bounds := b.ComputedBounds()
	
	// Determine the background color based on button state
	bg := b.backgroundColor
	if b.disabled {
		bg = color.RGBA{150, 150, 150, 255}
	} else if b.pressed {
		bg = b.pressedColor
	} else if b.hovered {
		bg = b.hoverColor
	}
	
	// Draw the button background
	surface.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, bg)
	
	// Draw the button border
	surface.DrawRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{100, 100, 100, 255})
	
	// Calculate text position to center it
	textWidth := len(b.text) * b.fontSize / 2
	textX := bounds.X + (bounds.Width - textWidth) / 2
	textY := bounds.Y + (bounds.Height - b.fontSize) / 2
	
	// Draw text with slight offset when pressed
	if b.pressed && !b.disabled {
		textX += 1
		textY += 1
	}
	
	// Determine text color
	textColor := b.textColor
	if b.disabled {
		textColor = color.RGBA{100, 100, 100, 255}
	}
	
	// Draw the text
	surface.DrawText(b.text, textX, textY, textColor, b.fontSize)
	
	// Draw children (if any)
	for _, child := range b.Children() {
		child.Draw(surface)
	}
}

// SetOnClick sets the click handler
func (b *Button) SetOnClick(handler func()) {
	b.onClick = handler
}

// SetBackgroundColor sets the button background color
func (b *Button) SetBackgroundColor(color color.RGBA) {
	b.backgroundColor = color
}

// SetTextColor sets the button text color
func (b *Button) SetTextColor(color color.RGBA) {
	b.textColor = color
}

// SetText sets the button text
func (b *Button) SetText(text string) {
	b.text = text
}

// SetFontSize sets the button font size
func (b *Button) SetFontSize(size int) {
	b.fontSize = size
}

// HandleMouseDown handles mouse down events
func (b *Button) HandleMouseDown(x, y int) bool {
	if b.disabled {
		return false
	}
	
	bounds := b.ComputedBounds()
	if PointInRect(Point{x, y}, bounds) {
		b.pressed = true
		
		// Check if any children handle the event
		for i := len(b.Children()) - 1; i >= 0; i-- {
			child := b.Children()[i]
			if child.HandleMouseDown(x, y) {
				return true
			}
		}
		
		return true
	}
	return false
}

// HandleMouseUp handles mouse up events
func (b *Button) HandleMouseUp(x, y int) bool {
	wasPressed := b.pressed
	b.pressed = false
	
	if b.disabled {
		return false
	}
	
	bounds := b.ComputedBounds()
	if wasPressed && PointInRect(Point{x, y}, bounds) {
		// Execute onClick handler
		if b.onClick != nil {
			b.onClick()
		}
		
		return true
	}
	
	// Still try children even if this element didn't handle it
	for i := len(b.Children()) - 1; i >= 0; i-- {
		child := b.Children()[i]
		if child.HandleMouseUp(x, y) {
			return true
		}
	}
	
	return false
}

// HandleMouseMove handles mouse move events
func (b *Button) HandleMouseMove(x, y int) bool {
	wasHovered := b.hovered
	bounds := b.ComputedBounds()
	b.hovered = PointInRect(Point{x, y}, bounds)
	
	// Check if any children handle the event
	for i := len(b.Children()) - 1; i >= 0; i-- {
		child := b.Children()[i]
		if child.HandleMouseMove(x, y) {
			return true
		}
	}
	
	return b.hovered || wasHovered != b.hovered
}

// WasClicked returns whether the button was clicked (for compatibility)
func (b *Button) WasClicked() bool {
	return b.pressed
}

// IsHovered returns whether the button is hovered
func (b *Button) IsHovered() bool {
	return b.hovered
}

// IsPressed returns whether the button is pressed
func (b *Button) IsPressed() bool {
	return b.pressed
}

// Checkbox represents a checkbox in the UI
type Checkbox struct {
	*Node
	checked        bool
	checkedChanged func(bool)
}

// NewCheckbox creates a new checkbox
func NewCheckbox(id string) *Checkbox {
	return &Checkbox{
		Node: NewNode(id),
		checked: false,
	}
}

// SetChecked sets whether the checkbox is checked
func (c *Checkbox) SetChecked(checked bool) {
	c.checked = checked
}

// IsChecked returns whether the checkbox is checked
func (c *Checkbox) IsChecked() bool {
	return c.checked
}

// SetCheckedChanged sets the handler for when the checked state changes
func (c *Checkbox) SetCheckedChanged(handler func(bool)) {
	c.checkedChanged = handler
}

// Draw draws the checkbox
func (c *Checkbox) Draw(surface DrawSurface) {
	if !c.IsVisible() {
		return
	}
	
	bounds := c.ComputedBounds()
	
	// Draw checkbox background
	surface.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{255, 255, 255, 255})
	
	// Draw checkbox border
	surface.DrawRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{100, 100, 100, 255})
	
	// Draw check mark if checked
	if c.checked {
		// Simple X mark
		surface.DrawLine(
			bounds.X + 3, 
			bounds.Y + 3, 
			bounds.X + bounds.Width - 3, 
			bounds.Y + bounds.Height - 3, 
			color.RGBA{0, 0, 0, 255})
		surface.DrawLine(
			bounds.X + bounds.Width - 3, 
			bounds.Y + 3, 
			bounds.X + 3, 
			bounds.Y + bounds.Height - 3, 
			color.RGBA{0, 0, 0, 255})
	}
	
	// Draw children (if any)
	for _, child := range c.Children() {
		child.Draw(surface)
	}
}

// HandleMouseDown handles mouse down events
func (c *Checkbox) HandleMouseDown(x, y int) bool {
	bounds := c.ComputedBounds()
	if PointInRect(Point{x, y}, bounds) {
		// Toggle checked state
		c.checked = !c.checked
		
		// Call handler if set
		if c.checkedChanged != nil {
			c.checkedChanged(c.checked)
		}
		
		return true
	}
	return false
}

// HandleMouseMove handles mouse move events
func (c *Checkbox) HandleMouseMove(x, y int) bool {
	return false
} 