package test

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/aggnr/finch/components"
)

// TodoTestCase represents a test case specifically for the Todo component
type TodoTestCase struct {
	Name        string
	Description string
	Actions     []TodoTestAction
	Results     []TodoTestResult
}

// TodoTestAction represents a specific action in a Todo test case
type TodoTestAction struct {
	Type        string // "add", "toggle", "delete", "clearCompleted", "setText"
	TodoID      string // ID of the todo to act on (empty for new todos)
	TodoText    string // Text for new todos or to update existing
	Description string // Human-readable description
	Delay       time.Duration
}

// TodoTestResult represents the expected result of a Todo test case
type TodoTestResult struct {
	Type        string // "count", "completed", "text", "exists"
	TodoID      string // ID of the todo to check (empty for count checks)
	ExpectedVal interface{} // Expected value (bool, string, or int)
	Description string
}

// TodoTestRunner manages running Todo component tests
type TodoTestRunner struct {
	*components.FlexContainer
	rootContainer *components.FlexContainer
	todoList      *components.TodoList
	inputField    *components.TextArea
	addButton     *components.Button
	testCases     []*TodoTestCase
	currentTest   int
	isRunning     bool
	isInteractive bool
	stepMode      bool
	currentStep   int
	lastStepTime  time.Time
	statusLabel   *components.Label
	resultsLabel  *components.Label
	nextButton    *components.Button
	runButton     *components.Button
	stopButton    *components.Button
}

// NewTodoTestRunner creates a new Todo test runner
func NewTodoTestRunner(todoList *components.TodoList, inputField *components.TextArea, addButton *components.Button) *TodoTestRunner {
	runner := &TodoTestRunner{
		FlexContainer:  components.NewFlexContainer("todo_test_runner"),
		todoList:       todoList,
		inputField:     inputField,
		addButton:      addButton,
		testCases:      make([]*TodoTestCase, 0),
		currentTest:    0,
		isRunning:      false,
		isInteractive:  true,
		stepMode:       true,
		currentStep:    -1,
		lastStepTime:   time.Time{},
	}

	// Set up the UI for the test runner
	runner.SetBounds(components.Rect{
		X:      0,
		Y:      0,
		Width:  800,
		Height: 120,
	})
	runner.SetBackgroundColor(color.RGBA{230, 230, 230, 200})
	runner.SetBoxModel(components.BoxModel{
		Padding: components.Spacing{10, 10, 10, 10},
	})

	// Create status label
	runner.statusLabel = components.NewLabel("test_status", "Todo Test Runner Ready", 14, color.RGBA{0, 0, 0, 255})
	runner.statusLabel.SetBounds(components.Rect{X: 10, Y: 10, Width: 780, Height: 20})
	runner.AddChild(runner.statusLabel)

	// Create results label
	runner.resultsLabel = components.NewLabel("test_results", "", 14, color.RGBA{0, 100, 0, 255})
	runner.resultsLabel.SetBounds(components.Rect{X: 10, Y: 40, Width: 780, Height: 20})
	runner.AddChild(runner.resultsLabel)

	// Create buttons container
	buttonsContainer := components.NewFlexContainer("buttons_container")
	buttonsContainer.SetBounds(components.Rect{X: 10, Y: 70, Width: 780, Height: 40})
	buttonsContainer.SetFlexDirection(components.FlexRow)
	runner.AddChild(buttonsContainer)

	// Create Next button
	runner.nextButton = components.NewButton("next_button", "Next Step")
	runner.nextButton.SetBounds(components.Rect{X: 0, Y: 0, Width: 120, Height: 30})
	runner.nextButton.SetOnClick(func() {
		if runner.isRunning && runner.stepMode {
			runner.ExecuteNextStep()
		}
	})
	buttonsContainer.AddChild(runner.nextButton)

	// Create Run button
	runner.runButton = components.NewButton("run_button", "Run Test")
	runner.runButton.SetBounds(components.Rect{X: 130, Y: 0, Width: 120, Height: 30})
	runner.runButton.SetOnClick(func() {
		if !runner.isRunning && len(runner.testCases) > 0 {
			runner.StartTest(runner.currentTest)
		}
	})
	buttonsContainer.AddChild(runner.runButton)

	// Create Stop button
	runner.stopButton = components.NewButton("stop_button", "Stop Test")
	runner.stopButton.SetBounds(components.Rect{X: 260, Y: 0, Width: 120, Height: 30})
	runner.stopButton.SetOnClick(func() {
		if runner.isRunning {
			runner.StopTest()
		}
	})
	buttonsContainer.AddChild(runner.stopButton)

	return runner
}

