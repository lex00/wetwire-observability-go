package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lex00/wetwire-observability-go/testrunner"
)

func testCmd(args []string) int {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	persona := fs.String("persona", "", "Persona to evaluate against (sre, developer, security, beginner)")
	all := fs.Bool("all", false, "Evaluate against all personas")
	format := fs.String("format", "text", "Output format: text or json")
	listPersonas := fs.Bool("list-personas", false, "List available personas")
	threshold := fs.Int("threshold", 0, "Minimum passing score percentage (0-100)")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-obs test [options] <path>")
		fmt.Println()
		fmt.Println("Evaluate observability configurations against personas.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println()
		fmt.Println("Personas:")
		fmt.Println("  sre        Site Reliability Engineer - reliability, alerting, SLOs")
		fmt.Println("  developer  Developer - debugging, application metrics, tracing")
		fmt.Println("  security   Security Analyst - auth, compliance, threat detection")
		fmt.Println("  beginner   Beginner - basic monitoring setup")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  wetwire-obs test --persona sre ./monitoring")
		fmt.Println("  wetwire-obs test --all --format json ./monitoring")
		fmt.Println("  wetwire-obs test --persona beginner --threshold 70 ./monitoring")
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	// Handle list-personas
	if *listPersonas {
		printPersonas()
		return 0
	}

	// Get path
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: no path provided")
		fmt.Fprintln(os.Stderr, "Usage: wetwire-obs test [options] <path>")
		return 1
	}
	path := fs.Arg(0)

	// Create runner
	runner := testrunner.NewRunner()

	if *all {
		runner.WithAllPersonas()
	} else if *persona != "" {
		p := testrunner.GetPersona(*persona)
		if p == nil {
			fmt.Fprintf(os.Stderr, "Error: unknown persona %q\n", *persona)
			fmt.Fprintln(os.Stderr, "Available personas: sre, developer, security, beginner")
			return 1
		}
		runner.WithPersona(*persona)
	} else {
		// Default to beginner
		runner.WithPersona("beginner")
	}

	// Evaluate
	result, err := runner.Evaluate(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Output
	switch *format {
	case "json":
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	default:
		printTextResult(result)
	}

	// Check threshold
	if *threshold > 0 && result.Percentage < float64(*threshold) {
		fmt.Fprintf(os.Stderr, "\nFailed: score %.1f%% below threshold %d%%\n",
			result.Percentage, *threshold)
		return 1
	}

	return 0
}

func printPersonas() {
	fmt.Println("Available Personas:")
	fmt.Println()

	personas := testrunner.GetAllPersonas()
	for _, p := range personas {
		fmt.Printf("  %-12s %s\n", p.ID, p.Description)
		fmt.Printf("               Criteria: %d\n", len(p.Criteria))
		fmt.Println()
	}
}

func printTextResult(result *testrunner.Result) {
	fmt.Println()
	fmt.Println("Evaluation Results")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	for _, pr := range result.PersonaResults {
		fmt.Printf("Persona: %s\n", pr.PersonaName)
		fmt.Printf("Score: %d/%d (%.1f%%)\n", pr.Score, pr.MaxScore, pr.Percentage)
		fmt.Println()

		// Group criteria by status
		var passed, partial, failed, skipped []testrunner.CriterionResult
		for _, cr := range pr.Criteria {
			switch cr.Status {
			case testrunner.StatusPass:
				passed = append(passed, cr)
			case testrunner.StatusPartial:
				partial = append(partial, cr)
			case testrunner.StatusFail:
				failed = append(failed, cr)
			case testrunner.StatusSkip:
				skipped = append(skipped, cr)
			}
		}

		if len(passed) > 0 {
			fmt.Println("  PASS:")
			for _, cr := range passed {
				fmt.Printf("    [+] %s\n", cr.Name)
			}
		}

		if len(partial) > 0 {
			fmt.Println("  PARTIAL:")
			for _, cr := range partial {
				fmt.Printf("    [~] %s: %s\n", cr.Name, cr.Message)
			}
		}

		if len(failed) > 0 {
			fmt.Println("  FAIL:")
			for _, cr := range failed {
				fmt.Printf("    [-] %s: %s\n", cr.Name, cr.Message)
			}
		}

		if len(skipped) > 0 {
			fmt.Println("  SKIP:")
			for _, cr := range skipped {
				fmt.Printf("    [?] %s\n", cr.Name)
			}
		}

		fmt.Println()
	}

	// Overall
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Overall Score: %d/%d (%.1f%%)\n",
		result.TotalScore, result.MaxScore, result.Percentage)

	if len(result.Recommendations) > 0 {
		fmt.Println()
		fmt.Println("Recommendations:")
		for _, rec := range result.Recommendations {
			fmt.Printf("  - %s\n", rec)
		}
	}
}
