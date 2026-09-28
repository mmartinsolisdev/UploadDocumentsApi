package database

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "tenants.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

const legacyConfig = `{
  "gtmark": {
    "dsn": { "server": "HOST", "port": 1433, "database": "BD_GTMARK", "user": "U", "password": "P", "encrypt": false },
    "schema": { "saleType": false }
  }
}`

const v2Config = `{
  "sirenis": {
    "default": "OrigosVCSPT_Temp",
    "schema": { "saleType": true },
    "databases": [
      { "dsn": { "server": "HOST", "port": 1433, "database": "OrigosVCSPT_Temp", "user": "U", "password": "P" } },
      { "dsn": { "server": "HOST", "database": "OrigosVCSPT_Prod", "user": "U", "password": "P" } }
    ]
  }
}`

func TestLoadTenantsLegacy(t *testing.T) {
	if err := LoadTenants(writeConfig(t, legacyConfig)); err != nil {
		t.Fatalf("LoadTenants: %v", err)
	}

	dsn, err := ResolveDSN("gtmark", "")
	if err != nil {
		t.Fatalf("ResolveDSN default: %v", err)
	}
	if dsn.Database != "BD_GTMARK" {
		t.Fatalf("default database = %q, want BD_GTMARK", dsn.Database)
	}
	if _, err := ResolveDSN("gtmark", "BD_GTMARK"); err != nil {
		t.Fatalf("ResolveDSN explicit: %v", err)
	}
	if _, err := ResolveDSN("gtmark", "OTRA"); !errors.Is(err, ErrUnknownDatabase) {
		t.Fatalf("expected ErrUnknownDatabase, got %v", err)
	}
	if _, err := ResolveDSN("nope", ""); !errors.Is(err, ErrUnknownTenant) {
		t.Fatalf("expected ErrUnknownTenant, got %v", err)
	}
}

func TestLoadTenantsV2(t *testing.T) {
	if err := LoadTenants(writeConfig(t, v2Config)); err != nil {
		t.Fatalf("LoadTenants: %v", err)
	}

	dsn, err := ResolveDSN("sirenis", "")
	if err != nil {
		t.Fatalf("ResolveDSN default: %v", err)
	}
	if dsn.Database != "OrigosVCSPT_Temp" {
		t.Fatalf("default database = %q, want OrigosVCSPT_Temp", dsn.Database)
	}

	dsn, err = ResolveDSN("sirenis", "OrigosVCSPT_Prod")
	if err != nil {
		t.Fatalf("ResolveDSN second: %v", err)
	}
	if dsn.Server != "HOST" || dsn.Database != "OrigosVCSPT_Prod" {
		t.Fatalf("unexpected dsn: %+v", dsn)
	}

	if _, err := ResolveDSN("sirenis", "OTRA"); !errors.Is(err, ErrUnknownDatabase) {
		t.Fatalf("expected ErrUnknownDatabase, got %v", err)
	}
}

func TestLoadTenantsTrimsDefault(t *testing.T) {
	config := `{ "x": { "default": "  D  ", "databases": [ { "dsn": { "server": "H", "database": "D", "user": "U" } } ] } }`
	if err := LoadTenants(writeConfig(t, config)); err != nil {
		t.Fatalf("LoadTenants: %v", err)
	}
	dsn, err := ResolveDSN("x", "")
	if err != nil {
		t.Fatalf("ResolveDSN default: %v", err)
	}
	if dsn.Database != "D" {
		t.Fatalf("default database = %q, want D", dsn.Database)
	}
}

func TestLoadTenantsInvalid(t *testing.T) {
	cases := map[string]string{
		"empty catalog":        `{}`,
		"no dsn":               `{ "x": {} }`,
		"incomplete dsn":       `{ "x": { "dsn": { "server": "H", "user": "U" } } }`,
		"v2 without default":   `{ "x": { "databases": [ { "dsn": { "server": "H", "database": "D", "user": "U" } } ] } }`,
		"v2 default not found": `{ "x": { "default": "OTRA", "databases": [ { "dsn": { "server": "H", "database": "D", "user": "U" } } ] } }`,
		"v2 duplicate":         `{ "x": { "default": "D", "databases": [ { "dsn": { "server": "H", "database": "D", "user": "U" } }, { "dsn": { "server": "H2", "database": "D", "user": "U" } } ] } }`,
	}
	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			if err := LoadTenants(writeConfig(t, content)); err == nil {
				t.Fatalf("expected error for %s", name)
			}
		})
	}
}
