package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateConfig_JWTSecretFailClosed(t *testing.T) {
	base := func(secret string, authEnabled bool, env string) *Config {
		return &Config{
			App: AppConfig{
				Name:        "cc-starship",
				Environment: env,
			},
			Auth: AuthConfig{
				Enabled: authEnabled,
				JWT: JWTConfig{
					SecretKey: secret,
					Issuer:    "cc-starship-api",
				},
			},
			Databases: DatabasesConfig{Default: ""},
		}
	}

	tests := []struct {
		name        string
		cfg         *Config
		wantErr     bool
		errContains string
	}{
		{
			name:    "auth enabled with secure secret succeeds",
			cfg:     base("unit-test-secure-jwt-secret-value", true, "development"),
			wantErr: false,
		},
		{
			name:        "auth enabled with empty secret fails in development",
			cfg:         base("", true, "development"),
			wantErr:     true,
			errContains: "jwt.secretKey is required",
		},
		{
			name:        "auth enabled with empty secret fails in production",
			cfg:         base("", true, "production"),
			wantErr:     true,
			errContains: "jwt.secretKey is required",
		},
		{
			name:        "auth enabled with committed placeholder fails",
			cfg:         base("your-secret-key-here", true, "development"),
			wantErr:     true,
			errContains: "placeholder",
		},
		{
			name:        "auth enabled with former dev fallback fails",
			cfg:         base("dev-secret-key", true, "development"),
			wantErr:     true,
			errContains: "placeholder",
		},
		{
			name:    "auth disabled allows empty secret",
			cfg:     base("", false, "development"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestLoad_RequiresNonPlaceholderSecretFromEnv(t *testing.T) {
	tmp := t.TempDir()
	configDir := filepath.Join(tmp, "configs")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	configPath := filepath.Join(configDir, "app.yaml")
	yaml := []byte(`
app:
  name: "cc-starship"
  environment: "development"
auth:
  enabled: true
  jwt:
    secretKey: ""
    issuer: "cc-starship-api"
    expiryDuration: "24h"
databases:
  default: ""
`)
	if err := os.WriteFile(configPath, yaml, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	origWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWD)
	})

	t.Run("refuses empty and placeholder without env", func(t *testing.T) {
		t.Setenv("APP_AUTH_JWT_SECRETKEY", "")
		_, err := Load()
		if err == nil {
			t.Fatal("expected Load to fail without JWT secret")
		}
		if !strings.Contains(err.Error(), "jwt.secretKey") {
			t.Fatalf("unexpected error: %v", err)
		}

		t.Setenv("APP_AUTH_JWT_SECRETKEY", "your-secret-key-here")
		_, err = Load()
		if err == nil {
			t.Fatal("expected Load to fail with placeholder JWT secret")
		}
		if !strings.Contains(err.Error(), "placeholder") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("starts when env supplies non-placeholder secret", func(t *testing.T) {
		t.Setenv("APP_AUTH_JWT_SECRETKEY", "unit-test-secure-jwt-secret-from-env")
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if cfg.Auth.JWT.SecretKey != "unit-test-secure-jwt-secret-from-env" {
			t.Fatalf("secret not loaded from env, got %q", cfg.Auth.JWT.SecretKey)
		}
	})
}
