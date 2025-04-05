package main

import (
	"fmt"

	"github.com/vijay/finch-ui/finch"
)

// TodoItem represents a single todo item with its data
type TodoItem struct {
	ID   string
	Text string
	Done bool
}

func main() {
	// Create a new UI instance
	ui := finch.New()
	
	// Configure the page
	ui.SetPageConfig("Todo List App", "column")
	
	// Create reactive state
	todos := ui.State([]TodoItem{
		{ID: "todo_1", Text: "Buy groceries", Done: false},
		{ID: "todo_2", Text: "Finish project", Done: false},
	})
	
	// Add a title
	ui.Title("Todo List").Centered().Size(24)
	
	// Create input area
	ui.Container().Layout("row", func(c *finch.Container) {
		// Add a text input
		input := c.TextInput("Enter a new todo...")
		
		// Add a button
		c.Button("Add").OnClick(func() {
			// Get text from the input
			text := input.Value()
			
			if text != "" {
				// Add a new todo
				todos.Update(func(value interface{}) interface{} {
					items := value.([]TodoItem)
					newItem := TodoItem{
						ID:   fmt.Sprintf("todo_%d", len(items)+1),
						Text: text,
						Done: false,
					}
					return append(items, newItem)
				})
				
				// Clear the input
				input.Clear()
			}
		})
	})
	
	// Create todo list container
	listContainer := ui.Container().Background("#ffffff").Margin(10).Padding(10).Grow(1)
	
	// Create status text
	statusText := ui.Text("")
	
	// Create action buttons
	ui.Container().Layout("row", func(c *finch.Container) {
		// Add a clear completed button
		c.Button("Clear Completed").OnClick(func() {
			todos.Update(func(value interface{}) interface{} {
				items := value.([]TodoItem)
				result := []TodoItem{}
				
				// Filter out completed todos
				for _, item := range items {
					if !item.Done {
						result = append(result, item)
					}
				}
				
				return result
			})
		})
	})
	
	// Watch for changes to todos
	todos.Watch(func(value interface{}) {
		items := value.([]TodoItem)
		
		// Update the list container with the current todos
		listContainer.RemoveAllChildren()
		
		// Add each todo as a row
		for i, item := range items {
			index := i // Capture the index for closure
			
			// Create a row for this todo
			listContainer.Layout("row", func(row *finch.Container) {
				// Add checkbox
				checkbox := row.Checkbox(item.Text)
				checkbox.SetValue(item.Done)
				
				// Handle checkbox changes
				checkbox.OnChange(func(checked bool) {
					todos.Update(func(value interface{}) interface{} {
						items := value.([]TodoItem)
						items[index].Done = checked
						return items
					})
				})
				
				// Add delete button
				row.Button("×").OnClick(func() {
					todos.Update(func(value interface{}) interface{} {
						items := value.([]TodoItem)
						
						// Remove the item at this index
						return append(items[:index], items[index+1:]...)
					})
				})
			})
		}
		
		// Update status text
		completed := 0
		for _, item := range items {
			if item.Done {
				completed++
			}
		}
		
		if len(items) == 0 {
			statusText.SetText("No items")
		} else if len(items) == 1 {
			statusText.SetText("1 item")
		} else {
			statusText.SetText(fmt.Sprintf("%d items, %d completed", len(items), completed))
		}
	})
	
	// Run the UI
	ui.Run(800, 600)
} 