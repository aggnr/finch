package test

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/aggnr/finch/components"
)

// Global reference to the current UITestGame instance
var currentTestGame *UITestGame

// UITestFrame is the main container for the UI test
type UITestFrame struct {
	*components.BaseElement
	testCases     []*UITestCase
	currentTest   int
	playingTest   bool
	stepMode      bool
	currentStep   int
	controls      *TestControls
	logPanel      *LogPanel
	statusLabel   *components.Label
	testResult    *components.Label
	rootElement   components.Element // Root element of the UI being tested
}

// UITestCase represents a sequence of test actions
type UITestCase struct {
	Name        string
	Description string
	Actions     []UITestAction
	Results     []string
}

// UITestAction represents a single action to be performed during testing
type UITestAction struct {
	Type         string          // "click", "hover", "wait"
	TargetID     string          // ID of target element
	Target       components.Element // Reference to target
	X, Y         int             // Coordinates for actions like click
	Description  string          // Human-readable description
	Delay        time.Duration   // Delay after action
}

// TestControls contains the control buttons and UI for the test framework
type TestControls struct {
	*components.BaseElement
	playButton    *components.Button
	pauseButton   *components.Button
	stepButton    *components.Button
	stopButton    *components.Button
	resetButton   *components.Button
	nextTestButton *components.Button
	prevTestButton *components.Button
	stepModeCheckbox *components.Button // Using button as a toggle
	speedButton *components.Button // Button to cycle through speeds
	stepModeActive bool
	speedIndex int // 0=slow, 1=normal, 2=fast
	testCaseLabel *components.Label
}

// LogPanel displays test event logs
type LogPanel struct {
	*components.BaseElement
	logs        []string
	maxLogs     int
	scrollY     int
	logLabels   []*components.Label
}

// NewUITestFrame creates a new UI test frame
func NewUITestFrame(rootUI components.Element) *UITestFrame {
	// Create the test frame container
	frame := &UITestFrame{
		BaseElement: components.NewBaseElement("ui_test_frame"),
		testCases:   make([]*UITestCase, 0),
		currentTest: 0,
		playingTest: false,
		stepMode:    false,
		currentStep: -1,
		rootElement: rootUI,
	}
	
	// Calculate the layout - test frame appears at the bottom of the window
	frameHeight := 200
	frameWidth := components.ScreenWidth
	frameX := 0
	frameY := components.ScreenHeight - frameHeight
	frame.SetBounds(components.Rect{X: frameX, Y: frameY, Width: frameWidth, Height: frameHeight})
	
	// Create and add test controls
	frame.controls = createTestControls(frame)
	frame.controls.SetBounds(components.Rect{X: 10, Y: 10, Width: frameWidth - 20, Height: 40})
	frame.AddChild(frame.controls)
	
	// Create status label
	frame.statusLabel = components.NewLabel("status_label", "Test Framework Ready", 14, color.RGBA{0, 0, 0, 255})
	frame.statusLabel.SetBounds(components.Rect{X: 10, Y: 60, Width: 400, Height: 20})
	frame.AddChild(frame.statusLabel)
	
	// Create test result label
	frame.testResult = components.NewLabel("test_result", "", 14, color.RGBA{0, 100, 0, 255})
	frame.testResult.SetBounds(components.Rect{X: 10, Y: 80, Width: 400, Height: 20})
	frame.AddChild(frame.testResult)
	
	// Create log panel at bottom of frame - position it to use the remaining space
	logPanelHeight := frameHeight - 110 // Allow space for controls and status
	frame.logPanel = createLogPanel()
	frame.logPanel.SetBounds(components.Rect{X: 10, Y: 105, Width: frameWidth - 20, Height: logPanelHeight})
	frame.AddChild(frame.logPanel)
	
	// Set up button handlers
	frame.setupControlHandlers()
	
	return frame
}

