package components

import (
	"image/color"
	"strings"
)

// PositionType defines how an element is positioned
type PositionType int

const (
	PositionRelative PositionType = iota // Positioned relative to parent's content area
	PositionAbsolute                     // Positioned relative to nearest positioned ancestor
	PositionFixed                        // Positioned relative to the viewport
)

// FlexDirection defines the direction of flex items
type FlexDirection int

const (
	FlexRow FlexDirection = iota
	FlexColumn
)

// Alignment defines how items align on the cross axis
type Alignment int

const (
	AlignStart Alignment = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

// BoxModel represents the CSS-like box model for an element
type BoxModel struct {
	Margin  Spacing
	Padding Spacing
	Border  Border
}

// Spacing represents spacing values for top, right, bottom, left
type Spacing struct {
	Top, Right, Bottom, Left int
}

// Border represents border properties
type Border struct {
	Width Spacing
	Color color.RGBA
	Style BorderStyle
}

// BorderStyle defines the style of a border
type BorderStyle int

const (
	BorderNone BorderStyle = iota
	BorderSolid
	BorderDashed
	BorderDotted
)

// TextAlignment defines text alignment options
type TextAlignment int

const (
	TextAlignLeft TextAlignment = iota
	TextAlignCenter
	TextAlignRight
)

// NodeElement extends the base Element interface with DOM-like capabilities
type NodeElement interface {
	Element
	
	// DOM-specific methods
	GetPositionType() PositionType
	SetPositionType(posType PositionType)
	GetBoxModel() BoxModel
	SetBoxModel(box BoxModel)
	GetRelativePosition() Point
	SetRelativePosition(pos Point)
	QuerySelector(selector string) NodeElement
	QuerySelectorAll(selector string) []NodeElement
	
	// Class management
	AddClass(className string)
	RemoveClass(className string)
	HasClass(className string) bool
	GetClassNames() []string
	
	// Computed values
	ComputedBounds() Rect // Returns the absolute screen position after calculations
}

// parseSelectorString parses a selector string and returns its type and value
func parseSelectorString(selector string) (string, string) {
	if selector == "" {
		return "", ""
	}
	
	if strings.HasPrefix(selector, "#") {
		return "id", selector[1:]
	} else if strings.HasPrefix(selector, ".") {
		return "class", selector[1:]
	} else if strings.HasPrefix(selector, "[") && strings.HasSuffix(selector, "]") {
		// Attribute selector (simplified)
		return "attr", selector[1:len(selector)-1]
	} else {
		return "tag", selector
	}
} 