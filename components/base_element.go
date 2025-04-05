package components

import (
	"fmt"
)

// BaseElement provides default implementations for the Element interface
type BaseElement struct {
	id       string
	bounds   Rect
	parent   Element
	children []Element
	mouseOver bool
	pressed   bool
}

// NewBaseElement creates a new base element
func NewBaseElement(id string) *BaseElement {
	return &BaseElement{
		id:       id,
		bounds:   Rect{},
		parent:   nil,
		children: make([]Element, 0),
		mouseOver: false,
		pressed:   false,
	}
}

// ID returns the element's ID
func (b *BaseElement) ID() string {
	return b.id
}

// SetID sets the element's ID
func (b *BaseElement) SetID(id string) {
	b.id = id
}

// Bounds returns the element's bounds
func (b *BaseElement) Bounds() Rect {
	return b.bounds
}

// SetBounds sets the element's bounds
func (b *BaseElement) SetBounds(bounds Rect) {
	b.bounds = bounds
}

// Parent returns the element's parent
func (b *BaseElement) Parent() Element {
	return b.parent
}

// SetParent sets the element's parent
func (b *BaseElement) SetParent(parent Element) {
	b.parent = parent
}

// Children returns the element's children
func (b *BaseElement) Children() []Element {
	return b.children
}

// AddChild adds a child element
func (b *BaseElement) AddChild(child Element) {
	b.children = append(b.children, child)
	child.SetParent(b)
	fmt.Printf("Added child %s to %s\n", child.ID(), b.id)
}

// RemoveChild removes a child element
func (b *BaseElement) RemoveChild(child Element) {
	for i, c := range b.children {
		if c == child {
			b.children = append(b.children[:i], b.children[i+1:]...)
			break
		}
	}
}

// RemoveAllChildren removes all child elements
func (b *BaseElement) RemoveAllChildren() {
	b.children = make([]Element, 0)
}

// IsMouseOver checks if the mouse is over the element
func (b *BaseElement) IsMouseOver(x, y int) bool {
	p := Point{X: x, Y: y}
	result := PointInRect(p, b.bounds)
	if result {
		fmt.Printf("Mouse is over %s at (%d,%d)\n", b.id, x, y)
	}
	return result
}

// HandleMouseDown handles mouse down events
func (b *BaseElement) HandleMouseDown(x, y int) bool {
	if b.IsMouseOver(x, y) {
		b.pressed = true
		fmt.Printf("MouseDown on %s\n", b.id)
		
		// Check if any children handle the event
		for i := len(b.children) - 1; i >= 0; i-- {
			child := b.children[i]
			if child.HandleMouseDown(x, y) {
				return true
			}
		}
		
		return true
	}
	return false
}

// HandleMouseUp handles mouse up events
func (b *BaseElement) HandleMouseUp(x, y int) bool {
	wasPressed := b.pressed
	b.pressed = false
	
	if wasPressed && b.IsMouseOver(x, y) {
		fmt.Printf("MouseUp on %s\n", b.id)
		
		// Check if any children handle the event
		for i := len(b.children) - 1; i >= 0; i-- {
			child := b.children[i]
			if child.HandleMouseUp(x, y) {
				return true
			}
		}
		
		return true
	}
	
	// Still try children even if this element didn't handle it
	for i := len(b.children) - 1; i >= 0; i-- {
		child := b.children[i]
		if child.HandleMouseUp(x, y) {
			return true
		}
	}
	
	return false
}

// HandleMouseMove handles mouse move events
func (b *BaseElement) HandleMouseMove(x, y int) bool {
	wasOver := b.mouseOver
	b.mouseOver = b.IsMouseOver(x, y)
	
	if b.mouseOver != wasOver {
		if b.mouseOver {
			fmt.Printf("MouseEnter on %s\n", b.id)
		} else {
			fmt.Printf("MouseLeave on %s\n", b.id)
		}
	}
	
	// Check if any children handle the event
	for i := len(b.children) - 1; i >= 0; i-- {
		child := b.children[i]
		if child.HandleMouseMove(x, y) {
			return true
		}
	}
	
	return b.mouseOver
}

// Draw draws the base element and its children
func (b *BaseElement) Draw(surface DrawSurface) {
	// The base element doesn't draw anything itself
	// But it does draw its children
	for _, child := range b.children {
		child.Draw(surface)
	}
}

// Update updates the element state
func (b *BaseElement) Update() {
	// Update all children
	for _, child := range b.children {
		child.Update()
	}
} 