// createTestControls creates the test control buttons
func createTestControls(frame *UITestFrame) *TestControls {
	controls := &TestControls{
		BaseElement: components.NewBaseElement("test_controls"),
		stepModeActive: false,
		speedIndex: 1, // Default to normal speed
	}
	
	// Calculate absolute coordinates
	frameY := components.ScreenHeight - 200 // 200 is the frame height
	
	// Create play button
	playButton := components.NewButton("play_button", "▶ Play")
	playButton.SetBounds(components.Rect{X: 10, Y: frameY + 10, Width: 80, Height: 30})
	playButton.SetBackgroundColor(color.RGBA{60, 179, 113, 255})
	playButton.SetTextColor(color.RGBA{255, 255, 255, 255})
	controls.playButton = playButton
	controls.AddChild(playButton)
	
	// Create pause button
	pauseButton := components.NewButton("pause_button", "⏸ Pause")
	pauseButton.SetBounds(components.Rect{X: 100, Y: frameY + 10, Width: 80, Height: 30})
	pauseButton.SetBackgroundColor(color.RGBA{255, 165, 0, 255})
	pauseButton.SetTextColor(color.RGBA{255, 255, 255, 255})
	controls.pauseButton = pauseButton
	controls.AddChild(pauseButton)
	
	// Create step button
	stepButton := components.NewButton("step_button", "⏭ Step")
	stepButton.SetBounds(components.Rect{X: 190, Y: frameY + 10, Width: 80, Height: 30})
	stepButton.SetBackgroundColor(color.RGBA{30, 144, 255, 255})
	stepButton.SetTextColor(color.RGBA{255, 255, 255, 255})
	controls.stepButton = stepButton
	controls.AddChild(stepButton)
	
	// Create stop button
	stopButton := components.NewButton("stop_button", "⏹ Stop")
	stopButton.SetBounds(components.Rect{X: 280, Y: frameY + 10, Width: 80, Height: 30})
	stopButton.SetBackgroundColor(color.RGBA{220, 20, 60, 255})
	stopButton.SetTextColor(color.RGBA{255, 255, 255, 255})
	controls.stopButton = stopButton
	controls.AddChild(stopButton)
	
	// Create test navigation buttons
	prevTestButton := components.NewButton("prev_test_button", "◀ Prev")
	prevTestButton.SetBounds(components.Rect{X: 370, Y: frameY + 10, Width: 80, Height: 30})
	prevTestButton.SetBackgroundColor(color.RGBA{100, 100, 100, 255})
	prevTestButton.SetTextColor(color.RGBA{255, 255, 255, 255})
	controls.prevTestButton = prevTestButton
	controls.AddChild(prevTestButton)
	
	nextTestButton := components.NewButton("next_test_button", "Next ▶")
	nextTestButton.SetBounds(components.Rect{X: 460, Y: frameY + 10, Width: 80, Height: 30})
	nextTestButton.SetBackgroundColor(color.RGBA{100, 100, 100, 255})
	nextTestButton.SetTextColor(color.RGBA{255, 255, 255, 255})
	controls.nextTestButton = nextTestButton
	controls.AddChild(nextTestButton)
	
	// Create reset button
	resetButton := components.NewButton("reset_button", "↺ Reset")
	resetButton.SetBounds(components.Rect{X: 550, Y: frameY + 10, Width: 80, Height: 30})
	resetButton.SetBackgroundColor(color.RGBA{128, 0, 128, 255})
	resetButton.SetTextColor(color.RGBA{255, 255, 255, 255})
	controls.resetButton = resetButton
	controls.AddChild(resetButton)
	
	// Create step mode toggle
	stepModeButton := components.NewButton("step_mode_button", "□ Step Mode")
	stepModeButton.SetBounds(components.Rect{X: 640, Y: frameY + 10, Width: 100, Height: 30})
	stepModeButton.SetBackgroundColor(color.RGBA{70, 70, 70, 255})
	stepModeButton.SetTextColor(color.RGBA{255, 255, 255, 255})
	controls.stepModeCheckbox = stepModeButton
	controls.AddChild(stepModeButton)
	
	// Create speed button
	speedButton := components.NewButton("speed_button", "🕒 Normal")
	speedButton.SetBounds(components.Rect{X: 750, Y: frameY + 10, Width: 100, Height: 30})
	speedButton.SetBackgroundColor(color.RGBA{70, 130, 180, 255})
	speedButton.SetTextColor(color.RGBA{255, 255, 255, 255})
	controls.speedButton = speedButton
	controls.AddChild(speedButton)
	
	// Create test case label
	testCaseLabel := components.NewLabel("test_case_label", "No Test Selected", 14, color.RGBA{0, 0, 0, 255})
	testCaseLabel.SetBounds(components.Rect{X: 860, Y: frameY + 15, Width: 200, Height: 20})
	controls.testCaseLabel = testCaseLabel
	controls.AddChild(testCaseLabel)
	
	return controls
}

// createLogPanel creates the log panel
func createLogPanel() *LogPanel {
	logPanel := &LogPanel{
		BaseElement: components.NewBaseElement("log_panel"),
		logs:        make([]string, 0),
		maxLogs:     6, // Reduced to avoid overlapping
		scrollY:     0,
		logLabels:   make([]*components.Label, 0),
	}
	
	// Create initial log labels with more spacing between them
	const lineHeight = 20 // Increased from 18 to give more space
	for i := 0; i < logPanel.maxLogs; i++ {
		label := components.NewLabel(fmt.Sprintf("log_%d", i), "", 12, color.RGBA{50, 50, 50, 255})
		label.SetBounds(components.Rect{X: 5, Y: i*lineHeight, Width: components.ScreenWidth - 30, Height: 16})
		logPanel.logLabels = append(logPanel.logLabels, label)
		logPanel.AddChild(label)
	}
	
	return logPanel
}

