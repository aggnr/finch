package test

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/aggnr/finch/components"
)

// DOMTestGame is a game that runs DOM-based UI tests
type DOMTestGame struct {
	rootElement        components.Element
	domInspector       *components.DOMInspector
	testManager        *DOMTestManager
	renderer           *components.EbitenRenderer
	mouseX, mouseY     int
	mousePressed       bool
	inspectorEnabled   bool
	recordingEnabled   bool
	recordedTestCase   *DOMTestCase
	recordStartTime    time.Time
	lastEventTime      time.Time
}

// DOMTestManager manages DOM-based test cases
type DOMTestManager struct {
	*components.BaseElement
	testCases       []*DOMTestCase
	currentTest     int
	playingTest     bool
	stepMode        bool
	currentStep     int
	speedMultiplier float64
	logPanel        *components.DOMNode
	statusLabel     *components.Label
	testResult      *components.Label
	controls        *TestControls
}

// DOMTestCase represents a test case for the DOM-based UI
type DOMTestCase struct {
	Name        string
	Description string
	Actions     []DOMTestAction
	Results     []string
}

// DOMTestAction represents an action in a DOM test case
type DOMTestAction struct {
	Type             string                 // "click", "hover", "wait", "assertValue", etc.
	Selector         string                 // DOM selector for target element
	SelectorType     string                 // "id", "class", "tag", "xpath"
	Target           components.DOMElement  // Reference to target element
	X, Y             int                    // Coordinates for actions like click
	RelativePosition bool                   // Whether coordinates are relative to target element
	Value            string                 // Value for input actions
	Description      string                 // Human-readable description
	Delay            time.Duration          // Delay after action
}

// NewDOMTestGame creates a new game for running DOM-based tests
func NewDOMTestGame(rootUI components.Element) *DOMTestGame {
	game := &DOMTestGame{
		rootElement:      rootUI,
		mouseX:           0,
		mouseY:           0,
		mousePressed:     false,
		inspectorEnabled: false,
		recordingEnabled: false,
	}
	
	// Create test manager
	game.testManager = NewDOMTestManager()
	
	// Create DOM inspector (if rootUI is a DOMElement)
	if domRoot, ok := rootUI.(components.DOMElement); ok {
		game.domInspector = components.NewDOMInspector(domRoot)
	}
	
	// Store global reference
	currentTestGame = game
	
	return game
}

// NewDOMTestManager creates a new test manager
func NewDOMTestManager() *DOMTestManager {
	manager := &DOMTestManager{
		BaseElement:     components.NewBaseElement("dom_test_manager"),
		testCases:       make([]*DOMTestCase, 0),
		currentTest:     0,
		playingTest:     false,
		stepMode:        false,
		currentStep:     -1,
		speedMultiplier: 1.0,
	}
	
	// Create test manager UI
	manager.SetBounds(components.Rect{
		X:      0, 
		Y:      components.ScreenHeight - 200, 
		Width:  components.ScreenWidth, 
		Height: 200,
	})
	
	// Create status label
	manager.statusLabel = components.NewLabel("status_label", "DOM Test Framework Ready", 14, color.RGBA{0, 0, 0, 255})
	manager.statusLabel.SetBounds(components.Rect{X: 10, Y: 60, Width: 400, Height: 20})
	manager.AddChild(manager.statusLabel)
	
	// Create test result label
	manager.testResult = components.NewLabel("test_result", "", 14, color.RGBA{0, 100, 0, 255})
	manager.testResult.SetBounds(components.Rect{X: 10, Y: 80, Width: 400, Height: 20})
	manager.AddChild(manager.testResult)
	
	// Create test controls
	manager.controls = createTestControls(nil) // We'll set this up properly later
	manager.controls.SetBounds(components.Rect{X: 10, Y: 10, Width: components.ScreenWidth - 20, Height: 40})
	manager.AddChild(manager.controls)
	
	return manager
}

