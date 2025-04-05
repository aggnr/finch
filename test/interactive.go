package test

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"time"

	"github.com/aggnr/finch/components"
)

// InteractiveAction represents a single action in a test scenario
type InteractiveAction struct {
	Type      string       // "click", "hover", "input", "wait"
	Target    components.Element // Target element for the action
	X, Y      int          // Coordinates for actions like click
	Value     interface{}  // Additional value for the action (e.g., text input)
	Delay     time.Duration // Time to wait before next action
	Message   string       // Description of the action
}

// ActionResult stores the result of a test action
type ActionResult struct {
	Action    InteractiveAction
	Success   bool
	Message   string
	Timestamp time.Time
	Screenshot string   // Path to screenshot file
}

// TestCase represents a sequence of actions to test UI components
type TestCase struct {
	Name        string
	Description string
	Actions     []InteractiveAction
	Results     []ActionResult
	CurrentStep int
	Completed   bool
	Success     bool
}

// InteractiveTest manages the execution of UI tests with interactive control
type InteractiveTest struct {
	rootElement    components.Element
	testCases      []*TestCase
	currentTest    int
	running        bool
	paused         bool
	stepByStep     bool
	speed          float64 // Speed multiplier (1.0 = normal)
	surface        components.DrawSurface
	simulatedMouse image.Point
	showMouse      bool
	logFile        *os.File
	screenshotDir  string
}

// NewInteractiveTest creates a new interactive test manager
func NewInteractiveTest(root components.Element) *InteractiveTest {
	// Create logs and screenshots directories if they don't exist
	os.MkdirAll("logs", 0755)
	os.MkdirAll("screenshots", 0755)
	
	// Create log file with timestamp
	logFile, err := os.Create(fmt.Sprintf("logs/ui_test_%s.log", 
		time.Now().Format("20060102_150405")))
	if err != nil {
		fmt.Println("Error creating log file:", err)
	}
	
	return &InteractiveTest{
		rootElement:    root,
		testCases:      make([]*TestCase, 0),
		surface:        NewMemorySurface(components.ScreenWidth, components.ScreenHeight),
		simulatedMouse: image.Point{-100, -100}, // Off-screen initially
		showMouse:      true,
		logFile:        logFile,
		speed:          1.0,
		stepByStep:     false,
		screenshotDir:  "screenshots",
	}
}

// Close closes resources held by the test
func (t *InteractiveTest) Close() {
	if t.logFile != nil {
		t.logFile.Close()
	}
}

// Log writes a message to the log file and console
func (t *InteractiveTest) Log(message string) {
	if t.logFile != nil {
		fmt.Fprintf(t.logFile, "[%s] %s\n", time.Now().Format("15:04:05.000"), message)
	}
	fmt.Println(message)
}

// AddTestCase adds a test case to the test suite
func (t *InteractiveTest) AddTestCase(testCase *TestCase) {
	t.testCases = append(t.testCases, testCase)
}

// NewTestCase creates a new test case
func NewTestCase(name, description string) *TestCase {
	return &TestCase{
		Name:        name,
		Description: description,
		Actions:     make([]InteractiveAction, 0),
		Results:     make([]ActionResult, 0),
	}
}

// AddClickAction adds a click action to a test case
func (tc *TestCase) AddClickAction(target components.Element, x, y int, message string) {
	tc.Actions = append(tc.Actions, InteractiveAction{
		Type:    "click",
		Target:  target,
		X:       x,
		Y:       y,
		Delay:   500 * time.Millisecond,
		Message: message,
	})
}

// AddHoverAction adds a hover action to a test case
func (tc *TestCase) AddHoverAction(target components.Element, x, y int, message string) {
	tc.Actions = append(tc.Actions, InteractiveAction{
		Type:    "hover",
		Target:  target,
		X:       x,
		Y:       y,
		Delay:   300 * time.Millisecond,
		Message: message,
	})
}

