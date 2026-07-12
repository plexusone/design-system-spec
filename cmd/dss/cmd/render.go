package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/design-system-spec/sdk/go"
	"github.com/plexusone/structured-evaluation/rubric"
)

var (
	renderOutput   string
	renderTitle    string
	renderEvalFile string
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Generate HTML documentation for the design system",
	Long: `Generate a static HTML site documenting the design system spec.

The output includes:
  - Overview page with stats and principles
  - Component gallery with all components
  - Individual component detail pages
  - Design tokens visualization
  - Evaluation dashboard (if eval data provided)

Examples:
  dss render --output ./docs
  dss render -d ./specs/v3 --output ./docs/v3
  dss render --output ./docs --eval ./evals/v3.json
  dss render --output ./docs --title "Material Design 3"`,
	RunE: runRender,
}

func init() {
	renderCmd.Flags().StringVarP(&renderOutput, "output", "o", "docs", "output directory for HTML files")
	renderCmd.Flags().StringVar(&renderTitle, "title", "", "documentation title (default: design system name)")
	renderCmd.Flags().StringVar(&renderEvalFile, "eval", "", "evaluation JSON file to include")

	rootCmd.AddCommand(renderCmd)
}

func runRender(cmd *cobra.Command, args []string) error {
	dir := getSpecDir()
	ctx := context.Background()

	// Load design system
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	fmt.Printf("Rendering HTML documentation for: %s v%s\n", ds.Meta.Name, ds.Meta.Version)

	// Build options
	opts := &dss.HTMLOptions{
		OutputDir: renderOutput,
		Title:     renderTitle,
	}

	// Load evaluation data if provided
	if renderEvalFile != "" {
		evalData, err := os.ReadFile(renderEvalFile)
		if err != nil {
			return fmt.Errorf("reading eval file: %w", err)
		}

		var evalResult rubric.Rubric
		if err := json.Unmarshal(evalData, &evalResult); err != nil {
			return fmt.Errorf("parsing eval JSON: %w", err)
		}

		opts.EvalResult = &evalResult
		fmt.Printf("Including evaluation data: score %d/5 (%s)\n", evalResult.IntScore, evalResult.IntScore.String())
	} else {
		// Generate fresh evaluation
		service := dss.NewService(ds)
		evalResult, err := service.Evaluate(ctx, nil)
		if err != nil {
			fmt.Printf("Warning: could not generate evaluation: %v\n", err)
		} else {
			opts.EvalResult = evalResult
		}
	}

	// Generate HTML
	if err := ds.GenerateHTML(opts); err != nil {
		return fmt.Errorf("generating HTML: %w", err)
	}

	fmt.Printf("\nGenerated HTML documentation in: %s\n", renderOutput)
	fmt.Printf("Open %s/index.html in a browser to view.\n", renderOutput)

	return nil
}