// Draw draws the test frame and its children
func (f *UITestFrame) Draw(surface components.DrawSurface) {
	// Get the frame bounds
	bounds := f.Bounds()
	
	// Draw background with a gradient effect
	surface.FillRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{230, 230, 230, 255})
	
	// Draw a border around the frame
	surface.DrawRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, color.RGBA{50, 50, 150, 255})
	
	// Draw a separator line below the controls
	controlsBottom := bounds.Y + 50
	surface.DrawLine(bounds.X, controlsBottom, bounds.X + bounds.Width, controlsBottom, color.RGBA{150, 150, 150, 255})
	
	// Draw title text
	surface.DrawText("FINCH UI TEST FRAMEWORK", bounds.X + 10, bounds.Y + 5, color.RGBA{50, 50, 150, 255}, 12)
	
	// Draw log panel background and border
	logPanelBounds := f.logPanel.Bounds()
	surface.FillRect(logPanelBounds.X, logPanelBounds.Y, logPanelBounds.Width, logPanelBounds.Height, color.RGBA{245, 245, 245, 255})
	surface.DrawRect(logPanelBounds.X, logPanelBounds.Y, logPanelBounds.Width, logPanelBounds.Height, color.RGBA{100, 100, 150, 255})
	
	// Draw log panel title
	surface.DrawText("TEST LOG OUTPUT", logPanelBounds.X + 5, logPanelBounds.Y - 5, color.RGBA{50, 50, 120, 255}, 12)
	
	// Draw all child elements
	for _, child := range f.Children() {
		child.Draw(surface)
	}
	
	// Draw debug info at the bottom of the test frame
	debugY := bounds.Y + bounds.Height - 20
	debugX := bounds.X + 10
	
	// Get mouse position from the game instance
	mouseX, mouseY := 0, 0
	if currentTestGame != nil {
		mouseX, mouseY = currentTestGame.mouseX, currentTestGame.mouseY
	}
	
	// Show mouse position
	mouseInfo := fmt.Sprintf("Mouse: (%d,%d)", mouseX, mouseY)
	surface.DrawText(mouseInfo, debugX, debugY, color.RGBA{50, 50, 50, 255}, 10)
	
	// Show component under mouse if any
	if mouseX > 0 && mouseY > 0 && f.rootElement != nil {
		if element := findElementAtPosition(f.rootElement, mouseX, mouseY); element != nil {
			elementBounds := element.Bounds()
			elementInfo := fmt.Sprintf("Element: %s at (%d,%d,%d,%d)", 
				element.ID(), elementBounds.X, elementBounds.Y, 
				elementBounds.Width, elementBounds.Height)
			surface.DrawText(elementInfo, debugX + 150, debugY, color.RGBA{50, 50, 50, 255}, 10)
		}
	}
	
	// Display hint about inspector mode
	surface.DrawText("Press 'I' to toggle component inspector", debugX + 500, debugY, color.RGBA{50, 50, 50, 255}, 10)
}

// findElementAtPosition recursively finds the element at the given position
func findElementAtPosition(element components.Element, x, y int) components.Element {
	// Check if point is inside this element
	if !components.PointInRect(components.Point{X: x, Y: y}, element.Bounds()) {
		return nil
	}
	
	// Check children first (in reverse order to get topmost element)
	for i := len(element.Children()) - 1; i >= 0; i-- {
		child := element.Children()[i]
		if found := findElementAtPosition(child, x, y); found != nil {
			return found
		}
	}
	
	// If no matching child, return this element
	return element
}

// AddTestCase adds a test case to the test frame
func (f *UITestFrame) AddTestCase(testCase *UITestCase) {
	f.testCases = append(f.testCases, testCase)
	
	// Update label if this is the first test case
	if len(f.testCases) == 1 {
		f.updateTestCaseLabel()
	}
}

// setupControlHandlers sets up the button event handlers
func (f *UITestFrame) setupControlHandlers() {
	// Play button
	f.controls.playButton.SetOnClick(func() {
		fmt.Println("Play button clicked")
		if len(f.testCases) == 0 {
			f.Log("No test cases to run")
			return
		}
		
		f.playingTest = true
		f.statusLabel.SetText("Running test: " + f.testCases[f.currentTest].Name)
		f.Log("Started test: " + f.testCases[f.currentTest].Name)
	})
	
	// Pause button
	f.controls.pauseButton.SetOnClick(func() {
		fmt.Println("Pause button clicked")
		if f.playingTest {
			f.playingTest = false
			f.statusLabel.SetText("Paused test: " + f.testCases[f.currentTest].Name)
			f.Log("Paused test")
		}
	})
	
	// Step button
	f.controls.stepButton.SetOnClick(func() {
		fmt.Println("Step button clicked")
		if len(f.testCases) == 0 {
			f.Log("No test cases to run")
			return
		}
		
		f.ExecuteNextStep()
	})
	
	// Stop button
	f.controls.stopButton.SetOnClick(func() {
		fmt.Println("Stop button clicked")
		if f.playingTest {
			f.playingTest = false
			f.statusLabel.SetText("Stopped test: " + f.testCases[f.currentTest].Name)
			f.Log("Stopped test")
		}
	})
	
	// Reset button
	f.controls.resetButton.SetOnClick(func() {
		fmt.Println("Reset button clicked")
		f.ResetTest()
	})
	
	// Step mode toggle
	f.controls.stepModeCheckbox.SetOnClick(func() {
		fmt.Println("Step mode button clicked")
		f.controls.stepModeActive = !f.controls.stepModeActive
		if f.controls.stepModeActive {
			f.controls.stepModeCheckbox.SetText("☑ Step Mode")
			f.stepMode = true
			f.Log("Step mode enabled")
		} else {
			f.controls.stepModeCheckbox.SetText("□ Step Mode")
			f.stepMode = false
			f.Log("Step mode disabled")
		}
	})
	
	// Speed button
	f.controls.speedButton.SetOnClick(func() {
		fmt.Println("Speed button clicked")
		// Cycle through speeds: Slow -> Normal -> Fast -> Slow
		f.controls.speedIndex = (f.controls.speedIndex + 1) % 3
		
		// Update button text based on speed
		switch f.controls.speedIndex {
		case 0:
			f.controls.speedButton.SetText("🐢 Slow")
			f.controls.speedButton.SetBackgroundColor(color.RGBA{70, 130, 180, 255})
			f.Log("Test speed set to SLOW")
		case 1:
			f.controls.speedButton.SetText("🕒 Normal")
			f.controls.speedButton.SetBackgroundColor(color.RGBA{70, 130, 180, 255})
			f.Log("Test speed set to NORMAL")
		case 2:
			f.controls.speedButton.SetText("🚀 Fast")
			f.controls.speedButton.SetBackgroundColor(color.RGBA{70, 130, 180, 255})
			f.Log("Test speed set to FAST")
		}
	})
	
	// Next test button
	f.controls.nextTestButton.SetOnClick(func() {
		fmt.Println("Next test button clicked")
		if len(f.testCases) == 0 {
			return
		}
		
		f.currentTest = (f.currentTest + 1) % len(f.testCases)
		f.updateTestCaseLabel()
		f.ResetTest()
	})
	
	// Previous test button
	f.controls.prevTestButton.SetOnClick(func() {
		fmt.Println("Previous test button clicked")
		if len(f.testCases) == 0 {
			return
		}
		
		f.currentTest = (f.currentTest - 1 + len(f.testCases)) % len(f.testCases)
		f.updateTestCaseLabel()
		f.ResetTest()
	})
}

