package components

import (
	"image/color"
)

// Text represents a simple text element in the UI
type Text struct {
	*Node
	text      string
	fontSize  int
	textColor color.RGBA
	bold      bool
	italic    bool
	alignment TextAlignment
}

// NewText creates a new text element
func NewText(id string, text string, fontSize int, textColor color.RGBA) *Text {
	return &Text{
		Node:      NewNode(id),
		text:      text,
		fontSize:  fontSize,
		textColor: textColor,
		bold:      false,
		italic:    false,
		alignment: TextAlignLeft,
	}
}

// SetText sets the text content
func (t *Text) SetText(text string) {
	t.text = text
}

// GetText returns the text content
func (t *Text) GetText() string {
	return t.text
}

// SetFontSize sets the font size
func (t *Text) SetFontSize(size int) {
	t.fontSize = size
}

// SetTextColor sets the text color
func (t *Text) SetTextColor(color color.RGBA) {
	t.textColor = color
}

// SetBold sets whether the text is bold
func (t *Text) SetBold(bold bool) {
	t.bold = bold
}

// SetItalic sets whether the text is italic
func (t *Text) SetItalic(italic bool) {
	t.italic = italic
}

// SetAlignment sets the text alignment
func (t *Text) SetAlignment(alignment TextAlignment) {
	t.alignment = alignment
}

// Draw draws the text
func (t *Text) Draw(surface DrawSurface) {
	if !t.IsVisible() {
		return
	}
	
	bounds := t.ComputedBounds()
	
	// Calculate text position based on alignment
	textWidth := len(t.text) * t.fontSize / 2
	textX := bounds.X
	
	if t.alignment == TextAlignCenter {
		textX = bounds.X + (bounds.Width - textWidth) / 2
	} else if t.alignment == TextAlignRight {
		textX = bounds.X + bounds.Width - textWidth
	}
	
	// Draw the text
	surface.DrawText(t.text, textX, bounds.Y, t.textColor, t.fontSize)
	
	// Draw children (if any)
	for _, child := range t.Children() {
		child.Draw(surface)
	}
}

// HandleMouseDown handles mouse down events
func (t *Text) HandleMouseDown(x, y int) bool {
	// Text doesn't handle mouse events directly, but we check children
	for i := len(t.Children()) - 1; i >= 0; i-- {
		child := t.Children()[i]
		if child.HandleMouseDown(x, y) {
			return true
		}
	}
	return false
}

// HandleMouseMove handles mouse move events
func (t *Text) HandleMouseMove(x, y int) bool {
	// Text doesn't handle mouse events directly, but we check children
	for i := len(t.Children()) - 1; i >= 0; i-- {
		child := t.Children()[i]
		if child.HandleMouseMove(x, y) {
			return true
		}
	}
	return false
}

// Label represents a label element in the UI
type Label struct {
	*Node
	text      string
	fontSize  int
	textColor color.RGBA
	bold      bool
	italic    bool
	alignment TextAlignment
}

// NewLabel creates a new label
func NewLabel(id string, text string, fontSize int, textColor color.RGBA) *Label {
	return &Label{
		Node:      NewNode(id),
		text:      text,
		fontSize:  fontSize,
		textColor: textColor,
		bold:      false,
		italic:    false,
		alignment: TextAlignLeft,
	}
}

// SetText sets the label text
func (l *Label) SetText(text string) {
	l.text = text
}

// GetText returns the label text
func (l *Label) GetText() string {
	return l.text
}

// SetFontSize sets the font size
func (l *Label) SetFontSize(size int) {
	l.fontSize = size
}

// SetTextColor sets the text color
func (l *Label) SetTextColor(color color.RGBA) {
	l.textColor = color
}

// SetBold sets whether the text is bold
func (l *Label) SetBold(bold bool) {
	l.bold = bold
}

// SetItalic sets whether the text is italic
func (l *Label) SetItalic(italic bool) {
	l.italic = italic
}

// SetAlignment sets the text alignment
func (l *Label) SetAlignment(alignment TextAlignment) {
	l.alignment = alignment
}

// SetTextAlignment is an alias for SetAlignment
func (l *Label) SetTextAlignment(alignment TextAlignment) {
	l.SetAlignment(alignment)
}

// Draw draws the label
func (l *Label) Draw(surface DrawSurface) {
	if !l.IsVisible() {
		return
	}
	
	bounds := l.ComputedBounds()
	
	// Calculate text position based on alignment
	textWidth := len(l.text) * l.fontSize / 2
	textX := bounds.X
	
	if l.alignment == TextAlignCenter {
		textX = bounds.X + (bounds.Width - textWidth) / 2
	} else if l.alignment == TextAlignRight {
		textX = bounds.X + bounds.Width - textWidth
	}
	
	// Center text vertically in the label
	textY := bounds.Y + (bounds.Height - l.fontSize) / 2
	
	// Draw the text
	surface.DrawText(l.text, textX, textY, l.textColor, l.fontSize)
	
	// Draw children (if any)
	for _, child := range l.Children() {
		child.Draw(surface)
	}
}

// HandleMouseDown handles mouse down events
func (l *Label) HandleMouseDown(x, y int) bool {
	// Label doesn't handle mouse events directly, but we check children
	for i := len(l.Children()) - 1; i >= 0; i-- {
		child := l.Children()[i]
		if child.HandleMouseDown(x, y) {
			return true
		}
	}
	return false
}

// HandleMouseMove handles mouse move events
func (l *Label) HandleMouseMove(x, y int) bool {
	// Label doesn't handle mouse events directly, but we check children
	for i := len(l.Children()) - 1; i >= 0; i-- {
		child := l.Children()[i]
		if child.HandleMouseMove(x, y) {
			return true
		}
	}
	return false
} 