// NewTodoTestCase creates a new Todo test case
func NewTodoTestCase(name, description string) *TodoTestCase {
	return &TodoTestCase{
		Name:        name,
		Description: description,
		Actions:     make([]TodoTestAction, 0),
		Results:     make([]TodoTestResult, 0),
	}
}

// AddAddTodoAction adds an action to add a new todo
func (tc *TodoTestCase) AddAddTodoAction(text, description string, delay time.Duration) {
	tc.Actions = append(tc.Actions, TodoTestAction{
		Type:        "add",
		TodoText:    text,
		Description: description,
		Delay:       delay,
	})
}

// AddToggleTodoAction adds an action to toggle a todo's completion status
func (tc *TodoTestCase) AddToggleTodoAction(todoID, description string, delay time.Duration) {
	tc.Actions = append(tc.Actions, TodoTestAction{
		Type:        "toggle",
		TodoID:      todoID,
		Description: description,
		Delay:       delay,
	})
}

// AddDeleteTodoAction adds an action to delete a todo
func (tc *TodoTestCase) AddDeleteTodoAction(todoID, description string, delay time.Duration) {
	tc.Actions = append(tc.Actions, TodoTestAction{
		Type:        "delete",
		TodoID:      todoID,
		Description: description,
		Delay:       delay,
	})
}

// AddClearCompletedAction adds an action to clear all completed todos
func (tc *TodoTestCase) AddClearCompletedAction(description string, delay time.Duration) {
	tc.Actions = append(tc.Actions, TodoTestAction{
		Type:        "clearCompleted",
		Description: description,
		Delay:       delay,
	})
}

// AddSetTextAction adds an action to update a todo's text
func (tc *TodoTestCase) AddSetTextAction(todoID, text, description string, delay time.Duration) {
	tc.Actions = append(tc.Actions, TodoTestAction{
		Type:        "setText",
		TodoID:      todoID,
		TodoText:    text,
		Description: description,
		Delay:       delay,
	})
}

// AddExpectCount adds an expectation for the number of todos
func (tc *TodoTestCase) AddExpectCount(expectedCount int, description string) {
	tc.Results = append(tc.Results, TodoTestResult{
		Type:        "count",
		ExpectedVal: expectedCount,
		Description: description,
	})
}

// AddExpectTodoExists adds an expectation for a todo to exist
func (tc *TodoTestCase) AddExpectTodoExists(todoID string, shouldExist bool, description string) {
	tc.Results = append(tc.Results, TodoTestResult{
		Type:        "exists",
		TodoID:      todoID,
		ExpectedVal: shouldExist,
		Description: description,
	})
}

// AddExpectTodoCompleted adds an expectation for a todo's completion status
func (tc *TodoTestCase) AddExpectTodoCompleted(todoID string, isCompleted bool, description string) {
	tc.Results = append(tc.Results, TodoTestResult{
		Type:        "completed",
		TodoID:      todoID,
		ExpectedVal: isCompleted,
		Description: description,
	})
}

// AddExpectTodoText adds an expectation for a todo's text
func (tc *TodoTestCase) AddExpectTodoText(todoID, expectedText, description string) {
	tc.Results = append(tc.Results, TodoTestResult{
		Type:        "text",
		TodoID:      todoID,
		ExpectedVal: expectedText,
		Description: description,
	})
}