// AddWaitAction adds a wait action to a test case
func (tc *TestCase) AddWaitAction(duration time.Duration, message string) {
	tc.Actions = append(tc.Actions, InteractiveAction{
		Type:    "wait",
		Delay:   duration,
		Message: message,
	})
}

// Run runs the interactive test
func (t *InteractiveTest) Run() {
	if len(t.testCases) == 0 {
		t.Log("No test cases defined")
		return
	}
	
	t.Log("Starting interactive UI test...")
	t.Log(fmt.Sprintf("Found %d test cases", len(t.testCases)))
	
	reader := bufio.NewReader(os.Stdin)
	
	// Show available commands
	t.printHelp()
	
	t.currentTest = 0
	keepRunning := true
	
	for keepRunning {
		fmt.Print("\nCommand (type 'help' for commands): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		switch input {
		case "help", "h":
			t.printHelp()
			
		case "list", "l":
			t.listTestCases()
			
		case "start", "run", "r":
			t.startTests()
			
		case "stop", "s":
			t.stopTests()
			
		case "pause", "p":
			t.pauseTests()
			
		case "stepby", "step":
			t.toggleStepByStep()
			
		case "next", "n":
			t.executeNextStep()
			
		case "screenshot", "ss":
			t.takeScreenshot("manual")
			
		case "speed+":
			t.speed *= 1.5
			t.Log(fmt.Sprintf("Speed increased to %.1fx", t.speed))
			
		case "speed-":
			t.speed /= 1.5
			if t.speed < 0.1 {
				t.speed = 0.1
			}
			t.Log(fmt.Sprintf("Speed decreased to %.1fx", t.speed))
			
		case "reset":
			t.resetTests()
			
		case "exit", "quit", "q":
			keepRunning = false
			
		default:
			if strings.HasPrefix(input, "case ") {
				// Select a specific test case
				var caseNum int
				fmt.Sscanf(input, "case %d", &caseNum)
				t.selectTestCase(caseNum - 1) // Convert to 0-based index
			} else {
				t.Log("Unknown command. Type 'help' for available commands.")
			}
		}
	}
	
	t.Log("Interactive test completed.")
}

// printHelp shows available commands
func (t *InteractiveTest) printHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  help       - Show this help message")
	fmt.Println("  list       - List all test cases")
	fmt.Println("  case N     - Select test case number N")
	fmt.Println("  start/run  - Start/resume test execution")
	fmt.Println("  stop       - Stop test execution")
	fmt.Println("  pause      - Pause test execution")
	fmt.Println("  stepby     - Toggle step-by-step mode")
	fmt.Println("  next       - Execute next step (in step-by-step mode)")
	fmt.Println("  screenshot - Take a screenshot")
	fmt.Println("  speed+     - Increase test speed")
	fmt.Println("  speed-     - Decrease test speed")
	fmt.Println("  reset      - Reset all tests to initial state")
	fmt.Println("  exit/quit  - Exit the test")
}

// listTestCases lists all test cases
func (t *InteractiveTest) listTestCases() {
	fmt.Println("\nTest Cases:")
	for i, tc := range t.testCases {
		status := "Not Run"
		if tc.Completed {
			if tc.Success {
				status = "Passed"
			} else {
				status = "Failed"
			}
		} else if i == t.currentTest && t.running {
			status = "Running"
		}
		
		fmt.Printf("  %d. %s - %s [%s]\n", i+1, tc.Name, tc.Description, status)
	}
}

// selectTestCase selects a specific test case
func (t *InteractiveTest) selectTestCase(index int) {
	if index >= 0 && index < len(t.testCases) {
		t.currentTest = index
		t.Log(fmt.Sprintf("Selected test case: %s", t.testCases[index].Name))
	} else {
		t.Log(fmt.Sprintf("Invalid test case index. Valid range is 1-%d", len(t.testCases)))
	}
}