// NewDOMTestCase creates a new DOM test case
func NewDOMTestCase(name, description string) *DOMTestCase {
	return &DOMTestCase{
		Name:        name,
		Description: description,
		Actions:     make([]DOMTestAction, 0),
		Results:     make([]string, 0),
	}
}

// AddClickAction adds a click action to a DOM test case using a selector
func (tc *DOMTestCase) AddClickAction(selector string, x, y int, relativePos bool, description string, delay time.Duration) {
	tc.Actions = append(tc.Actions, DOMTestAction{
		Type:             "click",
		Selector:         selector,
		SelectorType:     getSelectorType(selector),
		X:                x,
		Y:                y,
		RelativePosition: relativePos,
		Description:      description,
		Delay:            delay,
	})
}

// AddHoverAction adds a hover action to a DOM test case using a selector
func (tc *DOMTestCase) AddHoverAction(selector string, x, y int, relativePos bool, description string, delay time.Duration) {
	tc.Actions = append(tc.Actions, DOMTestAction{
		Type:             "hover",
		Selector:         selector,
		SelectorType:     getSelectorType(selector),
		X:                x,
		Y:                y,
		RelativePosition: relativePos,
		Description:      description,
		Delay:            delay,
	})
}

// AddWaitAction adds a wait action to a DOM test case
func (tc *DOMTestCase) AddWaitAction(duration time.Duration, description string) {
	tc.Actions = append(tc.Actions, DOMTestAction{
		Type:        "wait",
		Description: description,
		Delay:       duration,
	})
}

// AddInputAction adds an input action to a DOM test case using a selector
func (tc *DOMTestCase) AddInputAction(selector string, value string, description string, delay time.Duration) {
	tc.Actions = append(tc.Actions, DOMTestAction{
		Type:         "input",
		Selector:     selector,
		SelectorType: getSelectorType(selector),
		Value:        value,
		Description:  description,
		Delay:        delay,
	})
}

// AddAssertAction adds an assertion action to a DOM test case
func (tc *DOMTestCase) AddAssertAction(selector string, expectedValue string, description string) {
	tc.Actions = append(tc.Actions, DOMTestAction{
		Type:         "assert",
		Selector:     selector,
		SelectorType: getSelectorType(selector),
		Value:        expectedValue,
		Description:  description,
	})
}

// getSelectorType determines the type of selector
func getSelectorType(selector string) string {
	if len(selector) == 0 {
		return ""
	}
	
	switch selector[0] {
	case '#':
		return "id"
	case '.':
		return "class"
	case '/':
		return "xpath"
	default:
		return "tag"
	}
}