// updateTestCaseLabel updates the test case label with the current test name
func (f *UITestFrame) updateTestCaseLabel() {
	if len(f.testCases) > 0 {
		f.controls.testCaseLabel.SetText(fmt.Sprintf("Test %d/%d: %s", 
			f.currentTest+1, 
			len(f.testCases), 
			f.testCases[f.currentTest].Name))
	} else {
		f.controls.testCaseLabel.SetText("No Tests Available")
	}
}

// ResetTest resets the current test to initial state
func (f *UITestFrame) ResetTest() {
	f.currentStep = -1
	f.testResult.SetText("")
	f.playingTest = false
	f.statusLabel.SetText("Test reset: Ready to run")
	f.Log("Test reset")
}

// Update updates the test frame and processes test actions if tests are running
func (f *UITestFrame) Update() {
	// Debug test controls visibility
	if f.controls != nil && f.controls.playButton != nil {
		bounds := f.controls.playButton.Bounds()
		// Only print occasionally to avoid log spam
		if time.Now().Second() % 5 == 0 {
			fmt.Printf("Play button position: X=%d, Y=%d, W=%d, H=%d\n", 
				bounds.X, bounds.Y, bounds.Width, bounds.Height)
		}
	}
	
	// Update controls and other UI elements
	for _, child := range f.Children() {
		child.Update()
	}
	
	// Process test steps if playing and not in step mode
	if f.playingTest && !f.stepMode {
		f.ExecuteNextStep()
	}
}

// ExecuteNextStep executes the next step in the current test case
func (f *UITestFrame) ExecuteNextStep() {
	// Check if we have test cases
	if len(f.testCases) == 0 {
		f.Log("No test cases to run")
		return
	}
	
	// Get current test case
	testCase := f.testCases[f.currentTest]
	
	// Move to next step
	f.currentStep++
	
	// Check if test is complete
	if f.currentStep >= len(testCase.Actions) {
		f.statusLabel.SetText("Test completed: " + testCase.Name)
		f.testResult.SetText("Test Passed!")
		f.testResult.SetTextColor(color.RGBA{0, 128, 0, 255})
		f.Log("Test completed successfully")
		
		// If in step mode, don't auto-advance; wait for next button click
		if f.stepMode {
			f.currentStep = -1
			f.playingTest = false
			return
		}
		
		// Check if this was the last test case
		if f.currentTest == len(f.testCases) - 1 {
			f.Log("All test cases completed")
			f.statusLabel.SetText("All test cases completed")
			f.playingTest = false
			f.currentStep = -1
			return
		}
		
		// Move to next test case
		f.currentTest++
		f.updateTestCaseLabel()
		f.currentStep = -1
		
		// Brief pause before starting next test case
		time.Sleep(1 * time.Second)
		
		// Continue testing with the next test case
		f.statusLabel.SetText("Running test: " + f.testCases[f.currentTest].Name)
		f.Log("Starting next test: " + f.testCases[f.currentTest].Name)
		f.testResult.SetText("")
		return
	}
	
	// Get current action
	action := testCase.Actions[f.currentStep]
	
	// Log the action
	f.Log(fmt.Sprintf("Step %d/%d: %s", 
		f.currentStep+1, 
		len(testCase.Actions), 
		action.Description))
	
	// Update status
	f.statusLabel.SetText(fmt.Sprintf("Running step %d/%d: %s", 
		f.currentStep+1, 
		len(testCase.Actions), 
		action.Description))
	
	// Execute the action
	f.executeAction(action)
}

