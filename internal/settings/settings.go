package settings

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type Settings struct {
	General struct {
		Organization string `koanf:"organization"`
		SecretsKeyID string `koanf:"secrets_key_id"`
	} `koanf:"general"`

	Sites struct {
		Ordered []string `koanf:"ordered"`
	} `koanf:"sites"`

	Environments struct {
		Default struct {
			Domain      string `koanf:"domain"`
			Name        string `koanf:"name"`
			ProfileName string `koanf:"profile_name"`
			Region      string `koanf:"region"`
		} `koanf:"default"`

		Permanent []string `koanf:"permanent"`

		Temporary struct {
			Secrets struct {
				FirebaseCredentials string `koanf:"firebase_credentials"`
				MailPassword        string `koanf:"mail_password"`
			} `koanf:"secrets"`
		} `koanf:"temporary"`
	} `koanf:"environments"`
	DefaultEnvironment string `koanf:"default_environment"`
}

// LoadSettings loads settings from a specified file path using Koanf
func LoadSettings(path string) (*Settings, error) {
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load settings file: %w", err)
	}

	var settings Settings
	if err := k.Unmarshal("", &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}
	return &settings, nil
}

// ValidateSettings checks whether all required fields are present
func ValidateSettings(settings *Settings) error {
	if settings.General.Organization == "" {
		return fmt.Errorf("organization field in general config is required")
	}
	if len(settings.Sites.Ordered) == 0 {
		return fmt.Errorf("at least one site must be listed in sites.ordered")
	}
	if settings.Environments.Default.Name == "" {
		return fmt.Errorf("default environment name is required")
	}
	return nil
}

// GetSettings loads the settings from TERRABUTLER_ROOT/configs/settings.yml
func GetSettings() *Settings {
	root := os.Getenv("TERRABUTLER_ROOT")
	return GetSettingsFromRoot(root)
}

// GetSettingsFromRoot allows tests or external callers to specify the root directory
func GetSettingsFromRoot(root string) *Settings {
	settingsPath := fmt.Sprintf("%s/internal/configs/settings.yml", root)

	settings, err := LoadSettings(settingsPath)
	if err != nil {
		panic(fmt.Errorf("error loading settings: %w", err))
	}

	if err := ValidateSettings(settings); err != nil {
		panic(fmt.Errorf("invalid settings: %w", err))
	}

	return settings
}
