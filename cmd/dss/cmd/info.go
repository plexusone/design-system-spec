package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	dss "github.com/plexusone/design-system-spec/sdk/go"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display information about the design system",
	Long: `Display information about the design system specification.

Shows metadata, foundation counts, component counts, and validation status.

Examples:
  dss info
  dss info -d ./design-system`,
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	dir := getSpecDir()

	// Load design system
	ds, err := dss.LoadDesignSystem(dir)
	if err != nil {
		return fmt.Errorf("loading design system: %w", err)
	}

	fmt.Printf("Design System: %s\n", ds.Meta.Name)
	fmt.Printf("Version: %s\n", ds.Meta.Version)
	if ds.Meta.Description != "" {
		fmt.Printf("Description: %s\n", ds.Meta.Description)
	}
	fmt.Println()

	// Principles
	if len(ds.Principles) > 0 {
		fmt.Printf("Principles: %d\n", len(ds.Principles))
		for _, p := range ds.Principles {
			fmt.Printf("  - %s\n", p.Name)
		}
		fmt.Println()
	}

	// Foundations
	fmt.Println("Foundations:")
	if len(ds.Foundations.Colors) > 0 {
		fmt.Printf("  Colors: %d tokens\n", len(ds.Foundations.Colors))
	}
	if len(ds.Foundations.Typography.FontFamilies) > 0 {
		fmt.Printf("  Font Families: %d\n", len(ds.Foundations.Typography.FontFamilies))
	}
	if len(ds.Foundations.Typography.FontSizes) > 0 {
		fmt.Printf("  Font Sizes: %d\n", len(ds.Foundations.Typography.FontSizes))
	}
	if len(ds.Foundations.Spacing.Scale) > 0 {
		fmt.Printf("  Spacing: %d values\n", len(ds.Foundations.Spacing.Scale))
	}
	if len(ds.Foundations.BorderRadius) > 0 {
		fmt.Printf("  Border Radius: %d values\n", len(ds.Foundations.BorderRadius))
	}
	fmt.Println()

	// Components
	if len(ds.Components) > 0 {
		fmt.Printf("Components: %d\n", len(ds.Components))
		for _, c := range ds.Components {
			variantCount := len(c.Variants)
			if variantCount > 0 {
				fmt.Printf("  - %s (%d variants)\n", c.Name, variantCount)
			} else {
				fmt.Printf("  - %s\n", c.Name)
			}
		}
		fmt.Println()
	}

	// Accessibility
	if ds.Accessibility != nil && ds.Accessibility.WCAGLevel != "" {
		fmt.Printf("Accessibility Target: WCAG %s\n", ds.Accessibility.WCAGLevel)
	}

	// Validation
	if err := ds.Validate(); err != nil {
		fmt.Printf("\n⚠ Validation Issues: %v\n", err)
	} else {
		fmt.Println("✓ Spec validation passed")
	}

	return nil
}