// AddTestCase adds a test case to the test runner
func (tr *TodoTestRunner) AddTestCase(testCase *TodoTestCase) {
	tr.testCases = append(tr.testCases, testCase)
	tr.Log(fmt.Sprintf("Added test case: %s", testCase.Name))
}

// Log logs a message to the status label
func (tr *TodoTestRunner) Log(message string) {
	tr.statusLabel.SetText(message)
	fmt.Println("[TodoTestRunner]", message)
}

// StartTest starts running a test case
func (tr *TodoTestRunner) StartTest(testIndex int) {
	if testIndex < 0 || testIndex >= len(tr.testCases) {
		tr.Log("Invalid test index")
		return
	}

	tr.currentTest = testIndex
	tr.currentStep = -1
	tr.isRunning = true
	tr.lastStepTime = time.Now()

	testCase := tr.testCases[tr.currentTest]
	tr.Log(fmt.Sprintf("Starting test: %s - %s", testCase.Name, testCase.Description))

	// If not in step mode, continue with execution
	if !tr.stepMode {
		tr.ExecuteNextStep()
	}
}

// StopTest stops the current test
func (tr *TodoTestRunner) StopTest() {
	tr.isRunning = false
	tr.Log("Test stopped")
}

// Update updates the test runner
func (tr *TodoTestRunner) Update() {
	// If not running a test, return
	if !tr.isRunning {
		return
	}

	// If in step mode, wait for user to press the Next button
	if tr.stepMode {
		return
	}

	// Check if it's time to execute the next step
	currentTest := tr.testCases[tr.currentTest]
	if tr.currentStep < len(currentTest.Actions) {
		// Get the current action
		action := currentTest.Actions[tr.currentStep]
		
		// Check if the delay has passed
		if time.Since(tr.lastStepTime) >= action.Delay {
			tr.ExecuteNextStep()
		}
	}
}

// ExecuteNextStep executes the next step of the current test case
func (tr *TodoTestRunner) ExecuteNextStep() {
	// Check if we have test cases
	if len(tr.testCases) == 0 || tr.currentTest < 0 || tr.currentTest >= len(tr.testCases) {
		tr.Log("No test case to execute")
		tr.isRunning = false
		return
	}

	// Get the current test case
	currentTest := tr.testCases[tr.currentTest]

	// Increment step counter
	tr.currentStep++

	// Check if we've reached the end of the actions
	if tr.currentStep >= len(currentTest.Actions) {
		// Execute the expectations/assertions
		tr.VerifyTestResults()
		
		// If we're done, stop the test
		tr.Log(fmt.Sprintf("Test '%s' completed", currentTest.Name))
		tr.isRunning = false
		return
	}

	// Get the current action
	action := currentTest.Actions[tr.currentStep]
	
	// Log the current action
	tr.Log(fmt.Sprintf("Step %d: %s", tr.currentStep+1, action.Description))

	// Execute the action based on its type
	switch action.Type {
	case "add":
		tr.executeAddTodoAction(action)
	case "toggle":
		tr.executeToggleTodoAction(action)
	case "delete":
		tr.executeDeleteTodoAction(action)
	case "clearCompleted":
		tr.executeClearCompletedAction(action)
	case "setText":
		tr.executeSetTextAction(action)
	default:
		tr.Log(fmt.Sprintf("Unknown action type: %s", action.Type))
	}

	// Update the timestamp for the next step
	tr.lastStepTime = time.Now()
}

