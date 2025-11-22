package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type AppSettings struct {
    CliqHubBaseURL string `mapstructure:"cliq_hub_base_url"`
}

type SettingsService struct {
	vp         *viper.Viper
	configFile string
}

func NewSettingsService() (*SettingsService, error) {
	vp := viper.New()
	vp.SetEnvPrefix("CLIQ")
	vp.AutomaticEnv()

    vp.SetDefault("cliq_hub_base_url", "http://localhost:8080")

	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("get user config dir: %w", err)
	}
	appDir := filepath.Join(cfgDir, "cliq")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return nil, fmt.Errorf("ensure app config dir: %w", err)
	}
	cfgFile := filepath.Join(appDir, "settings.yaml")
	vp.SetConfigFile(cfgFile)
	vp.SetConfigType("yaml")

	if _, err := os.Stat(cfgFile); err == nil {
		if err := vp.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
	} else {
		if err := vp.WriteConfigAs(cfgFile); err != nil {
			return nil, fmt.Errorf("write default config: %w", err)
		}
	}

	return &SettingsService{vp: vp, configFile: cfgFile}, nil
}

func (s *SettingsService) Load() (*AppSettings, error) {
    var cfg AppSettings
    if err := s.vp.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("unmarshal settings: %w", err)
    }
    return &cfg, nil
}

func (s *SettingsService) Save(in *AppSettings) error {
    if in == nil {
        return errors.New("settings is nil")
    }
    if err := validateURL(in.CliqHubBaseURL); err != nil {
        return err
    }
    s.vp.Set("cliq_hub_base_url", in.CliqHubBaseURL)
    return s.vp.WriteConfigAs(s.configFile)
}

func (s *SettingsService) Update(partial map[string]any) error {
	if partial == nil {
		return nil
	}
	// Validate known keys
    if v, ok := partial["cliq_hub_base_url"]; ok {
        if str, ok2 := v.(string); ok2 {
            if err := validateURL(str); err != nil {
                return err
            }
            s.vp.Set("cliq_hub_base_url", str)
        }
    }
    return s.vp.WriteConfigAs(s.configFile)
}

func validateURL(u string) error {
	if u == "" {
		return errors.New("hub_base_url cannot be empty")
	}
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("hub_base_url must use http or https")
	}
	return nil
}
