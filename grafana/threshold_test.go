package grafana

import "testing"

func TestBaseStep(t *testing.T) {
	step := BaseStep("green")
	if step.Value != nil {
		t.Errorf("BaseStep value should be nil, got %v", step.Value)
	}
	if step.Color != "green" {
		t.Errorf("BaseStep color = %q, want green", step.Color)
	}
}

func TestStep(t *testing.T) {
	step := Step(50, "yellow")
	if step.Value == nil || *step.Value != 50 {
		t.Errorf("Step value = %v, want 50", step.Value)
	}
	if step.Color != "yellow" {
		t.Errorf("Step color = %q, want yellow", step.Color)
	}
}

func TestAbsoluteThresholds(t *testing.T) {
	thresholds := AbsoluteThresholds(
		BaseStep("green"),
		Step(50, "yellow"),
		Step(80, "red"),
	)

	if thresholds.Mode != ThresholdModeAbsolute {
		t.Errorf("Mode = %q, want absolute", thresholds.Mode)
	}
	if len(thresholds.Steps) != 3 {
		t.Errorf("len(Steps) = %d, want 3", len(thresholds.Steps))
	}
	if thresholds.Steps[0].Color != "green" {
		t.Errorf("Steps[0].Color = %q, want green", thresholds.Steps[0].Color)
	}
	if *thresholds.Steps[1].Value != 50 {
		t.Errorf("Steps[1].Value = %v, want 50", *thresholds.Steps[1].Value)
	}
}

func TestPercentageThresholds(t *testing.T) {
	thresholds := PercentageThresholds(
		BaseStep("green"),
		Step(50, "yellow"),
		Step(80, "red"),
	)

	if thresholds.Mode != ThresholdModePercentage {
		t.Errorf("Mode = %q, want percentage", thresholds.Mode)
	}
}

func TestGreenYellowRed(t *testing.T) {
	thresholds := GreenYellowRed(50, 80)

	if len(thresholds.Steps) != 3 {
		t.Errorf("len(Steps) = %d, want 3", len(thresholds.Steps))
	}
	if thresholds.Steps[0].Color != ColorGreen {
		t.Errorf("Steps[0].Color = %q, want green", thresholds.Steps[0].Color)
	}
	if *thresholds.Steps[1].Value != 50 {
		t.Errorf("Steps[1].Value = %v, want 50", *thresholds.Steps[1].Value)
	}
	if thresholds.Steps[1].Color != ColorYellow {
		t.Errorf("Steps[1].Color = %q, want yellow", thresholds.Steps[1].Color)
	}
	if *thresholds.Steps[2].Value != 80 {
		t.Errorf("Steps[2].Value = %v, want 80", *thresholds.Steps[2].Value)
	}
	if thresholds.Steps[2].Color != ColorRed {
		t.Errorf("Steps[2].Color = %q, want red", thresholds.Steps[2].Color)
	}
}

func TestRedYellowGreen(t *testing.T) {
	thresholds := RedYellowGreen(20, 80)

	if len(thresholds.Steps) != 3 {
		t.Errorf("len(Steps) = %d, want 3", len(thresholds.Steps))
	}
	if thresholds.Steps[0].Color != ColorRed {
		t.Errorf("Steps[0].Color = %q, want red", thresholds.Steps[0].Color)
	}
	if thresholds.Steps[1].Color != ColorYellow {
		t.Errorf("Steps[1].Color = %q, want yellow", thresholds.Steps[1].Color)
	}
	if thresholds.Steps[2].Color != ColorGreen {
		t.Errorf("Steps[2].Color = %q, want green", thresholds.Steps[2].Color)
	}
}

func TestSingleColor(t *testing.T) {
	thresholds := SingleColor("blue")

	if len(thresholds.Steps) != 1 {
		t.Errorf("len(Steps) = %d, want 1", len(thresholds.Steps))
	}
	if thresholds.Steps[0].Color != "blue" {
		t.Errorf("Steps[0].Color = %q, want blue", thresholds.Steps[0].Color)
	}
}