// executeAction performs a single test action
func (f *UITestFrame) executeAction(action UITestAction) {
	// Get the UITestGame instance to update the virtual cursor
	game := currentTestGame
	
	// Get delay multiplier based on speed setting
	var delayMultiplier float64
	switch f.controls.speedIndex {
	case 0: // Slow
		delayMultiplier = 2.0
	case 1: // Normal
		delayMultiplier = 1.0
	case 2: // Fast
		delayMultiplier = 0.5
	default:
		delayMultiplier = 1.0
	}
	
	switch action.Type {
	case "click":
		// Find the target if needed
		if action.Target == nil && action.TargetID != "" {
			fmt.Printf("Looking for target element with ID: %s\n", action.TargetID)
			action.Target = f.FindElementByID(action.TargetID)
		}
		
		if action.Target != nil {
			// Log the element found
			bounds := action.Target.Bounds()
			fmt.Printf("Found element: %s at (%d,%d,%d,%d)\n", 
				action.Target.ID(), bounds.X, bounds.Y, bounds.Width, bounds.Height)
			f.Log(fmt.Sprintf("Found element: %s at (%d,%d)", action.Target.ID(), bounds.X, bounds.Y))
			
			// Get element bounds to calculate center if x,y are not specified
			x, y := action.X, action.Y
			
			// If coordinates are not specified, click the center of the element
			if x == 0 && y == 0 {
				x = bounds.X + bounds.Width/2
				y = bounds.Y + bounds.Height/2
				fmt.Printf("Calculated center point: (%d,%d)\n", x, y)
			} else {
				fmt.Printf("Using specified click point: (%d,%d)\n", x, y)
			}
			
			// Update virtual cursor position
			if game != nil {
				game.virtualCursor.x = x
				game.virtualCursor.y = y
				game.virtualCursor.active = true
				fmt.Printf("Moving virtual cursor to: (%d,%d)\n", x, y)
			}
			
			// Add visual delay before clicking to make it visible
			time.Sleep(time.Duration(float64(500 * time.Millisecond) * delayMultiplier))
			
			// Simulate mouse down
			fmt.Printf("Simulating mouse down on %s at (%d,%d)\n", action.Target.ID(), x, y)
			f.Log(fmt.Sprintf("Mouse down on %s at (%d,%d)", action.Target.ID(), x, y))
			action.Target.HandleMouseDown(x, y)
			
			// Record result
			result := fmt.Sprintf("Clicked element %s at (%d,%d)", action.TargetID, x, y)
			f.testCases[f.currentTest].Results = append(f.testCases[f.currentTest].Results, result)
			
			// Update virtual cursor click time
			if game != nil {
				game.virtualCursor.clickTime = time.Now()
			}
			
			// Small delay to simulate real interaction
			time.Sleep(time.Duration(float64(300 * time.Millisecond) * delayMultiplier))
			
			// Simulate mouse up
			fmt.Printf("Simulating mouse up on %s at (%d,%d)\n", action.Target.ID(), x, y)
			f.Log(fmt.Sprintf("Mouse up on %s at (%d,%d)", action.Target.ID(), x, y))
			action.Target.HandleMouseUp(x, y)
			
			// Add time to see the result of the interaction
			time.Sleep(time.Duration(float64(700 * time.Millisecond) * delayMultiplier))
		} else {
			fmt.Printf("Error: Could not find target element '%s'\n", action.TargetID)
			f.Log(fmt.Sprintf("Error: Could not find target element '%s'", action.TargetID))
		}
		
	case "hover":
		// Find the target if needed
		if action.Target == nil && action.TargetID != "" {
			action.Target = f.FindElementByID(action.TargetID)
		}
		
		if action.Target != nil {
			// Get element bounds to calculate center if x,y are not specified
			bounds := action.Target.Bounds()
			x, y := action.X, action.Y
			
			// If coordinates are not specified, hover over the center of the element
			if x == 0 && y == 0 {
				x = bounds.X + bounds.Width/2
				y = bounds.Y + bounds.Height/2
			}
			
			// Update virtual cursor position
			if game != nil {
				game.virtualCursor.x = x
				game.virtualCursor.y = y
				game.virtualCursor.active = true
			}
			
			// Add visual delay before hovering to make it visible
			time.Sleep(time.Duration(float64(300 * time.Millisecond) * delayMultiplier))
			
			// Simulate mouse move
			f.Log(fmt.Sprintf("Mouse move on %s at (%d,%d)", action.Target.ID(), x, y))
			action.Target.HandleMouseMove(x, y)
			
			// Record result
			result := fmt.Sprintf("Hovered over element %s at (%d,%d)", action.TargetID, x, y)
			f.testCases[f.currentTest].Results = append(f.testCases[f.currentTest].Results, result)
		} else {
			f.Log(fmt.Sprintf("Error: Could not find target element %s", action.TargetID))
		}
		
	case "wait":
		// Just wait for the specified duration
		result := fmt.Sprintf("Waited for %v", action.Delay)
		f.testCases[f.currentTest].Results = append(f.testCases[f.currentTest].Results, result)
	}
	
	// Add delay after action
	time.Sleep(time.Duration(float64(action.Delay) * delayMultiplier))
}

