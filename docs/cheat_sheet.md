# Finch UI Cheat Sheet

## Installation & Import

```go
go get github.com/vijay/finch-ui
```

```go
import "github.com/vijay/finch-ui/finch"
```

## Basic Structure

```go
// Create a new UI
ui := finch.New()

// Configure the page
ui.SetPageConfig("My App", "column")

// Add components and handle events
// ...

// Run the application
ui.Run(800, 600)
```

## Display Text

```go
// Add a title
ui.Title("My Title")

// Add standard text
ui.Text("Simple text")

// Chain style methods
ui.Title("Styled Title").Centered().Size(24)

// Set text programmatically
text := ui.Text("Initial text")
text.SetText("Updated text")
```

## Interactive Widgets

```go
// Button
ui.Button("Click Me").OnClick(func() {
    // Handle click
})

// Text input
input := ui.TextInput("Placeholder text")
text := input.Value()  // Get current value
input.Clear()          // Clear the input

// Checkbox
checkbox := ui.Checkbox("Enable feature")
isChecked := checkbox.Value()
checkbox.SetValue(true)
checkbox.OnChange(func(checked bool) {
    // Handle change
})
```

## Layout

```go
// Container with row layout
ui.Container().Layout("row", func(c *finch.Container) {
    c.Text("Left")
    c.Text("Right")
})

// Container with column layout
ui.Container().Layout("column", func(c *finch.Container) {
    c.Text("Top")
    c.Text("Bottom")
})

// Container with styling
ui.Container().
    Background("#ffffff").
    Padding(10).
    Margin(5).
    Width(400).
    Height(300)

// Multi-column layout
ui.Columns(3, func(cols []*finch.Column) {
    cols[0].Text("Column 1")
    cols[1].Text("Column 2")
    cols[2].Text("Column 3")
})

// Tabs
ui.Tabs([]string{"Info", "Settings"}, func(tabs []*finch.Tab) {
    tabs[0].Text("Info content")
    tabs[1].Button("Save settings")
})
```

## State Management

```go
// Create state
counter := ui.State(0)

// Update state with a transform function
counter.Update(func(value interface{}) interface{} {
    return value.(int) + 1
})

// Watch for state changes
counter.Watch(func(value interface{}) {
    count := value.(int)
    fmt.Println("Counter:", count)
})

// Get current state value
currentValue := counter.Value().(int)
```

## Building a Todo App

```go
// Define data structure
type TodoItem struct {
    ID   string
    Text string
    Done bool
}

// Create state for todo items
todos := ui.State([]TodoItem{})

// Add a new todo
todos.Update(func(value interface{}) interface{} {
    items := value.([]TodoItem)
    newItem := TodoItem{
        ID:   fmt.Sprintf("todo_%d", len(items)+1),
        Text: "New todo item",
        Done: false,
    }
    return append(items, newItem)
})

// Remove a todo
todos.Update(func(value interface{}) interface{} {
    items := value.([]TodoItem)
    result := []TodoItem{}
    for _, item := range items {
        if item.ID != "todo_to_remove" {
            result = append(result, item)
        }
    }
    return result
})

// Watch for changes to todos
todos.Watch(func(value interface{}) {
    items := value.([]TodoItem)
    for _, item := range items {
        // Update UI
    }
})
```

## Common Patterns

### Create a Counter

```go
// Create state
counter := ui.State(0)

// Display counter
display := ui.Text("0")

// Create increment/decrement buttons
ui.Container().Layout("row", func(c *finch.Container) {
    c.Button("-").OnClick(func() {
        counter.Update(func(v interface{}) interface{} {
            return v.(int) - 1
        })
    })
    
    c.Button("+").OnClick(func() {
        counter.Update(func(v interface{}) interface{} {
            return v.(int) + 1
        })
    })
})

// Update display when counter changes
counter.Watch(func(v interface{}) {
    display.SetText(fmt.Sprintf("%d", v.(int)))
})
```

### Form with Validation

```go
// Create state for form validity
isValid := ui.State(false)

// Create input
input := ui.TextInput("Enter your name")

// Validate on change
input.OnChange(func(text string) {
    isValid.Update(func(_ interface{}) interface{} {
        return len(text) >= 3
    })
})

// Create submit button that's enabled only when valid
submitBtn := ui.Button("Submit")

// Update button state when validity changes
isValid.Watch(func(v interface{}) {
    valid := v.(bool)
    if valid {
        submitBtn.SetEnabled(true)
    } else {
        submitBtn.SetEnabled(false)
    }
})
``` 