package components

import (
	"image/color"
)

// TextArea represents a multi-line text input
type TextArea struct {
	*Node
	text        string
	fontSize    int
	textColor   color.RGBA
	onChange    func(string)
	focused     bool
	placeholder string
}

// NewTextArea creates a new text area
func NewTextArea(id string) *TextArea {
	return &TextArea{
		Node:        NewNode(id),
		text:        "",
		fontSize:    14,
		textColor:   color.RGBA{0, 0, 0, 255},
		onChange:    nil,
		focused:     false,
		placeholder: "",
	}
}

// SetText sets the text content
func (t *TextArea) SetText(text string) {
	t.text = text
	if t.onChange != nil {
		t.onChange(t.text)
	}
}

// GetText returns the text content
func (t *TextArea) GetText() string {
	return t.text
}

// SetFontSize sets the font size
func (t *TextArea) SetFontSize(size int) {
	t.fontSize = size
}

// SetTextColor sets the text color
func (t *TextArea) SetTextColor(color color.RGBA) {
	t.textColor = color
}

// SetOnChange sets the change handler
func (t *TextArea) SetOnChange(handler func(string)) {
	t.onChange = handler
}

// SetPlaceholder sets the placeholder text shown when the text area is empty
func (t *TextArea) SetPlaceholder(placeholder string) {
	t.placeholder = placeholder
}

// Draw draws the text area
func (t *TextArea) Draw(surface DrawSurface) {
	if !t.IsVisible() {
		return
	}
	
	bounds := t.ComputedBounds()
	
	// Draw background
	surface.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{255, 255, 255, 255})
	
	// Draw border
	surface.DrawRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{100, 100, 100, 255})
	
	// Draw text or placeholder if empty
	if t.text != "" {
		surface.DrawText(t.text, bounds.X + 5, bounds.Y + 5, t.textColor, t.fontSize)
	} else if t.placeholder != "" {
		// Draw placeholder with a lighter color
		surface.DrawText(t.placeholder, bounds.X + 5, bounds.Y + 5, color.RGBA{180, 180, 180, 255}, t.fontSize)
	}
	
	// Draw children (if any)
	for _, child := range t.Children() {
		child.Draw(surface)
	}
}

// HandleMouseDown handles mouse down events
func (t *TextArea) HandleMouseDown(x, y int) bool {
	bounds := t.ComputedBounds()
	if PointInRect(Point{x, y}, bounds) {
		t.focused = true
		return true
	} else {
		t.focused = false
	}
	
	return false
}

// Select represents a dropdown select box
type Select struct {
	*Node
	options        []string
	selectedIndex  int
	onChange       func(int)
	isOpen         bool
	backgroundColor color.RGBA
	textColor      color.RGBA
	fontSize       int
}

// NewSelect creates a new select box
func NewSelect(id string, options []string) *Select {
	return &Select{
		Node:           NewNode(id),
		options:        options,
		selectedIndex:  -1,
		onChange:       nil,
		isOpen:         false,
		backgroundColor: color.RGBA{240, 240, 240, 255},
		textColor:      color.RGBA{0, 0, 0, 255},
		fontSize:       14,
	}
}

// SetOptions sets the available options
func (s *Select) SetOptions(options []string) {
	s.options = options
	if s.selectedIndex >= len(options) {
		s.selectedIndex = -1
	}
}

// GetOptions returns the available options
func (s *Select) GetOptions() []string {
	return s.options
}

// SetSelectedIndex sets the selected option index
func (s *Select) SetSelectedIndex(index int) {
	if index >= -1 && index < len(s.options) {
		s.selectedIndex = index
		if s.onChange != nil {
			s.onChange(index)
		}
	}
}

// GetSelectedIndex returns the selected option index
func (s *Select) GetSelectedIndex() int {
	return s.selectedIndex
}

// GetSelectedOption returns the selected option text
func (s *Select) GetSelectedOption() string {
	if s.selectedIndex >= 0 && s.selectedIndex < len(s.options) {
		return s.options[s.selectedIndex]
	}
	return ""
}

// SetOnChange sets the change handler
func (s *Select) SetOnChange(handler func(int)) {
	s.onChange = handler
}