// Log adds a message to the log panel
func (f *UITestFrame) Log(message string) {
	// Add timestamp to log
	logEntry := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), message)
	
	// Add to logs
	f.logPanel.logs = append(f.logPanel.logs, logEntry)
	
	// Limit log size
	if len(f.logPanel.logs) > 100 {
		f.logPanel.logs = f.logPanel.logs[len(f.logPanel.logs)-100:]
	}
	
	// Update log display
	f.updateLogDisplay()
}

// updateLogDisplay updates the log labels with current log entries
func (f *UITestFrame) updateLogDisplay() {
	startIndex := len(f.logPanel.logs) - f.logPanel.maxLogs
	if startIndex < 0 {
		startIndex = 0
	}
	
	for i := 0; i < f.logPanel.maxLogs; i++ {
		logIndex := startIndex + i
		if logIndex < len(f.logPanel.logs) {
			f.logPanel.logLabels[i].SetText(f.logPanel.logs[logIndex])
		} else {
			f.logPanel.logLabels[i].SetText("")
		}
	}
}

// NewUITestCase creates a new UI test case
func NewUITestCase(name, description string) *UITestCase {
	return &UITestCase{
		Name:        name,
		Description: description,
		Actions:     make([]UITestAction, 0),
		Results:     make([]string, 0),
	}
}

// AddClickAction adds a click action to a test case
func (tc *UITestCase) AddClickAction(targetID string, x, y int, description string, delay time.Duration) {
	tc.Actions = append(tc.Actions, UITestAction{
		Type:        "click",
		TargetID:    targetID,
		X:           x,
		Y:           y,
		Description: description,
		Delay:       delay,
	})
}

// AddHoverAction adds a hover action to a test case
func (tc *UITestCase) AddHoverAction(targetID string, x, y int, description string, delay time.Duration) {
	tc.Actions = append(tc.Actions, UITestAction{
		Type:        "hover",
		TargetID:    targetID,
		X:           x,
		Y:           y,
		Description: description,
		Delay:       delay,
	})
}

// AddWaitAction adds a wait action to a test case
func (tc *UITestCase) AddWaitAction(duration time.Duration, description string) {
	tc.Actions = append(tc.Actions, UITestAction{
		Type:        "wait",
		Description: description,
		Delay:       duration,
	})
}

// FindElementByID recursively searches for an element with the given ID
func (f *UITestFrame) FindElementByID(id string) components.Element {
	// First check if the ID matches any of the control buttons
	if f.controls != nil {
		// Check control buttons directly
		if f.controls.ID() == id {
			return f.controls
		}
		
		// Search in control buttons children
		if result := findElementByIDRecursive(f.controls, id); result != nil {
			return result
		}
	}
	
	// Search in all direct children including the app container
	for _, child := range f.Children() {
		if child.ID() == id {
			return child
		}
		
		// Search in the child's descendants
		if result := findElementByIDRecursive(child, id); result != nil {
			return result
		}
	}
	
	return nil
}

// findElementByIDRecursive is a helper function that recursively searches for an element with the given ID
func findElementByIDRecursive(root components.Element, id string) components.Element {
	if root.ID() == id {
		return root
	}
	
	for _, child := range root.Children() {
		if found := findElementByIDRecursive(child, id); found != nil {
			return found
		}
	}
	
	return nil
}

// UITestGame implements the Ebiten game interface for running UI tests
type UITestGame struct {
	rootElement   components.Element
	testFrame     *UITestFrame
	renderer      *components.EbitenRenderer
	mouseX, mouseY int
	mousePressed   bool
	clickedButton  string
	virtualCursor  struct {
		x, y      int
		active    bool
		clickTime time.Time
	}
}

// NewUITestGame creates a new UI test game
func NewUITestGame(rootUI components.Element) *UITestGame {
	// Create test frame
	testFrame := NewUITestFrame(rootUI)
	
	// Position the target UI in the top part of the screen (above the test frame)
	targetUIBounds := components.Rect{
		X: 0,
		Y: 0,
		Width: components.ScreenWidth,
		Height: components.ScreenHeight - testFrame.Bounds().Height,
	}
	rootUI.SetBounds(targetUIBounds)
	
	// Add root UI to test frame
	testFrame.AddChild(rootUI)
	
	// Create game
	game := &UITestGame{
		rootElement: testFrame,
		testFrame:   testFrame,
		renderer:    nil,
	}
	
	// Store reference to current game
	currentTestGame = game
	
	return game
}

