package components

import (
	"strings"
)

// Node provides a base implementation of the NodeElement interface
type Node struct {
	*BaseElement
	positionType    PositionType
	boxModel        BoxModel
	relativePos     Point
	flexDirection   FlexDirection
	alignItems      Alignment
	justifyContent  Alignment
	classNames      []string
	visible         bool
}

// NewNode creates a new node
func NewNode(id string) *Node {
	return &Node{
		BaseElement:    NewBaseElement(id),
		positionType:   PositionRelative,
		boxModel:       BoxModel{},
		relativePos:    Point{0, 0},
		flexDirection:  FlexRow,
		alignItems:     AlignStart,
		justifyContent: AlignStart,
		classNames:     make([]string, 0),
		visible:        true,
	}
}

// GetPositionType returns the position type
func (d *Node) GetPositionType() PositionType {
	return d.positionType
}

// SetPositionType sets the position type
func (d *Node) SetPositionType(posType PositionType) {
	d.positionType = posType
}

// GetBoxModel returns the box model
func (d *Node) GetBoxModel() BoxModel {
	return d.boxModel
}

// SetBoxModel sets the box model
func (d *Node) SetBoxModel(box BoxModel) {
	d.boxModel = box
}

// GetRelativePosition returns the relative position
func (d *Node) GetRelativePosition() Point {
	return d.relativePos
}

// SetRelativePosition sets the relative position
func (d *Node) SetRelativePosition(pos Point) {
	d.relativePos = pos
}

// IsVisible returns whether the element is visible
func (d *Node) IsVisible() bool {
	return d.visible
}

// SetVisible sets whether the element is visible
func (d *Node) SetVisible(visible bool) {
	d.visible = visible
}

// AddClass adds a class name to the element
func (d *Node) AddClass(className string) {
	// Don't add duplicate class names
	if !d.HasClass(className) {
		d.classNames = append(d.classNames, className)
	}
}

// RemoveClass removes a class name from the element
func (d *Node) RemoveClass(className string) {
	for i, c := range d.classNames {
		if c == className {
			d.classNames = append(d.classNames[:i], d.classNames[i+1:]...)
			break
		}
	}
}

// HasClass checks if element has a class
func (d *Node) HasClass(className string) bool {
	for _, c := range d.classNames {
		if c == className {
			return true
		}
	}
	return false
}

// GetClassNames returns all class names
func (d *Node) GetClassNames() []string {
	return d.classNames
}

// ComputedBounds calculates and returns the absolute screen position
func (d *Node) ComputedBounds() Rect {
	var bounds Rect
	
	// Start with the element's own bounds
	bounds = d.Bounds()
	
	// If positioned absolutely or relatively, adjust based on parent
	if d.positionType != PositionFixed && d.Parent() != nil {
		// Get parent's content area (without considering margins)
		var parentBounds Rect
		
		// If parent is a DOM element, use its computed bounds
		if domParent, ok := d.Parent().(NodeElement); ok {
			parentBounds = domParent.ComputedBounds()
			
			// Apply parent's padding to get the content area
			parentBoxModel := domParent.GetBoxModel()
			parentBounds.X += parentBoxModel.Padding.Left
			parentBounds.Y += parentBoxModel.Padding.Top
			parentBounds.Width -= (parentBoxModel.Padding.Left + parentBoxModel.Padding.Right)
			parentBounds.Height -= (parentBoxModel.Padding.Top + parentBoxModel.Padding.Bottom)
		} else {
			// For non-DOM parents, just use their bounds
			parentBounds = d.Parent().Bounds()
		}
		
		// For relative positioning, add relative offset to parent's content area
		if d.positionType == PositionRelative {
			bounds.X = parentBounds.X + d.relativePos.X
			bounds.Y = parentBounds.Y + d.relativePos.Y
		} else if d.positionType == PositionAbsolute {
			// For absolute positioning, position relative to parent's bounds
			bounds.X = parentBounds.X + d.relativePos.X
			bounds.Y = parentBounds.Y + d.relativePos.Y
		}
	}
	
	// Apply margin (affects position but not size)
	bounds.X += d.boxModel.Margin.Left
	bounds.Y += d.boxModel.Margin.Top
	
	// Return the computed bounds
	return bounds
}

// Draw draws the element and its children
func (d *Node) Draw(surface DrawSurface) {
	// If not visible, don't draw
	if !d.visible {
		return
	}
	
	// Get the computed bounds
	bounds := d.ComputedBounds()
	
	// Draw borders if they exist
	if d.boxModel.Border.Style != BorderNone {
		borderColor := d.boxModel.Border.Color
		
		// Top border
		if d.boxModel.Border.Width.Top > 0 {
			surface.FillRect(
				bounds.X, 
				bounds.Y, 
				bounds.Width, 
				d.boxModel.Border.Width.Top, 
				borderColor)
		}
		
		// Right border
		if d.boxModel.Border.Width.Right > 0 {
			surface.FillRect(
				bounds.X + bounds.Width - d.boxModel.Border.Width.Right, 
				bounds.Y, 
				d.boxModel.Border.Width.Right, 
				bounds.Height, 
				borderColor)
		}
		
		// Bottom border
		if d.boxModel.Border.Width.Bottom > 0 {
			surface.FillRect(
				bounds.X, 
				bounds.Y + bounds.Height - d.boxModel.Border.Width.Bottom, 
				bounds.Width, 
				d.boxModel.Border.Width.Bottom, 
				borderColor)
		}
		
		// Left border
		if d.boxModel.Border.Width.Left > 0 {
			surface.FillRect(
				bounds.X, 
				bounds.Y, 
				d.boxModel.Border.Width.Left, 
				bounds.Height, 
				borderColor)
		}
	}
	
	// Draw all children
	for _, child := range d.Children() {
		child.Draw(surface)
	}
}

// QuerySelector finds the first element matching the selector
func (d *Node) QuerySelector(selector string) NodeElement {
	// Simple selector implementation. In a full implementation, this would be more robust.
	// Currently supports:
	// - ID selectors: #id
	// - Class selectors: .class
	// - Tag/type selectors: tag
	
	selectorType, selectorValue := parseSelectorString(selector)
	
	// Check if this element matches
	if selectorType == "id" && d.ID() == selectorValue {
		return d
	} else if selectorType == "class" {
		for _, class := range d.classNames {
			if class == selectorValue {
				return d
			}
		}
	} else if selectorType == "tag" && strings.Contains(d.ID(), selectorValue) {
		// Simple tag selector implementation
		return d
	}
	
	// If not this element, search children
	for _, child := range d.Children() {
		if domChild, ok := child.(NodeElement); ok {
			if result := domChild.QuerySelector(selector); result != nil {
				return result
			}
		}
	}
	
	return nil
}

// QuerySelectorAll finds all elements matching the selector
func (d *Node) QuerySelectorAll(selector string) []NodeElement {
	results := make([]NodeElement, 0)
	
	selectorType, selectorValue := parseSelectorString(selector)
	
	// Check if this element matches
	if selectorType == "id" && d.ID() == selectorValue {
		results = append(results, d)
	} else if selectorType == "class" {
		for _, class := range d.classNames {
			if class == selectorValue {
				results = append(results, d)
				break
			}
		}
	} else if selectorType == "tag" && strings.Contains(d.ID(), selectorValue) {
		results = append(results, d)
	}
	
	// Search children
	for _, child := range d.Children() {
		if domChild, ok := child.(NodeElement); ok {
			childResults := domChild.QuerySelectorAll(selector)
			results = append(results, childResults...)
		}
	}
	
	return results
} 