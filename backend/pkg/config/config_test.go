package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateConfigSkipsJWTWhenFeaturesEnableAuthFalse(t *testing.T) {
	cfg := &Config{}
	cfg.App.Environment = "production"
	cfg.Auth.Enabled = true
	cfg.Auth.JWT.SecretKey = ""
	cfg.Features.EnableAuth = false

	if err := validateConfig(cfg); err != nil {
		t.Fatalf("validateConfig should allow missing JWT when features.enable_auth is false: %v", err)
	}
}

func TestValidateConfigRequiresJWTWhenBothAuthSwitchesEnabled(t *testing.T) {
	cfg := &Config{}
	cfg.App.Environment = "production"
	cfg.Auth.Enabled = true
	cfg.Auth.JWT.SecretKey = ""
	cfg.Features.EnableAuth = true

	if err := validateConfig(cfg); err == nil {
		t.Fatal("expected validateConfig to require jwt.secretKey when both auth switches are enabled")
	}
}

func TestLoadBootstrapAdminFromPrefixedEnv(t *testing.T) {
	t.Setenv("APP_AUTH_BOOTSTRAPADMIN_EMAIL", "bootstrap@example.com")
	t.Setenv("APP_AUTH_BOOTSTRAPADMIN_PASSWORD", "bootstrap-secret")
	t.Setenv("APP_AUTH_BOOTSTRAPADMIN_FIRSTNAME", "Bootstrap")
	t.Setenv("APP_AUTH_BOOTSTRAPADMIN_LASTNAME", "Admin")
	t.Setenv("APP_APP_ENVIRONMENT", "development")

	// Avoid picking up a repo configs/app.yaml that could shadow env values.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if got := cfg.Auth.BootstrapAdmin.Email; got != "bootstrap@example.com" {
		t.Fatalf("BootstrapAdmin.Email = %q, want bootstrap@example.com", got)
	}
	if got := cfg.Auth.BootstrapAdmin.Password; got != "bootstrap-secret" {
		t.Fatalf("BootstrapAdmin.Password = %q, want bootstrap-secret", got)
	}
	if got := cfg.Auth.BootstrapAdmin.FirstName; got != "Bootstrap" {
		t.Fatalf("BootstrapAdmin.FirstName = %q, want Bootstrap", got)
	}
	if got := cfg.Auth.BootstrapAdmin.LastName; got != "Admin" {
		t.Fatalf("BootstrapAdmin.LastName = %q, want Admin", got)
	}

	// Ensure we did not accidentally create a config file in the temp dir.
	if _, err := os.Stat(filepath.Join(tmp, "configs", "app.yaml")); !os.IsNotExist(err) {
		t.Fatalf("unexpected configs/app.yaml presence: %v", err)
	}
}
