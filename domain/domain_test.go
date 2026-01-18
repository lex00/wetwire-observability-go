package domain

import (
	"testing"

	coredomain "github.com/lex00/wetwire-core-go/domain"
)

func TestObservabilityDomainImplementsInterface(t *testing.T) {
	// Compile-time check that ObservabilityDomain implements Domain
	var _ coredomain.Domain = (*ObservabilityDomain)(nil)
}

func TestObservabilityDomainImplementsListerDomain(t *testing.T) {
	// Compile-time check that ObservabilityDomain implements ListerDomain
	var _ coredomain.ListerDomain = (*ObservabilityDomain)(nil)
}

func TestObservabilityDomainName(t *testing.T) {
	d := &ObservabilityDomain{}
	if d.Name() != "observability" {
		t.Errorf("expected name 'observability', got %q", d.Name())
	}
}

func TestObservabilityDomainVersion(t *testing.T) {
	d := &ObservabilityDomain{}
	v := d.Version()
	if v == "" {
		t.Error("version should not be empty")
	}
}

func TestObservabilityDomainBuilder(t *testing.T) {
	d := &ObservabilityDomain{}
	b := d.Builder()
	if b == nil {
		t.Error("builder should not be nil")
	}
}

func TestObservabilityDomainLinter(t *testing.T) {
	d := &ObservabilityDomain{}
	l := d.Linter()
	if l == nil {
		t.Error("linter should not be nil")
	}
}

func TestObservabilityDomainInitializer(t *testing.T) {
	d := &ObservabilityDomain{}
	i := d.Initializer()
	if i == nil {
		t.Error("initializer should not be nil")
	}
}

func TestObservabilityDomainValidator(t *testing.T) {
	d := &ObservabilityDomain{}
	v := d.Validator()
	if v == nil {
		t.Error("validator should not be nil")
	}
}

func TestObservabilityDomainLister(t *testing.T) {
	d := &ObservabilityDomain{}
	l := d.Lister()
	if l == nil {
		t.Error("lister should not be nil")
	}
}

func TestCreateRootCommand(t *testing.T) {
	cmd := CreateRootCommand(&ObservabilityDomain{})
	if cmd == nil {
		t.Fatal("root command should not be nil")
	}
	if cmd.Use != "wetwire-observability" {
		t.Errorf("expected Use 'wetwire-observability', got %q", cmd.Use)
	}
}

func TestLinterWithDisableOption(t *testing.T) {
	d := &ObservabilityDomain{}
	linter := d.Linter()

	// Test that disabled rules are respected
	ctx := &Context{}
	result, err := linter.Lint(ctx, ".", LintOpts{
		Disable: []string{"WOB001", "WOB002"},
	})

	if err != nil {
		t.Fatalf("lint with disable option failed: %v", err)
	}
	if result == nil {
		t.Fatal("lint result should not be nil")
	}
}

func TestLinterWithFixOption(t *testing.T) {
	d := &ObservabilityDomain{}
	linter := d.Linter()

	// Test that fix option is handled
	ctx := &Context{}
	result, err := linter.Lint(ctx, ".", LintOpts{
		Fix: true,
	})

	if err != nil {
		t.Fatalf("lint with fix option failed: %v", err)
	}
	if result == nil {
		t.Fatal("lint result should not be nil")
	}
}

func TestLinterWithFixAndDisableOptions(t *testing.T) {
	d := &ObservabilityDomain{}
	linter := d.Linter()

	// Test both options together
	ctx := &Context{}
	result, err := linter.Lint(ctx, ".", LintOpts{
		Fix:     true,
		Disable: []string{"WOB001"},
	})

	if err != nil {
		t.Fatalf("lint with fix and disable options failed: %v", err)
	}
	if result == nil {
		t.Fatal("lint result should not be nil")
	}
}
