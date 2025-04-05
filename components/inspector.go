package components

import (
	"fmt"
	"image/color"
	"strings"
)

// Inspector is a debugging tool for inspecting the UI element tree
type Inspector struct {
	*Node
	root         NodeElement
	selectedNode NodeElement
	expanded     map[NodeElement]bool
	onSelect     func(NodeElement)
}

// NewInspector creates a new inspector for the given root element
func NewInspector(id string, root NodeElement) *Inspector {
	return &Inspector{
		Node:         NewNode(id),
		root:         root,
		selectedNode: nil,
		expanded:     make(map[NodeElement]bool),
		onSelect:     nil,
	}
}

// SetOnSelect sets the handler for when a node is selected
func (i *Inspector) SetOnSelect(handler func(NodeElement)) {
	i.onSelect = handler
}

// Draw draws the inspector
func (i *Inspector) Draw(surface DrawSurface) {
	if !i.IsVisible() {
		return
	}
	
	bounds := i.ComputedBounds()
	
	// Draw background
	surface.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{240, 240, 240, 255})
	
	// Draw border
	surface.DrawRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{180, 180, 180, 255})
	
	// Draw title
	surface.DrawText("UI Inspector", bounds.X + 5, bounds.Y + 5, color.RGBA{0, 0, 0, 255}, 16)
	
	// Draw element tree
	i.drawNode(surface, i.root, bounds.X + 10, bounds.Y + 30, 0)
}

// drawNode recursively draws a node and its children
func (i *Inspector) drawNode(surface DrawSurface, node NodeElement, x, y int, depth int) int {
	indent := depth * 15
	lineHeight := 20
	
	// Check if this node is selected
	isSelected := node == i.selectedNode
	
	// Determine node label
	label := fmt.Sprintf("%s (%T)", node.ID(), node)
	
	// Trim the label to fit
	if len(label) > 30 {
		label = label[:27] + "..."
	}
	
	// Draw selection highlight if selected
	if isSelected {
		surface.FillRect(x - 5, y, 200, lineHeight, color.RGBA{200, 200, 255, 255})
	}
	
	// Draw expand/collapse indicator if the node has children
	if len(node.Children()) > 0 {
		if i.expanded[node] {
			surface.DrawText("-", x, y, color.RGBA{0, 0, 0, 255}, 14)
		} else {
			surface.DrawText("+", x, y, color.RGBA{0, 0, 0, 255}, 14)
		}
	}
	
	// Draw node label
	surface.DrawText(label, x + 15, y, color.RGBA{0, 0, 0, 255}, 14)
	
	// Move to next line
	y += lineHeight
	
	// If expanded and has children, draw children
	if i.expanded[node] {
		for _, child := range node.Children() {
			// Skip non-DOM elements
			if domChild, ok := child.(NodeElement); ok {
				y = i.drawNode(surface, domChild, x + indent, y, depth + 1)
			}
		}
	}
	
	return y
}

// HandleMouseDown handles mouse down events
func (i *Inspector) HandleMouseDown(x, y int) bool {
	bounds := i.ComputedBounds()
	if !PointInRect(Point{x, y}, bounds) {
		return false
	}
	
	// Handle clicks on tree nodes
	nodeY := bounds.Y + 30
	i.handleNodeClick(i.root, bounds.X + 10, &nodeY, 0, x, y)
	
	return true
}

// handleNodeClick recursively handles clicks on tree nodes
func (i *Inspector) handleNodeClick(node NodeElement, x int, yPtr *int, depth int, clickX, clickY int) bool {
	indent := depth * 15
	lineHeight := 20
	y := *yPtr
	
	// Check if click is on this node
	if clickY >= y && clickY < y + lineHeight {
		if len(node.Children()) > 0 && clickX >= x && clickX < x + 15 {
			// Toggle expanded state
			i.expanded[node] = !i.expanded[node]
			return true
		} else if clickX >= x + 15 && clickX < x + 200 {
			// Select node
			i.selectedNode = node
			if i.onSelect != nil {
				i.onSelect(node)
			}
			return true
		}
	}
	
	// Move to next line
	*yPtr += lineHeight
	
	// If expanded and has children, check children
	if i.expanded[node] {
		for _, child := range node.Children() {
			// Skip non-DOM elements
			if domChild, ok := child.(NodeElement); ok {
				if i.handleNodeClick(domChild, x + indent, yPtr, depth + 1, clickX, clickY) {
					return true
				}
			}
		}
	}
	
	return false
}

// HighlightNode temporarily highlights a node in the UI
func (i *Inspector) HighlightNode(node NodeElement) {
	// In a real implementation, this would draw a highlight around the node
	// and scroll the inspector to show it
}

// DumpNodeTree returns a string representation of the node tree
func (i *Inspector) DumpNodeTree() string {
	var sb strings.Builder
	i.dumpNodeTreeRecursive(&sb, i.root, 0)
	return sb.String()
}

// dumpNodeTreeRecursive builds a string representation of the node tree
func (i *Inspector) dumpNodeTreeRecursive(sb *strings.Builder, node NodeElement, depth int) {
	// Add indentation
	for j := 0; j < depth; j++ {
		sb.WriteString("  ")
	}
	
	// Add node info
	sb.WriteString(fmt.Sprintf("%s (%T)\n", node.ID(), node))
	
	// Recursively process children
	for _, child := range node.Children() {
		if domChild, ok := child.(NodeElement); ok {
			i.dumpNodeTreeRecursive(sb, domChild, depth + 1)
		}
	}
} 