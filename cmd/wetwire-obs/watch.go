// Command watch monitors source files and auto-rebuilds on changes.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	coredomain "github.com/lex00/wetwire-core-go/domain"
	"github.com/lex00/wetwire-observability-go/domain"
	"github.com/spf13/cobra"
)

func newWatchCmd() *cobra.Command {
	var (
		lintOnly bool
		debounce time.Duration
		output   string
	)

	cmd := &cobra.Command{
		Use:   "watch [path]",
		Short: "Monitor source files for changes and automatically rebuild",
		Long: `Watch monitors source files for changes and automatically rebuilds.

The watch command:
  - Monitors the source directory for .go file changes
  - Runs lint on each change
  - Rebuilds if lint passes (unless --lint-only)
  - Debounces rapid changes to avoid excessive rebuilds

Examples:
  wetwire-obs watch ./monitoring
  wetwire-obs watch --lint-only ./monitoring
  wetwire-obs watch --debounce 1s ./monitoring
  wetwire-obs watch -o prometheus.json ./monitoring`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			return runWatch(absPath, lintOnly, debounce, output)
		},
	}

	cmd.Flags().BoolVar(&lintOnly, "lint-only", false, "Only run lint, skip build")
	cmd.Flags().DurationVar(&debounce, "debounce", 500*time.Millisecond, "Debounce duration for rapid changes")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file for build")

	return cmd
}

func runWatch(path string, lintOnly bool, debounce time.Duration, output string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Add directory and subdirectories
	if err := addWatchDirs(watcher, path); err != nil {
		return fmt.Errorf("failed to add watch paths: %w", err)
	}

	fmt.Printf("Watching: %s\n", path)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initial build
	fmt.Println("Running initial lint/build...")
	runLintBuild(path, lintOnly, output)

	// Debounce timer
	var debounceTimer *time.Timer
	rebuildChan := make(chan struct{}, 1)

	fmt.Println("\nWatching for changes... (Ctrl+C to stop)")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			// Only watch .go files
			if !strings.HasSuffix(event.Name, ".go") {
				continue
			}

			// Skip generated files
			if strings.Contains(event.Name, "_gen.go") {
				continue
			}

			// Only process write/create events
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}

			// Debounce
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(debounce, func() {
				select {
				case rebuildChan <- struct{}{}:
				default:
				}
			})

		case <-rebuildChan:
			fmt.Printf("\n[%s] Change detected, rebuilding...\n", time.Now().Format("15:04:05"))
			runLintBuild(path, lintOnly, output)

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "Watch error: %v\n", err)

		case <-sigChan:
			fmt.Println("\nStopping watch...")
			return nil
		}
	}
}

func addWatchDirs(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip hidden and vendor directories
			name := filepath.Base(path)
			if strings.HasPrefix(name, ".") && path != root {
				return filepath.SkipDir
			}
			if name == "vendor" {
				return filepath.SkipDir
			}
			return watcher.Add(path)
		}
		return nil
	})
}

func runLintBuild(path string, lintOnly bool, output string) {
	d := &domain.ObservabilityDomain{}
	ctx := &coredomain.Context{}

	// Run lint
	lintResult, err := d.Linter().Lint(ctx, path, coredomain.LintOpts{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lint error: %v\n", err)
		return
	}

	if !lintResult.Success {
		fmt.Println("Lint failed:")
		for _, e := range lintResult.Errors {
			fmt.Printf("  %s: %s\n", e.Path, e.Message)
		}
		return
	}

	fmt.Println("Lint passed")

	if lintOnly {
		return
	}

	// Run build
	buildOpts := coredomain.BuildOpts{
		Format: "pretty",
		Output: output,
	}

	buildResult, err := d.Builder().Build(ctx, path, buildOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Build error: %v\n", err)
		return
	}

	if !buildResult.Success {
		fmt.Println("Build failed:")
		for _, e := range buildResult.Errors {
			fmt.Printf("  %s: %s\n", e.Path, e.Message)
		}
		return
	}

	if output != "" {
		fmt.Printf("Build successful, wrote %s\n", output)
	} else {
		fmt.Println("Build successful")
	}
}