// startTests starts or resumes test execution
func (t *InteractiveTest) startTests() {
	if len(t.testCases) == 0 {
		t.Log("No test cases defined")
		return
	}
	
	t.running = true
	t.paused = false
	
	t.Log(fmt.Sprintf("Starting test case: %s", t.testCases[t.currentTest].Name))
	
	// If not in step-by-step mode, execute the test
	if !t.stepByStep {
		t.executeTestCase(t.testCases[t.currentTest])
	} else {
		t.Log("Step-by-step mode active. Use 'next' to execute each step.")
	}
}

// stopTests stops test execution
func (t *InteractiveTest) stopTests() {
	t.running = false
	t.paused = false
	t.Log("Test execution stopped")
}

// pauseTests pauses test execution
func (t *InteractiveTest) pauseTests() {
	if t.running {
		t.paused = !t.paused
		if t.paused {
			t.Log("Test execution paused")
		} else {
			t.Log("Test execution resumed")
		}
	} else {
		t.Log("No test is currently running")
	}
}

// toggleStepByStep toggles step-by-step execution mode
func (t *InteractiveTest) toggleStepByStep() {
	t.stepByStep = !t.stepByStep
	if t.stepByStep {
		t.Log("Step-by-step mode enabled")
	} else {
		t.Log("Step-by-step mode disabled")
	}
}

// executeNextStep executes the next step in the current test case
func (t *InteractiveTest) executeNextStep() {
	if !t.running {
		t.Log("No test is currently running. Use 'start' to begin test execution.")
		return
	}
	
	if t.paused {
		t.Log("Test is paused. Use 'pause' to resume.")
		return
	}
	
	currentTest := t.testCases[t.currentTest]
	
	if currentTest.CurrentStep >= len(currentTest.Actions) {
		t.Log("Current test case has no more steps.")
		currentTest.Completed = true
		
		// Move to next test case if available
		if t.currentTest < len(t.testCases) - 1 {
			t.currentTest++
			t.Log(fmt.Sprintf("Moving to next test case: %s", t.testCases[t.currentTest].Name))
		} else {
			t.Log("All test cases completed")
			t.running = false
		}
		return
	}
	
	// Execute the current action
	action := currentTest.Actions[currentTest.CurrentStep]
	t.Log(fmt.Sprintf("Executing step %d/%d: %s", 
		currentTest.CurrentStep+1, 
		len(currentTest.Actions),
		action.Message))
	
	result := t.executeAction(action)
	
	// Store result
	currentTest.Results = append(currentTest.Results, result)
	
	// Move to next step
	currentTest.CurrentStep++
	
	// Take screenshot after action
	
	// If not in step-by-step mode and there are more steps, continue execution
	if !t.stepByStep && currentTest.CurrentStep < len(currentTest.Actions) {
		// Add delay based on action and speed
		delay := time.Duration(float64(action.Delay) / t.speed)
		t.Log(fmt.Sprintf("Waiting for %.2f seconds before next step...", delay.Seconds()))
		time.Sleep(delay)
		
		// Execute next step (when not in step-by-step mode)
		t.executeNextStep()
	}
}

// executeTestCase executes all steps in a test case
func (t *InteractiveTest) executeTestCase(testCase *TestCase) {
	if t.stepByStep {
		t.Log("In step-by-step mode. Use 'next' to execute each step.")
		return
	}
	
	t.Log(fmt.Sprintf("Executing test case: %s", testCase.Name))
	
	// Reset test case if already run
	if testCase.CurrentStep > 0 {
		testCase.CurrentStep = 0
		testCase.Results = make([]ActionResult, 0)
		testCase.Completed = false
		testCase.Success = false
	}
	
	// Initial screenshot
	t.takeScreenshot(fmt.Sprintf("test%d_initial", t.currentTest+1))
	
	// Execute first step - the rest will chain via executeNextStep
	t.executeNextStep()
}

