package components

import (
	"image/color"
)

// FlexContainer represents a flex container for layout
type FlexContainer struct {
	*Node
	backgroundColor color.RGBA
	flexDirection   FlexDirection
	alignItems      Alignment
	justifyContent  Alignment
	spacing         int // Space between items
}

// NewFlexContainer creates a new flex container
func NewFlexContainer(id string) *FlexContainer {
	return &FlexContainer{
		Node:           NewNode(id),
		backgroundColor: color.RGBA{0, 0, 0, 0}, // Transparent by default
		flexDirection:   FlexRow,                // Default to row
		alignItems:      AlignStart,             // Default to start alignment
		justifyContent:  AlignStart,             // Default to start justification
		spacing:         5,                      // Default spacing between items
	}
}

// SetFlexDirection sets the flex direction
func (f *FlexContainer) SetFlexDirection(direction FlexDirection) {
	f.flexDirection = direction
}

// SetAlignItems sets the align items property
func (f *FlexContainer) SetAlignItems(alignment Alignment) {
	f.alignItems = alignment
}

// SetJustifyContent sets the justify content property
func (f *FlexContainer) SetJustifyContent(alignment Alignment) {
	f.justifyContent = alignment
}

// SetBackgroundColor sets the background color
func (f *FlexContainer) SetBackgroundColor(color color.RGBA) {
	f.backgroundColor = color
}

// Draw draws the flex container and its children
func (f *FlexContainer) Draw(surface DrawSurface) {
	if !f.IsVisible() {
		return
	}
	
	bounds := f.ComputedBounds()
	
	// Draw background if not transparent
	if f.backgroundColor.A > 0 {
		surface.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, f.backgroundColor)
	}
	
	// Perform layout calculations for children here...
	// (Simplified - a real implementation would position children according to flex rules)
	
	// Draw children
	for _, child := range f.Children() {
		child.Draw(surface)
	}
}

// HandleMouseDown handles mouse down events
func (f *FlexContainer) HandleMouseDown(x, y int) bool {
	bounds := f.ComputedBounds()
	if PointInRect(Point{x, y}, bounds) {
		// Check if any children handle the event (in reverse order for proper z-index)
		for i := len(f.Children()) - 1; i >= 0; i-- {
			child := f.Children()[i]
			if child.HandleMouseDown(x, y) {
				return true
			}
		}
		
		// If no children handled it, this container handles it
		return true
	}
	return false
}

// HandleMouseUp handles mouse up events
func (f *FlexContainer) HandleMouseUp(x, y int) bool {
	// Check if any children handle the event (in reverse order for proper z-index)
	for i := len(f.Children()) - 1; i >= 0; i-- {
		child := f.Children()[i]
		if child.HandleMouseUp(x, y) {
			return true
		}
	}
	
	return false
}

// HandleMouseMove handles mouse move events
func (f *FlexContainer) HandleMouseMove(x, y int) bool {
	// Check if any children handle the event (in reverse order for proper z-index)
	for i := len(f.Children()) - 1; i >= 0; i-- {
		child := f.Children()[i]
		if child.HandleMouseMove(x, y) {
			return true
		}
	}
	
	return false
}

// SetSpacing sets the spacing between items
func (f *FlexContainer) SetSpacing(spacing int) {
	f.spacing = spacing
	f.updateLayout()
}

// AddChild adds a child element and updates layout
func (f *FlexContainer) AddChild(child Element) {
	f.Node.AddChild(child)
	f.updateLayout()
}

// RemoveChild removes a child element and updates layout
func (f *FlexContainer) RemoveChild(child Element) {
	f.Node.RemoveChild(child)
	f.updateLayout()
}

// updateLayout updates the layout of children
func (f *FlexContainer) updateLayout() {
	if len(f.Children()) == 0 {
		return
	}
	
	bounds := f.ComputedBounds()
	boxModel := f.GetBoxModel()
	
	// Calculate content area (inside padding)
	contentX := bounds.X + boxModel.Padding.Left
	contentY := bounds.Y + boxModel.Padding.Top
	contentWidth := bounds.Width - boxModel.Padding.Left - boxModel.Padding.Right
	contentHeight := bounds.Height - boxModel.Padding.Top - boxModel.Padding.Bottom
	
	// Simplified flex layout algorithm
	if f.flexDirection == FlexRow {
		// Row layout - items side by side
		x := contentX
		for _, child := range f.Children() {
			childBounds := child.Bounds()
			childHeight := childBounds.Height
			
			// Vertical alignment
			var y int
			switch f.alignItems {
			case AlignStart:
				y = contentY
			case AlignCenter:
				y = contentY + (contentHeight - childHeight) / 2
			case AlignEnd:
				y = contentY + contentHeight - childHeight
			case AlignStretch:
				childHeight = contentHeight
				y = contentY
			}
			
			// Set child position
			child.SetBounds(Rect{x, y, childBounds.Width, childHeight})
			
			// Move to next position
			x += childBounds.Width + f.spacing
		}
	} else {
		// Column layout - items stacked
		y := contentY
		for _, child := range f.Children() {
			childBounds := child.Bounds()
			childWidth := childBounds.Width
			
			// Horizontal alignment
			var x int
			switch f.alignItems {
			case AlignStart:
				x = contentX
			case AlignCenter:
				x = contentX + (contentWidth - childWidth) / 2
			case AlignEnd:
				x = contentX + contentWidth - childWidth
			case AlignStretch:
				childWidth = contentWidth
				x = contentX
			}
			
			// Set child position
			child.SetBounds(Rect{x, y, childWidth, childBounds.Height})
			
			// Move to next position
			y += childBounds.Height + f.spacing
		}
	}
} 