// Update updates the test game
func (g *DOMTestGame) Update() error {
	// Update mouse position
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	
	// Toggle inspector with I key
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		g.inspectorEnabled = !g.inspectorEnabled
		if g.domInspector != nil {
			if g.inspectorEnabled {
				g.domInspector.Enable()
			} else {
				g.domInspector.Disable()
			}
		}
	}
	
	// Toggle test recording with R key
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.recordingEnabled = !g.recordingEnabled
		if g.recordingEnabled {
			// Start recording a new test case
			g.recordedTestCase = NewDOMTestCase("Recorded Test", "Test case recorded from user interactions")
			g.recordStartTime = time.Now()
			g.lastEventTime = g.recordStartTime
		} else if g.recordedTestCase != nil && len(g.recordedTestCase.Actions) > 0 {
			// Add the recorded test case to the test manager
			g.testManager.AddTestCase(g.recordedTestCase)
			g.recordedTestCase = nil
		}
	}
	
	// If inspector is enabled, let it handle input first
	if g.inspectorEnabled && g.domInspector != nil {
		g.domInspector.HandleMouseMove(g.mouseX, g.mouseY)
		
		// Handle mouse press events
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.mousePressed = true
			g.domInspector.HandleMouseDown(g.mouseX, g.mouseY)
			return nil
		}
		
		// Handle mouse release events
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			if g.mousePressed {
				g.mousePressed = false
				g.domInspector.HandleMouseUp(g.mouseX, g.mouseY)
			}
			return nil
		}
	}
	
	// Propagate mouse move events
	g.rootElement.HandleMouseMove(g.mouseX, g.mouseY)
	
	// Handle mouse press events
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.mousePressed = true
		g.rootElement.HandleMouseDown(g.mouseX, g.mouseY)
		
		// If recording, add this click as an action
		if g.recordingEnabled && g.recordedTestCase != nil {
			// Calculate delay since last event
			delay := time.Since(g.lastEventTime)
			g.lastEventTime = time.Now()
			
			// Try to identify the clicked element
			var selector string
			var relativePos bool
			
			// If we have the inspector, use it to find the element
			if g.domInspector != nil {
				if element := g.domInspector.GetSelectedElement(); element != nil {
					selector = "#" + element.ID() // Use ID selector for simplicity
					
					// Calculate if using relative positioning makes sense
					bounds := element.ComputedBounds()
					if bounds.Width > 0 && bounds.Height > 0 {
						// Use relative position
						relativePos = true
						
						// Add click action with relative coordinates
						g.recordedTestCase.AddClickAction(
							selector,
							g.mouseX - bounds.X,
							g.mouseY - bounds.Y,
							true,
							fmt.Sprintf("Click on %s", element.ID()),
							delay,
						)
						return nil
					}
				}
			}
			
			// Fallback: Use absolute coordinates
			g.recordedTestCase.AddClickAction(
				"",
				g.mouseX,
				g.mouseY,
				false,
				fmt.Sprintf("Click at position (%d, %d)", g.mouseX, g.mouseY),
				delay,
			)
		}
	}
	
	// Handle mouse release events
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.mousePressed {
			g.mousePressed = false
			g.rootElement.HandleMouseUp(g.mouseX, g.mouseY)
		}
	}
	
	// Update test manager
	if g.testManager.playingTest {
		g.testManager.Update()
	}
	
	// Update all UI elements
	g.rootElement.Update()
	
	return nil
}

// Draw draws the test game
func (g *DOMTestGame) Draw(screen *ebiten.Image) {
	// Create renderer if needed
	if g.renderer == nil {
		g.renderer = components.NewEbitenRenderer(screen)
	}
	
	// Clear the screen
	g.renderer.Clear(color.RGBA{255, 255, 255, 255})
	
	// Draw all UI elements
	g.rootElement.Draw(g.renderer)
	
	// Draw test manager
	g.testManager.Draw(g.renderer)
	
	// Draw DOM inspector if enabled
	if g.inspectorEnabled && g.domInspector != nil {
		g.domInspector.Draw(g.renderer)
	}
	
	// Draw recording indicator
	if g.recordingEnabled {
		recordingText := fmt.Sprintf("RECORDING TEST (%d actions)", len(g.recordedTestCase.Actions))
		g.renderer.DrawText(recordingText, 10, 10, color.RGBA{255, 0, 0, 255}, 14)
	}
}

// Layout returns the game's screen layout
func (g *DOMTestGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return components.ScreenWidth, components.ScreenHeight
}

// AddTestCase adds a test case to the test manager
func (tm *DOMTestManager) AddTestCase(testCase *DOMTestCase) {
	tm.testCases = append(tm.testCases, testCase)
	tm.Log(fmt.Sprintf("Added test case: %s", testCase.Name))
}

// Log adds a message to the log panel
func (tm *DOMTestManager) Log(message string) {
	fmt.Println("[DOMTestManager]", message)
	// In a full implementation, this would add to a visible log panel
}

// Update updates the test manager
func (tm *DOMTestManager) Update() {
	// If not playing a test, return
	if !tm.playingTest {
		return
	}
	
	// Execute the next step
	tm.ExecuteNextStep()
}

