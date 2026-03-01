package dss

import (
	"encoding/json"
	"testing"
)

func TestDesignSystemValidate(t *testing.T) {
	tests := []struct {
		name    string
		ds      DesignSystem
		wantErr bool
	}{
		{
			name: "valid minimal system",
			ds: DesignSystem{
				Meta: Meta{
					Name:    "Test System",
					Version: "1.0.0",
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			ds: DesignSystem{
				Meta: Meta{
					Version: "1.0.0",
				},
			},
			wantErr: true,
		},
		{
			name: "missing version",
			ds: DesignSystem{
				Meta: Meta{
					Name: "Test System",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ds.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDesignSystemJSON(t *testing.T) {
	ds := DesignSystem{
		Meta: Meta{
			Name:          "Test System",
			Version:       "1.0.0",
			Description:   "A test design system",
			MaturityLevel: 3,
		},
		Foundations: Foundations{
			Colors: []ColorToken{
				{
					ID:       "primary-500",
					Value:    "#0066CC",
					Semantic: "primary",
					Usage:    "Primary brand color for main CTAs",
				},
			},
		},
		Components: []Component{
			{
				ID:          "button",
				Name:        "Button",
				Description: "Interactive button component",
				Variants: []Variant{
					{ID: "primary", Name: "Primary", IsDefault: true},
					{ID: "secondary", Name: "Secondary"},
				},
			},
		},
	}

	data, err := json.MarshalIndent(ds, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var parsed DesignSystem
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if parsed.Meta.Name != ds.Meta.Name {
		t.Errorf("Name mismatch: got %q, want %q", parsed.Meta.Name, ds.Meta.Name)
	}
	if len(parsed.Foundations.Colors) != 1 {
		t.Errorf("Colors count mismatch: got %d, want 1", len(parsed.Foundations.Colors))
	}
	if len(parsed.Components) != 1 {
		t.Errorf("Components count mismatch: got %d, want 1", len(parsed.Components))
	}
}

func TestGenerateSchema(t *testing.T) {
	schema := GenerateSchema()
	if schema == nil {
		t.Fatal("GenerateSchema returned nil")
	}

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	if len(data) < 100 {
		t.Errorf("Schema seems too small: %d bytes", len(data))
	}
}
