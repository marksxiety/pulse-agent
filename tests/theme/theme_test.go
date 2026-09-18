package theme_test

import (
	"reflect"
	"regexp"
	"testing"

	"pulse-agent/theme"
)

var hexColour = regexp.MustCompile(`^#[0-9a-f]{6}$`)

// The picker cursors through Names by index, so this order is user-visible and
// changing it silently remaps everyone's saved theme to a different palette.
func TestPresets_OrderIsStable(t *testing.T) {
	want := []string{
		"original",
		"catppuccin-mocha",
		"dracula",
		"nord",
		"gruvbox",
		"tokyo-night",
		"one-dark-pro",
		"github-dark",
		"ayu-mirage",
		"monokai-pro",
		"synthwave",
	}

	if len(theme.Presets) != len(want) {
		t.Fatalf("got %d presets, want %d", len(theme.Presets), len(want))
	}
	for i, name := range want {
		if theme.Presets[i].Name != name {
			t.Errorf("Presets[%d].Name = %q, want %q", i, theme.Presets[i].Name, name)
		}
	}
}

func TestNames_MatchPresets(t *testing.T) {
	if len(theme.Names) != len(theme.Presets) {
		t.Fatalf("Names has %d entries, Presets has %d", len(theme.Names), len(theme.Presets))
	}
	for i, preset := range theme.Presets {
		if theme.Names[i] != preset.Name {
			t.Errorf("Names[%d] = %q, want %q", i, theme.Names[i], preset.Name)
		}
	}
}

func TestPresets_NamesAreUnique(t *testing.T) {
	seen := make(map[string]int, len(theme.Presets))
	for i, preset := range theme.Presets {
		if previous, ok := seen[preset.Name]; ok {
			t.Errorf("duplicate name %q at indexes %d and %d", preset.Name, previous, i)
		}
		seen[preset.Name] = i
	}
}

// Every colour field is emitted to the webview as a CSS variable, so a blank
// value would silently produce an unstyled widget rather than a loud failure.
// Walking the struct by reflection means a newly added colour field is covered
// automatically instead of quietly escaping the guard.
func TestPresets_HaveCompletePalettes(t *testing.T) {
	for _, preset := range theme.Presets {
		value := reflect.ValueOf(preset)

		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			if field.Name == "Name" {
				continue
			}
			if field.Type.Kind() != reflect.String {
				t.Errorf("%s.%s is %s, want a colour string", preset.Name, field.Name, field.Type)
				continue
			}

			colour := value.Field(i).String()
			if !hexColour.MatchString(colour) {
				t.Errorf("%s.%s = %q, want a #rrggbb colour", preset.Name, field.Name, colour)
			}
		}
	}
}

func TestGet_ReturnsMatchingPreset(t *testing.T) {
	for _, preset := range theme.Presets {
		if got := theme.Get(preset.Name); got != preset {
			t.Errorf("Get(%q) = %+v, want %+v", preset.Name, got, preset)
		}
	}
}

func TestGet_UnknownNameFallsBackToOriginal(t *testing.T) {
	if got := theme.Get("does-not-exist"); got != theme.Original {
		t.Errorf("Get() = %+v, want the original theme", got)
	}
	if got := theme.Get(""); got != theme.Original {
		t.Errorf("Get(\"\") = %+v, want the original theme", got)
	}
}