// Update updates the game state
func (g *UITestGame) Update() error {
	// Get updated mouse position
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	
	// Debug mouse position every few seconds
	if time.Now().Second() % 3 == 0 && time.Now().Nanosecond() < 10000000 {
		fmt.Printf("Mouse position: (%d,%d)\n", g.mouseX, g.mouseY)
		
		// Check if mouse is over play button
		if g.testFrame.controls != nil && g.testFrame.controls.playButton != nil {
			playBounds := g.testFrame.controls.playButton.Bounds()
			if g.mouseX >= playBounds.X && 
			   g.mouseX < playBounds.X + playBounds.Width &&
			   g.mouseY >= playBounds.Y && 
			   g.mouseY < playBounds.Y + playBounds.Height {
				fmt.Printf("Mouse is over PLAY button!\n")
			}
		}
	}
	
	// Handle mouse events
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.mousePressed = true
		fmt.Printf("Mouse DOWN at (%d,%d)\n", g.mouseX, g.mouseY)
		
		// Check for direct button clicks
		if g.testFrame.controls != nil {
			// Play button
			if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.playButton.Bounds()) {
				fmt.Println("DIRECT PLAY BUTTON CLICK")
				if len(g.testFrame.testCases) > 0 {
					g.testFrame.playingTest = true
					g.testFrame.statusLabel.SetText("Running test: " + g.testFrame.testCases[g.testFrame.currentTest].Name)
					g.testFrame.Log("Started test: " + g.testFrame.testCases[g.testFrame.currentTest].Name)
				} else {
					g.testFrame.Log("No test cases to run")
				}
				return nil
			}
			
			// Pause button
			if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.pauseButton.Bounds()) {
				fmt.Println("DIRECT PAUSE BUTTON CLICK")
				if g.testFrame.playingTest {
					g.testFrame.playingTest = false
					g.testFrame.statusLabel.SetText("Paused test: " + g.testFrame.testCases[g.testFrame.currentTest].Name)
					g.testFrame.Log("Paused test")
				}
				return nil
			}
			
			// Step button
			if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.stepButton.Bounds()) {
				fmt.Println("DIRECT STEP BUTTON CLICK")
				if len(g.testFrame.testCases) > 0 {
					g.testFrame.ExecuteNextStep()
				} else {
					g.testFrame.Log("No test cases to run")
				}
				return nil
			}
			
			// Stop button
			if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.stopButton.Bounds()) {
				fmt.Println("DIRECT STOP BUTTON CLICK")
				if g.testFrame.playingTest {
					g.testFrame.playingTest = false
					g.testFrame.statusLabel.SetText("Stopped test: " + g.testFrame.testCases[g.testFrame.currentTest].Name)
					g.testFrame.Log("Stopped test")
				}
				return nil
			}
			
			// Reset button
			if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.resetButton.Bounds()) {
				fmt.Println("DIRECT RESET BUTTON CLICK")
				g.testFrame.ResetTest()
				return nil
			}
			
			// Step mode toggle
			if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.stepModeCheckbox.Bounds()) {
				fmt.Println("DIRECT STEP MODE BUTTON CLICK")
				g.testFrame.controls.stepModeActive = !g.testFrame.controls.stepModeActive
				if g.testFrame.controls.stepModeActive {
					g.testFrame.controls.stepModeCheckbox.SetText("☑ Step Mode")
					g.testFrame.stepMode = true
					g.testFrame.Log("Step mode enabled")
				} else {
					g.testFrame.controls.stepModeCheckbox.SetText("□ Step Mode")
					g.testFrame.stepMode = false
					g.testFrame.Log("Step mode disabled")
				}
				return nil
			}
			
			// Next test button
			if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.nextTestButton.Bounds()) {
				fmt.Println("DIRECT NEXT TEST BUTTON CLICK")
				if len(g.testFrame.testCases) > 0 {
					g.testFrame.currentTest = (g.testFrame.currentTest + 1) % len(g.testFrame.testCases)
					g.testFrame.updateTestCaseLabel()
					g.testFrame.ResetTest()
				}
				return nil
			}
			
			// Previous test button
			if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.prevTestButton.Bounds()) {
				fmt.Println("DIRECT PREV TEST BUTTON CLICK")
				if len(g.testFrame.testCases) > 0 {
					g.testFrame.currentTest = (g.testFrame.currentTest - 1 + len(g.testFrame.testCases)) % len(g.testFrame.testCases)
					g.testFrame.updateTestCaseLabel()
					g.testFrame.ResetTest()
				}
				return nil
			}
		}
		
		// Track which button was pressed for mouse up
		g.clickedButton = ""
		
		// Check if mouse is over any control button
		if g.testFrame.controls != nil {
			// Check all buttons
			if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.playButton.Bounds()) {
				g.clickedButton = "play_button"
			} else if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.pauseButton.Bounds()) {
				g.clickedButton = "pause_button"
			} else if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.stepButton.Bounds()) {
				g.clickedButton = "step_button"
			} else if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.stopButton.Bounds()) {
				g.clickedButton = "stop_button"
			} else if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.resetButton.Bounds()) {
				g.clickedButton = "reset_button"
			} else if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.prevTestButton.Bounds()) {
				g.clickedButton = "prev_test_button"
			} else if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.nextTestButton.Bounds()) {
				g.clickedButton = "next_test_button"
			} else if isPointInRect(g.mouseX, g.mouseY, g.testFrame.controls.stepModeCheckbox.Bounds()) {
				g.clickedButton = "step_mode_button"
			}
			
			if g.clickedButton != "" {
				fmt.Printf("PRESSED ON CONTROL: %s\n", g.clickedButton)
			}
		}
		
		// Always propagate the event to all elements
		g.rootElement.HandleMouseDown(g.mouseX, g.mouseY)
	}
	
	// Handle mouse release events
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.mousePressed {
			g.mousePressed = false
			fmt.Printf("Mouse UP at (%d,%d)\n", g.mouseX, g.mouseY)
			
			// Always propagate the event to all elements
			g.rootElement.HandleMouseUp(g.mouseX, g.mouseY)
			
			// Check if mouse is still over the same control button
			if g.clickedButton != "" {
				var buttonBounds components.Rect
				
				// Get bounds of the clicked button
				switch g.clickedButton {
				case "play_button":
					buttonBounds = g.testFrame.controls.playButton.Bounds()
				case "pause_button":
					buttonBounds = g.testFrame.controls.pauseButton.Bounds()
				case "step_button":
					buttonBounds = g.testFrame.controls.stepButton.Bounds()
				case "stop_button":
					buttonBounds = g.testFrame.controls.stopButton.Bounds()
				case "reset_button":
					buttonBounds = g.testFrame.controls.resetButton.Bounds()
				case "prev_test_button":
					buttonBounds = g.testFrame.controls.prevTestButton.Bounds()
				case "next_test_button":
					buttonBounds = g.testFrame.controls.nextTestButton.Bounds()
				case "step_mode_button":
					buttonBounds = g.testFrame.controls.stepModeCheckbox.Bounds()
				}
				
				// If mouse is still over the button, trigger its action directly
				if isPointInRect(g.mouseX, g.mouseY, buttonBounds) {
					fmt.Printf("TRIGGERING BUTTON ACTION: %s\n", g.clickedButton)
					
					// Trigger the appropriate action
					switch g.clickedButton {
					case "play_button":
						// Call play button handler
						g.testFrame.controls.playButton.HandleMouseUp(g.mouseX, g.mouseY)
					case "pause_button":
						g.testFrame.controls.pauseButton.HandleMouseUp(g.mouseX, g.mouseY)
					case "step_button":
						g.testFrame.controls.stepButton.HandleMouseUp(g.mouseX, g.mouseY)
					case "stop_button":
						g.testFrame.controls.stopButton.HandleMouseUp(g.mouseX, g.mouseY)
					case "reset_button":
						g.testFrame.controls.resetButton.HandleMouseUp(g.mouseX, g.mouseY)
					case "prev_test_button":
						g.testFrame.controls.prevTestButton.HandleMouseUp(g.mouseX, g.mouseY)
					case "next_test_button":
						g.testFrame.controls.nextTestButton.HandleMouseUp(g.mouseX, g.mouseY)
					case "step_mode_button":
						g.testFrame.controls.stepModeCheckbox.HandleMouseUp(g.mouseX, g.mouseY)
					}
				}
				
				g.clickedButton = ""
			}
		}
	}
	
	// Propagate mouse move events
	g.rootElement.HandleMouseMove(g.mouseX, g.mouseY)
	
	// Update test frame
	g.testFrame.Update()
	
	// Update all UI elements (this is redundant with the line above, but keeping it to be safe)
	g.rootElement.Update()
	
	return nil
}