// executeAddTodoAction executes an add todo action
func (tr *TodoTestRunner) executeAddTodoAction(action TodoTestAction) {
	// Set the text in the input field
	tr.inputField.SetText(action.TodoText)
	
	// Click the add button
	tr.addButton.SetOnClick(func() {
		// We just call the internal todo list's AddTodo method directly
		tr.todoList.AddTodo(action.TodoText)
	})
	
	// Trigger the click event
	bounds := tr.addButton.ComputedBounds()
	tr.addButton.HandleMouseDown(bounds.X + bounds.Width/2, bounds.Y + bounds.Height/2)
	tr.addButton.HandleMouseUp(bounds.X + bounds.Width/2, bounds.Y + bounds.Height/2)
}

// executeToggleTodoAction executes a toggle todo action
func (tr *TodoTestRunner) executeToggleTodoAction(action TodoTestAction) {
	// Find the todo with the given ID
	todo := tr.getTodoByID(action.TodoID)
	if todo == nil {
		tr.Log(fmt.Sprintf("Todo with ID %s not found", action.TodoID))
		return
	}

	// Find the checkbox and trigger it
	checkbox := todo.(*components.Todo).GetCheckbox()
	bounds := checkbox.ComputedBounds()
	checkbox.HandleMouseDown(bounds.X + bounds.Width/2, bounds.Y + bounds.Height/2)
	checkbox.HandleMouseUp(bounds.X + bounds.Width/2, bounds.Y + bounds.Height/2)
}

// executeDeleteTodoAction executes a delete todo action
func (tr *TodoTestRunner) executeDeleteTodoAction(action TodoTestAction) {
	// Find the todo with the given ID
	todo := tr.getTodoByID(action.TodoID)
	if todo == nil {
		tr.Log(fmt.Sprintf("Todo with ID %s not found", action.TodoID))
		return
	}

	// Find the delete button and trigger it
	deleteButton := todo.(*components.Todo).GetDeleteButton()
	bounds := deleteButton.ComputedBounds()
	deleteButton.HandleMouseDown(bounds.X + bounds.Width/2, bounds.Y + bounds.Height/2)
	deleteButton.HandleMouseUp(bounds.X + bounds.Width/2, bounds.Y + bounds.Height/2)
}

// executeClearCompletedAction executes a clear completed action
func (tr *TodoTestRunner) executeClearCompletedAction(action TodoTestAction) {
	tr.todoList.ClearCompleted()
}

// executeSetTextAction executes a set text action
func (tr *TodoTestRunner) executeSetTextAction(action TodoTestAction) {
	// Find the todo with the given ID
	todo := tr.getTodoByID(action.TodoID)
	if todo == nil {
		tr.Log(fmt.Sprintf("Todo with ID %s not found", action.TodoID))
		return
	}

	// Set the text
	todo.(*components.Todo).SetText(action.TodoText)
}

// getTodoByID gets a todo by its ID
func (tr *TodoTestRunner) getTodoByID(id string) components.Element {
	// Iterate through the children of the todo list to find a todo with the matching ID
	for _, child := range tr.todoList.Children() {
		if child.ID() == id {
			return child
		}
	}
	return nil
}

