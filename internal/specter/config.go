package specter

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Keybindings struct {
	Up                rune
	Down              rune
	ExpandCollapse    rune
	ExpandCollapseAll rune
	Quit              rune
}

type Config struct {
	Keybindings Keybindings
}

// keybindingsFile mirrors Config for YAML unmarshalling.
type keybindingsFile struct {
	Up                string `yaml:"up"`
	Down              string `yaml:"down"`
	ExpandCollapse    string `yaml:"expand_collapse"`
	ExpandCollapseAll string `yaml:"expand_collapse_all"`
	Quit              string `yaml:"quit"`
}

type configFile struct {
	Keybindings keybindingsFile `yaml:"keybindings"`
}

func DefaultConfig() *Config {
	return &Config{
		Keybindings: Keybindings{
			Up:                'k',
			Down:              'j',
			ExpandCollapse:    ' ',
			ExpandCollapseAll: 'e',
			Quit:              'q',
		},
	}
}

func DefaultConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "specter", "config.yaml"), nil
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var f configFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, err
	}

	d := DefaultConfig()
	return &Config{
		Keybindings: Keybindings{
			Up:                runeOrDefault(f.Keybindings.Up, d.Keybindings.Up),
			Down:              runeOrDefault(f.Keybindings.Down, d.Keybindings.Down),
			ExpandCollapse:    runeOrDefault(f.Keybindings.ExpandCollapse, d.Keybindings.ExpandCollapse),
			ExpandCollapseAll: runeOrDefault(f.Keybindings.ExpandCollapseAll, d.Keybindings.ExpandCollapseAll),
			Quit:              runeOrDefault(f.Keybindings.Quit, d.Keybindings.Quit),
		},
	}, nil
}

func runeOrDefault(s string, fallback rune) rune {
	for _, r := range s {
		return r
	}
	return fallback
}