// isPointInRect is a helper function to check if a point is inside a rectangle
func isPointInRect(x, y int, rect components.Rect) bool {
	return x >= rect.X && x < rect.X+rect.Width && y >= rect.Y && y < rect.Y+rect.Height
}

// Draw draws the game
func (g *UITestGame) Draw(screen *ebiten.Image) {
	// Create renderer if needed
	if g.renderer == nil {
		g.renderer = components.NewEbitenRenderer(screen)
	}
	
	// Clear the screen
	g.renderer.Clear(color.RGBA{255, 255, 255, 255})
	
	// Draw all UI elements
	g.rootElement.Draw(g.renderer)
	
	// Draw virtual cursor during test execution
	if g.testFrame.playingTest && g.virtualCursor.active {
		cursorSize := 10
		// Draw cursor circle
		g.renderer.FillCircle(g.virtualCursor.x, g.virtualCursor.y, cursorSize, color.RGBA{255, 0, 0, 200})
		
		// Draw click animation if recently clicked
		if time.Since(g.virtualCursor.clickTime) < 500*time.Millisecond {
			// Calculate size based on time elapsed (shrinking effect)
			size := cursorSize * 2 * int(1.0-float64(time.Since(g.virtualCursor.clickTime))/float64(500*time.Millisecond))
			g.renderer.DrawCircle(g.virtualCursor.x, g.virtualCursor.y, size, color.RGBA{255, 0, 0, 100})
		}
	}
}

// Layout returns the game's screen layout
func (g *UITestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return components.ScreenWidth, components.ScreenHeight
}

// RunUITests runs the UI tests in an interactive window
func RunUITests(targetUI components.Element, testCases []*UITestCase) {
	// Set up Ebiten
	ebiten.SetWindowSize(components.ScreenWidth, components.ScreenHeight)
	ebiten.SetWindowTitle("Finch UI Test Framework")
	
	// Create test game
	game := NewUITestGame(targetUI)
	
	// Add test cases
	for _, tc := range testCases {
		game.testFrame.AddTestCase(tc)
	}
	
	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
} 