// executeAction performs a single test action
func (t *InteractiveTest) executeAction(action InteractiveAction) ActionResult {
	result := ActionResult{
		Action:    action,
		Success:   true,
		Message:   "Action completed",
		Timestamp: time.Now(),
	}
	
	switch action.Type {
	case "click":
		// Update simulated mouse position
		t.simulatedMouse = image.Point{action.X, action.Y}
		
		// Render to show cursor position
		t.rootElement.Draw(t.surface)
		t.drawSimulatedMouse()
		t.takeScreenshot(fmt.Sprintf("test%d_click_before", t.currentTest+1))
		
		// Handle mouse down event
		downResult := t.rootElement.HandleMouseDown(action.X, action.Y)
		
		// Small delay to simulate real interaction
		time.Sleep(100 * time.Millisecond)
		
		// Render to show pressed state
		t.rootElement.Draw(t.surface)
		t.drawSimulatedMouse()
		t.takeScreenshot(fmt.Sprintf("test%d_click_down", t.currentTest+1))
		
		// Handle mouse up event
		upResult := t.rootElement.HandleMouseUp(action.X, action.Y)
		
		if !downResult && !upResult {
			result.Success = false
			result.Message = fmt.Sprintf("Click at (%d, %d) was not handled by any element", 
				action.X, action.Y)
		}
		
	case "hover":
		// Update simulated mouse position
		t.simulatedMouse = image.Point{action.X, action.Y}
		
		// Handle mouse move event
		moveResult := t.rootElement.HandleMouseMove(action.X, action.Y)
		
		if !moveResult {
			result.Success = false
			result.Message = fmt.Sprintf("Hover at (%d, %d) was not handled by any element", 
				action.X, action.Y)
		}
		
	case "wait":
		// Just wait for the specified delay
		// Nothing to do here, delay is handled at caller level
	}
	
	return result
}

// resetTests resets all tests to initial state
func (t *InteractiveTest) resetTests() {
	for _, testCase := range t.testCases {
		testCase.CurrentStep = 0
		testCase.Results = make([]ActionResult, 0)
		testCase.Completed = false
		testCase.Success = false
	}
	
	t.Log("All tests reset to initial state")
}

// takeScreenshot captures the current UI state and saves it to a file
func (t *InteractiveTest) takeScreenshot(prefix string) string {
	filename := fmt.Sprintf("%s/%s_%s.png", 
		t.screenshotDir, 
		prefix, 
		time.Now().Format("150405"))
	
	// Render UI
	t.rootElement.Draw(t.surface)
	
	// Draw simulated mouse if visible
	if t.showMouse {
		t.drawSimulatedMouse()
	}
	
	// Save image
	image := t.surface.(*MemorySurface).Image()
	
	f, err := os.Create(filename)
	if err != nil {
		t.Log(fmt.Sprintf("Error creating screenshot file: %v", err))
		return ""
	}
	defer f.Close()
	
	if err := png.Encode(f, image); err != nil {
		t.Log(fmt.Sprintf("Error encoding screenshot: %v", err))
		return ""
	}
	
	t.Log(fmt.Sprintf("Screenshot saved to %s", filename))
	return filename
}

// drawSimulatedMouse draws a cursor for automated testing
func (t *InteractiveTest) drawSimulatedMouse() {
	surface := t.surface.(*MemorySurface)
	
	// Draw a red crosshair
	red := color.RGBA{255, 0, 0, 255}
	
	// Draw circle
	surface.FillCircle(t.simulatedMouse.X, t.simulatedMouse.Y, 5, red)
	
	// Draw crosshair
	surface.DrawLine(
		t.simulatedMouse.X-8, t.simulatedMouse.Y,
		t.simulatedMouse.X+8, t.simulatedMouse.Y, red)
	surface.DrawLine(
		t.simulatedMouse.X, t.simulatedMouse.Y-8,
		t.simulatedMouse.X, t.simulatedMouse.Y+8, red)
}

// MemorySurface from previous implementation 