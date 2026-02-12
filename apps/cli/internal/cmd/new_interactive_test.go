package cmd

import (
	"reflect"
	"testing"

	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/prompts"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/validate"
)

func restoreInteractiveDeps() {
	chooseFlowFn = prompts.ChooseFlow
	basicFlowFn = prompts.BasicFlow
	destinationFlowFn = prompts.DestinationFlow
	componentsFlowFn = prompts.ComponentsFlow
	databaseFlowFn = prompts.DatabaseFlow
	storageFlowFn = prompts.StorageFlow
	reviewConfigFn = prompts.ReviewConfig
	validateProjectNameFn = validate.ProjectName
	validateModulePathFn = validate.ModulePath
	resolveProjectDestinationFn = validate.ResolveProjectDestination
	isNonEmptyDirFn = validate.IsNonEmptyDir
	showWelcomeFn = showWelcomeScreen
	runWithSpinnerFn = ui.RunWithSpinner
	scaffoldProjectFn = scaffold.ScaffoldProject
	printSummaryFn = ui.PrintSummary
}

func TestRunInteractiveBasicSkipsAdvancedPrompts(t *testing.T) {
	t.Cleanup(restoreInteractiveDeps)

	destinationCalled := false
	componentsCalled := false
	databaseCalled := false
	storageCalled := false

	showWelcomeFn = func() error { return nil }
	chooseFlowFn = func() (prompts.FlowChoice, error) { return prompts.FlowBasic, nil }
	basicFlowFn = func(cfg scaffold.ScaffoldConfiguration) (scaffold.ScaffoldConfiguration, error) {
		cfg.ProjectName = "demo"
		cfg.ModulePath = "github.com/acme/demo"
		return cfg, nil
	}
	destinationFlowFn = func(cfg *scaffold.ScaffoldConfiguration) error {
		destinationCalled = true
		return nil
	}
	componentsFlowFn = func(cfg *scaffold.ScaffoldConfiguration) error {
		componentsCalled = true
		return nil
	}
	databaseFlowFn = func(cfg *scaffold.ScaffoldConfiguration) error {
		databaseCalled = true
		return nil
	}
	storageFlowFn = func(cfg *scaffold.ScaffoldConfiguration) error {
		storageCalled = true
		return nil
	}
	reviewConfigFn = func(cfg scaffold.ScaffoldConfiguration) (prompts.ReviewAction, error) {
		return prompts.ReviewGenerate, nil
	}
	validateProjectNameFn = func(string) error { return nil }
	validateModulePathFn = func(string) error { return nil }
	resolveProjectDestinationFn = func(baseArg, projectName string) (string, error) { return "/tmp/demo", nil }
	isNonEmptyDirFn = func(path string) (bool, error) { return false, nil }
	runWithSpinnerFn = func(_ string, fn func() error) error { return fn() }
	scaffoldProjectFn = func(scaffold.ScaffoldConfiguration, bool) error { return nil }
	printSummaryFn = func(string) {}

	if err := runInteractive(); err != nil {
		t.Fatalf("runInteractive returned error: %v", err)
	}

	if destinationCalled {
		t.Fatalf("destination flow should not run in basic mode")
	}
	if componentsCalled || databaseCalled || storageCalled {
		t.Fatalf("advanced flows should not run in basic mode")
	}
}

func TestRunInteractiveAdvancedFlowOrder(t *testing.T) {
	t.Cleanup(restoreInteractiveDeps)

	var calls []string

	showWelcomeFn = func() error {
		calls = append(calls, "welcome")
		return nil
	}
	chooseFlowFn = func() (prompts.FlowChoice, error) {
		calls = append(calls, "choose")
		return prompts.FlowAdvanced, nil
	}
	basicFlowFn = func(cfg scaffold.ScaffoldConfiguration) (scaffold.ScaffoldConfiguration, error) {
		calls = append(calls, "basic")
		cfg.ProjectName = "demo"
		cfg.ModulePath = "github.com/acme/demo"
		return cfg, nil
	}
	destinationFlowFn = func(cfg *scaffold.ScaffoldConfiguration) error {
		calls = append(calls, "destination")
		return nil
	}
	componentsFlowFn = func(cfg *scaffold.ScaffoldConfiguration) error {
		calls = append(calls, "components")
		return nil
	}
	databaseFlowFn = func(cfg *scaffold.ScaffoldConfiguration) error {
		calls = append(calls, "database")
		return nil
	}
	storageFlowFn = func(cfg *scaffold.ScaffoldConfiguration) error {
		calls = append(calls, "storage")
		return nil
	}
	reviewConfigFn = func(cfg scaffold.ScaffoldConfiguration) (prompts.ReviewAction, error) {
		calls = append(calls, "review")
		return prompts.ReviewGenerate, nil
	}
	validateProjectNameFn = func(string) error { return nil }
	validateModulePathFn = func(string) error { return nil }
	resolveProjectDestinationFn = func(baseArg, projectName string) (string, error) { return "/tmp/demo", nil }
	isNonEmptyDirFn = func(path string) (bool, error) { return false, nil }
	runWithSpinnerFn = func(_ string, fn func() error) error { return fn() }
	scaffoldProjectFn = func(scaffold.ScaffoldConfiguration, bool) error { return nil }
	printSummaryFn = func(string) {}

	if err := runInteractive(); err != nil {
		t.Fatalf("runInteractive returned error: %v", err)
	}

	want := []string{"welcome", "choose", "basic", "destination", "components", "database", "storage", "review"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("unexpected call order: got %v, want %v", calls, want)
	}
}