// VerifyTestResults verifies the results of the current test case
func (tr *TodoTestRunner) VerifyTestResults() {
	// Get the current test case
	currentTest := tr.testCases[tr.currentTest]

	// Clear previous results
	tr.resultsLabel.SetText("")

	// Track test results
	passed := 0
	failed := 0

	// Check each expectation
	for i, result := range currentTest.Results {
		var success bool
		var message string

		switch result.Type {
		case "count":
			// Check the number of todos
			count := len(tr.todoList.Children())
			expectedCount := result.ExpectedVal.(int)
			success = count == expectedCount
			message = fmt.Sprintf("%d: Count - Expected %d, Got %d - %s", 
				i+1, expectedCount, count, ifThenElse(success, "PASSED", "FAILED"))

		case "exists":
			// Check if a todo exists
			todo := tr.getTodoByID(result.TodoID)
			exists := todo != nil
			expectedExists := result.ExpectedVal.(bool)
			success = exists == expectedExists
			message = fmt.Sprintf("%d: Exists - Todo %s should %s - %s", 
				i+1, result.TodoID, ifThenElse(expectedExists, "exist", "not exist"), ifThenElse(success, "PASSED", "FAILED"))

		case "completed":
			// Check if a todo is completed
			todo := tr.getTodoByID(result.TodoID)
			if todo == nil {
				success = false
				message = fmt.Sprintf("%d: Completed - Todo %s not found - FAILED", i+1, result.TodoID)
			} else {
				expectedCompleted := result.ExpectedVal.(bool)
				actualCompleted := todo.(*components.Todo).GetItem().Done
				success = actualCompleted == expectedCompleted
				message = fmt.Sprintf("%d: Completed - Todo %s should be %s - %s", 
					i+1, result.TodoID, ifThenElse(expectedCompleted, "completed", "not completed"), ifThenElse(success, "PASSED", "FAILED"))
			}

		case "text":
			// Check a todo's text
			todo := tr.getTodoByID(result.TodoID)
			if todo == nil {
				success = false
				message = fmt.Sprintf("%d: Text - Todo %s not found - FAILED", i+1, result.TodoID)
			} else {
				expectedText := result.ExpectedVal.(string)
				actualText := todo.(*components.Todo).GetItem().Text
				success = actualText == expectedText
				message = fmt.Sprintf("%d: Text - Todo %s should have text \"%s\" - %s", 
					i+1, result.TodoID, expectedText, ifThenElse(success, "PASSED", "FAILED"))
			}

		default:
			success = false
			message = fmt.Sprintf("%d: Unknown result type: %s - FAILED", i+1, result.Type)
		}

		// Update counts
		if success {
			passed++
		} else {
			failed++
		}

		// Log the result
		fmt.Println("[TodoTestResult]", message)
	}

	// Update the results label
	tr.resultsLabel.SetText(fmt.Sprintf("Results: %d passed, %d failed", passed, failed))
	if failed > 0 {
		tr.resultsLabel.SetTextColor(color.RGBA{200, 0, 0, 255})
	} else {
		tr.resultsLabel.SetTextColor(color.RGBA{0, 150, 0, 255})
	}
}

// Helper function for conditional expressions
func ifThenElse(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

// RunInteractiveTodoTest creates and returns a todo test runner with a basic test case
func RunInteractiveTodoTest(todoList *components.TodoList, inputField *components.TextArea, addButton *components.Button) *TodoTestRunner {
	// Create test runner
	runner := NewTodoTestRunner(todoList, inputField, addButton)

	// Create a simple test case
	testCase := NewTodoTestCase("Basic Todo Operations", "Tests adding, completing, and deleting todos")
	
	// Add test actions
	testCase.AddAddTodoAction("Buy groceries", "Add a todo to buy groceries", 500*time.Millisecond)
	testCase.AddAddTodoAction("Walk the dog", "Add a todo to walk the dog", 500*time.Millisecond)
	testCase.AddAddTodoAction("Pay bills", "Add a todo to pay bills", 500*time.Millisecond)
	testCase.AddToggleTodoAction("todo_1", "Complete the first todo", 500*time.Millisecond)
	testCase.AddDeleteTodoAction("todo_3", "Delete the third todo", 500*time.Millisecond)
	testCase.AddClearCompletedAction("Clear completed todos", 500*time.Millisecond)
	
	// Add expectations
	testCase.AddExpectCount(1, "Should have 1 todo remaining")
	testCase.AddExpectTodoExists("todo_1", false, "First todo should be deleted (was completed)")
	testCase.AddExpectTodoExists("todo_2", true, "Second todo should still exist")
	testCase.AddExpectTodoExists("todo_3", false, "Third todo should be deleted")
	testCase.AddExpectTodoText("todo_2", "Walk the dog", "Todo text should match")
	
	// Add the test case to the runner
	runner.AddTestCase(testCase)
	
	return runner
} 