// ExecuteNextStep executes the next step in the current test case
func (tm *DOMTestManager) ExecuteNextStep() {
	// Check if we have test cases
	if len(tm.testCases) == 0 {
		tm.Log("No test cases to run")
		return
	}
	
	// Get current test case
	testCase := tm.testCases[tm.currentTest]
	
	// Move to next step
	tm.currentStep++
	
	// Check if test is complete
	if tm.currentStep >= len(testCase.Actions) {
		tm.statusLabel.SetText("Test completed: " + testCase.Name)
		tm.testResult.SetText("Test Passed!")
		tm.testResult.SetTextColor(color.RGBA{0, 128, 0, 255})
		tm.Log("Test completed successfully")
		
		// If in step mode, don't auto-advance; wait for next button click
		if tm.stepMode {
			tm.currentStep = -1
			tm.playingTest = false
			return
		}
		
		// Move to next test case
		tm.currentTest++
		
		// If all tests are complete, stop playing
		if tm.currentTest >= len(tm.testCases) {
			tm.currentTest = 0
			tm.playingTest = false
			tm.Log("All test cases completed")
			tm.statusLabel.SetText("All test cases completed")
			return
		}
		
		// Reset step counter for next test
		tm.currentStep = -1
		
		// Show status for next test
		tm.statusLabel.SetText("Running test: " + tm.testCases[tm.currentTest].Name)
		tm.testResult.SetText("")
		
		return
	}
	
	// Get the current action
	action := testCase.Actions[tm.currentStep]
	
	// Show action description
	tm.statusLabel.SetText(action.Description)
	
	// Execute the action based on its type
	switch action.Type {
	case "click":
		tm.executeClickAction(action)
	case "hover":
		tm.executeHoverAction(action)
	case "wait":
		tm.executeWaitAction(action)
	case "input":
		tm.executeInputAction(action)
	case "assert":
		tm.executeAssertAction(action)
	default:
		tm.Log(fmt.Sprintf("Unknown action type: %s", action.Type))
	}
}

// executeClickAction executes a click action
func (tm *DOMTestManager) executeClickAction(action DOMTestAction) {
	// Find the target element if we have a selector
	if action.Selector != "" {
		// TODO: Implement selector lookup to find the element
		
		// For now, just log the action
		tm.Log(fmt.Sprintf("Would click element with selector: %s", action.Selector))
	} else {
		// Click at absolute coordinates
		tm.Log(fmt.Sprintf("Would click at position: (%d, %d)", action.X, action.Y))
	}
	
	// In a full implementation, this would:
	// 1. Find the target element using the selector
	// 2. Calculate the click position (relative or absolute)
	// 3. Simulate a mouse click on that position
	// 4. Wait for the specified delay
}

// executeHoverAction executes a hover action
func (tm *DOMTestManager) executeHoverAction(action DOMTestAction) {
	// Similar to executeClickAction, but simulates hover
	tm.Log(fmt.Sprintf("Would hover over selector: %s", action.Selector))
}

// executeWaitAction executes a wait action
func (tm *DOMTestManager) executeWaitAction(action DOMTestAction) {
	tm.Log(fmt.Sprintf("Would wait for: %v", action.Delay))
	
	// In a full implementation, this would create a timer that pauses execution
}

// executeInputAction executes an input action
func (tm *DOMTestManager) executeInputAction(action DOMTestAction) {
	tm.Log(fmt.Sprintf("Would input '%s' into selector: %s", action.Value, action.Selector))
}

// executeAssertAction executes an assertion action
func (tm *DOMTestManager) executeAssertAction(action DOMTestAction) {
	tm.Log(fmt.Sprintf("Would assert '%s' for selector: %s", action.Value, action.Selector))
}

// RunDOMTests runs UI tests using the DOM-based test framework
func RunDOMTests(targetUI components.Element, testCases []*DOMTestCase) {
	// Set up Ebiten
	ebiten.SetWindowSize(components.ScreenWidth, components.ScreenHeight)
	ebiten.SetWindowTitle("Finch UI DOM Test Framework")
	
	// Create test game
	game := NewDOMTestGame(targetUI)
	
	// Add test cases
	for _, tc := range testCases {
		game.testManager.AddTestCase(tc)
	}
	
	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		fmt.Printf("Error running game: %v", err)
	}
} 