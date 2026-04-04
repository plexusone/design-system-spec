package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/design-system-spec/sdk/go"
)

var (
	cssOutput   string
	typesOutput string
	llmOutput   string
	cssFormat   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code artifacts from design system spec",
	Long: `Generate CSS, TypeScript types, and LLM prompts from a design system specification.

By default, outputs to stdout. Use flags to write to files.

Examples:
  # Preview all outputs
  dss generate

  # Generate specific files
  dss generate --css ./src/index.css --types ./src/lib/types.ts --llm ./DESIGN_CONTEXT.md

  # Generate from a different directory
  dss generate -d ./design-system --css ./web/src/index.css`,
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVar(&cssOutput, "css", "", "output path for CSS (default: stdout)")
	generateCmd.Flags().StringVar(&typesOutput, "types", "", "output path for TypeScript types (default: stdout)")
	generateCmd.Flags().StringVar(&llmOutput, "llm", "", "output path for LLM prompt (default: stdout)")
	generateCmd.Flags().StringVar(&cssFormat, "css-format", "tailwind4", "CSS format: tailwind4, css-vars, scss, mkdocs-material")

	rootCmd.AddCommand(generateCmd)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	dir := getSpecDir()

	// Load the design system
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	// Validate
	if err := ds.Validate(); err != nil {
		return fmt.Errorf("design system validation failed: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Loaded design system: %s v%s\n", ds.Meta.Name, ds.Meta.Version)

	// Generate CSS
	cssOpts := dss.DefaultCSSOptions()
	switch cssFormat {
	case "tailwind4":
		cssOpts.Format = "tailwind4"
	case "vars", "css-vars":
		cssOpts.Format = "css-vars"
	case "scss":
		cssOpts.Format = "scss"
	case "mkdocs-material", "mkdocs":
		cssOpts.Format = "mkdocs-material"
	default:
		return fmt.Errorf("unknown CSS format: %s (valid: tailwind4, css-vars, scss, mkdocs-material)", cssFormat)
	}

	css, err := ds.GenerateCSS(cssOpts)
	if err != nil {
		return fmt.Errorf("generating CSS: %w", err)
	}

	if cssOutput != "" {
		if err := os.WriteFile(cssOutput, []byte(css), 0644); err != nil {
			return fmt.Errorf("writing CSS: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Generated CSS: %s\n", cssOutput)
	} else if typesOutput == "" && llmOutput == "" {
		// Only print to stdout if no files specified
		fmt.Println("=== CSS ===")
		fmt.Println(css)
	}

	// Generate LLM prompt
	llmOpts := dss.DefaultLLMPromptOptions()
	llmPrompt, err := ds.GenerateLLMPrompt(llmOpts)
	if err != nil {
		return fmt.Errorf("generating LLM prompt: %w", err)
	}

	if llmOutput != "" {
		if err := os.WriteFile(llmOutput, []byte(llmPrompt), 0644); err != nil {
			return fmt.Errorf("writing LLM prompt: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Generated LLM prompt: %s\n", llmOutput)
	} else if cssOutput == "" && typesOutput == "" {
		fmt.Println("\n=== LLM Prompt ===")
		fmt.Println(llmPrompt)
	}

	// Generate TypeScript types
	if len(ds.Components) > 0 {
		reactOpts := dss.DefaultReactOptions()
		types, err := ds.GenerateReactTypes(reactOpts)
		if err != nil {
			return fmt.Errorf("generating TypeScript types: %w", err)
		}

		if typesOutput != "" {
			if err := os.WriteFile(typesOutput, []byte(types), 0644); err != nil {
				return fmt.Errorf("writing TypeScript types: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Generated TypeScript types: %s\n", typesOutput)
		} else if cssOutput == "" && llmOutput == "" {
			fmt.Println("\n=== TypeScript Types ===")
			fmt.Println(types)
		}
	}

	fmt.Fprintf(os.Stderr, "Done!\n")
	return nil
}
