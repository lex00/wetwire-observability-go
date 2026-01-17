// Command test runs persona-based testing of observability configurations.
package main

import (
	"encoding/json"
	"fmt"

	"github.com/lex00/wetwire-observability-go/testrunner"
	"github.com/spf13/cobra"
)

func newTestCmd() *cobra.Command {
	var (
		persona      string
		all          bool
		format       string
		listPersonas bool
		threshold    int
	)

	cmd := &cobra.Command{
		Use:   "test <path>",
		Short: "Evaluate observability configurations against personas",
		Long: `Test evaluates observability configurations against different personas.

Available personas:
  - sre: Site Reliability Engineer - reliability, alerting, SLOs
  - developer: Developer - debugging, application metrics, tracing
  - security: Security Analyst - auth, compliance, threat detection
  - beginner: Beginner - basic monitoring setup

Examples:
  wetwire-obs test --persona sre ./monitoring
  wetwire-obs test --all --format json ./monitoring
  wetwire-obs test --persona beginner --threshold 70 ./monitoring`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if listPersonas {
				printPersonas()
				return nil
			}
			return runTest(args[0], persona, all, format, threshold)
		},
	}

	cmd.Flags().StringVarP(&persona, "persona", "p", "", "Persona to evaluate against")
	cmd.Flags().BoolVar(&all, "all", false, "Evaluate against all personas")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")
	cmd.Flags().BoolVar(&listPersonas, "list-personas", false, "List available personas")
	cmd.Flags().IntVar(&threshold, "threshold", 0, "Minimum passing score percentage (0-100)")

	return cmd
}

func runTest(path, persona string, all bool, format string, threshold int) error {
	// Create runner
	runner := testrunner.NewRunner()

	if all {
		runner.WithAllPersonas()
	} else if persona != "" {
		p := testrunner.GetPersona(persona)
		if p == nil {
			return fmt.Errorf("unknown persona %q\nAvailable: sre, developer, security, beginner", persona)
		}
		runner.WithPersona(persona)
	} else {
		// Default to beginner
		runner.WithPersona("beginner")
	}

	// Evaluate
	result, err := runner.Evaluate(path)
	if err != nil {
		return err
	}

	// Output
	switch format {
	case "json":
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	default:
		printTextResult(result)
	}

	// Check threshold
	if threshold > 0 && result.Percentage < float64(threshold) {
		return fmt.Errorf("score %.1f%% below threshold %d%%", result.Percentage, threshold)
	}

	return nil
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
	fmt.Println("==================================================")
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
	fmt.Println("--------------------------------------------------")
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
