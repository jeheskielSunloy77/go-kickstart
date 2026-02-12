package ui

import "testing"

func TestCopyConstantsAreNonEmpty(t *testing.T) {
	values := []string{
		WelcomeTitle,
		WelcomeNextLabel,
		FlowTitle,
		FlowDescription,
		ProjectNameTitle,
		ModuleTitle,
		DestinationTitle,
		ReviewTitle,
		ReviewActionTitle,
		ReviewSummaryHeading,
		OverwriteTitle,
		ContributionLine(),
	}

	for i, value := range values {
		if value == "" {
			t.Fatalf("copy value at index %d should not be empty", i)
		}
	}
}

func TestContributionLineIsDeterministic(t *testing.T) {
	first := ContributionLine()
	second := ContributionLine()
	if first != second {
		t.Fatalf("contribution line should be deterministic")
	}
}
