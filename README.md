# Finch UI Framework

Finch UI is a DOM-based GUI framework for Go applications that makes it easy to build beautiful, interactive user interfaces with minimal code.

## Key Features

- **Component-Based Architecture**: Create reusable, composable UI components
- **Intuitive API**: Streamlit-inspired declarative API for effortless development
- **Reactive State Management**: Automatic UI updates when state changes
- **Modern Layouts**: Flexible containers, columns, and tabs
- **Built on Ebiten**: Powered by the popular Go game engine for smooth rendering

## Getting Started

### Basic Setup

```go
// Create a new UI instance
ui := finch.New()

// Configure the UI
ui.SetPageConfig("My App", "column")

// Add components
ui.Title("Hello, Finch UI!").Centered()
ui.Text("This is a simple app built with Finch UI.")

// Add interactive elements
ui.Button("Click Me").OnClick(func() {
    fmt.Println("Button clicked!")
})

// Run the UI
ui.Run(800, 600)
```

### Managing State

Finch UI includes a reactive state system that automatically updates the UI when state changes:

```go
// Create a counter state
counter := ui.State(0)

// Display the counter
display := ui.Text("0")

// Create buttons to modify the state
ui.Button("+").OnClick(func() {
    counter.Update(func(value interface{}) interface{} {
        return value.(int) + 1
    })
})

// Watch for changes to update the UI
counter.Watch(func(value interface{}) {
    count := value.(int)
    display.SetText(fmt.Sprintf("%d", count))
})
```

### Layout Options

Finch UI provides flexible layout options:

```go
// Create a row layout
ui.Container().Layout("row", func(c *finch.Container) {
    c.Text("Left")
    c.Text("Center")
    c.Text("Right")
})

// Create columns with specific content
ui.Columns(3, func(cols []*finch.Column) {
    cols[0].Text("Column 1")
    cols[1].Text("Column 2")
    cols[2].Text("Column 3")
})

// Create tabs
ui.Tabs([]string{"Tab 1", "Tab 2"}, func(tabs []*finch.Tab) {
    tabs[0].Text("This is tab 1 content")
    tabs[1].Text("This is tab 2 content")
})
```

## Examples

### Todo App Example

A complete Todo application in under 100 lines of code:

```go
func main() {
    // Create a new UI instance
    ui := finch.New()
    
    // Configure the page
    ui.SetPageConfig("Todo List App", "column")
    
    // Create reactive state
    todos := ui.State([]TodoItem{})
    
    // Add a title
    ui.Title("Todo List").Centered()
    
    // Create input area
    ui.Container().Layout("row", func(c *finch.Container) {
        input := c.TextInput("Enter a new todo...")
        c.Button("Add").OnClick(func() {
            if text := input.Value(); text != "" {
                todos.Update(func(value interface{}) interface{} {
                    items := value.([]TodoItem)
                    newItem := TodoItem{ID: generateID(), Text: text}
                    return append(items, newItem)
                })
                input.Clear()
            }
        })
    })
    
    // Watch for changes to todos
    todos.Watch(func(value interface{}) {
        // Update UI with current todo items
    })
    
    // Run the UI
    ui.Run(800, 600)
}
```

## Design Philosophy

Finch UI is designed with inspiration from Streamlit, focusing on:

1. **Code Efficiency**: Create powerful UIs with minimal code
2. **Developer Experience**: Intuitive API that's easy to learn and use
3. **Reactivity**: Automatic UI updates when data changes
4. **Composition**: Build complex UIs from simple building blocks

## Project Structure

- **core/**: Core types and base functionality
- **components/**: Base UI components with DOM-like interface
- **finch/**: High-level, Streamlit-inspired API
- **examples/**: Example applications showcasing the framework

## License

MIT License 