// Draw draws the select box
func (s *Select) Draw(surface DrawSurface) {
	if !s.IsVisible() {
		return
	}
	
	bounds := s.ComputedBounds()
	
	// Draw background
	surface.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, s.backgroundColor)
	
	// Draw border
	surface.DrawRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{100, 100, 100, 255})
	
	// Draw selected option or placeholder
	text := "Select..."
	if s.selectedIndex >= 0 && s.selectedIndex < len(s.options) {
		text = s.options[s.selectedIndex]
	}
	
	surface.DrawText(text, bounds.X + 5, bounds.Y + (bounds.Height - s.fontSize) / 2, s.textColor, s.fontSize)
	
	// Draw dropdown arrow
	arrowX := bounds.X + bounds.Width - 20
	arrowY := bounds.Y + bounds.Height / 2
	
	// Simple triangle
	surface.DrawLine(arrowX, arrowY - 3, arrowX + 6, arrowY + 3, s.textColor)
	surface.DrawLine(arrowX + 6, arrowY + 3, arrowX + 12, arrowY - 3, s.textColor)
	
	// If open, draw dropdown list
	if s.isOpen {
		dropdownHeight := len(s.options) * 20
		
		// Draw dropdown background
		surface.FillRect(bounds.X, bounds.Y + bounds.Height, bounds.Width, dropdownHeight, s.backgroundColor)
		
		// Draw dropdown border
		surface.DrawRect(bounds.X, bounds.Y + bounds.Height, bounds.Width, dropdownHeight, color.RGBA{100, 100, 100, 255})
		
		// Draw options
		for i, option := range s.options {
			optionY := bounds.Y + bounds.Height + i * 20
			
			// Highlight selected option
			if i == s.selectedIndex {
				surface.FillRect(bounds.X, optionY, bounds.Width, 20, color.RGBA{200, 200, 255, 255})
			}
			
			// Draw option text
			surface.DrawText(option, bounds.X + 5, optionY + 3, s.textColor, s.fontSize)
		}
	}
	
	// Draw children (if any)
	for _, child := range s.Children() {
		child.Draw(surface)
	}
}

// HandleMouseDown handles mouse down events
func (s *Select) HandleMouseDown(x, y int) bool {
	bounds := s.ComputedBounds()
	
	// Check if click is in main select box
	if PointInRect(Point{x, y}, bounds) {
		s.isOpen = !s.isOpen
		return true
	}
	
	// If open, check if click is in dropdown area
	if s.isOpen {
		dropdownHeight := len(s.options) * 20
		dropdownBounds := Rect{bounds.X, bounds.Y + bounds.Height, bounds.Width, dropdownHeight}
		
		if PointInRect(Point{x, y}, dropdownBounds) {
			// Calculate which option was clicked
			optionIndex := (y - (bounds.Y + bounds.Height)) / 20
			if optionIndex >= 0 && optionIndex < len(s.options) {
				s.SetSelectedIndex(optionIndex)
				s.isOpen = false
				return true
			}
		} else {
			// Close dropdown if click outside
			s.isOpen = false
			return true
		}
	}
	
	return false
}

// Form represents a form container with submit capability
type Form struct {
	*Node
	onSubmit func(map[string]string)
}

// NewForm creates a new form
func NewForm(id string) *Form {
	return &Form{
		Node:     NewNode(id),
		onSubmit: nil,
	}
}

// SetOnSubmit sets the submit handler
func (f *Form) SetOnSubmit(handler func(map[string]string)) {
	f.onSubmit = handler
}

// Submit submits the form, collecting values from input elements
func (f *Form) Submit() {
	if f.onSubmit == nil {
		return
	}
	
	// Collect form data from child elements
	formData := make(map[string]string)
	f.collectFormData(f, formData)
	
	// Call the submit handler
	f.onSubmit(formData)
}

// collectFormData recursively collects form data from input elements
func (f *Form) collectFormData(element Element, formData map[string]string) {
	// Check if element is a form input and get its value
	if input, ok := element.(*TextArea); ok {
		formData[input.ID()] = input.GetText()
	} else if checkbox, ok := element.(*Checkbox); ok {
		if checkbox.IsChecked() {
			formData[checkbox.ID()] = "true"
		} else {
			formData[checkbox.ID()] = "false"
		}
	} else if select_, ok := element.(*Select); ok {
		formData[select_.ID()] = select_.GetSelectedOption()
	}
	
	// Recursively process children
	for _, child := range element.Children() {
		f.collectFormData(child, formData)
	}
}

// Draw draws the form
func (f *Form) Draw(surface DrawSurface) {
	if !f.IsVisible() {
		return
	}
	
	// Draw children (inputs, buttons, etc.)
	for _, child := range f.Children() {
		child.Draw(surface)
	}
}

// HandleMouseDown handles mouse down events
func (f *Form) HandleMouseDown(x, y int) bool {
	// Check if any children handle the event
	for i := len(f.Children()) - 1; i >= 0; i-- {
		child := f.Children()[i]
		if child.HandleMouseDown(x, y) {
			return true
		}
	}
	
	return false
} 