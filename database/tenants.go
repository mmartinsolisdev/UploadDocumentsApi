package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// TenantDSN contiene los datos de conexión de una base.
type TenantDSN struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	Encrypt  bool   `json:"encrypt"`
}

// TenantDatabase es una base de datos dentro de la entrada de un cliente (formato v2).
type TenantDatabase struct {
	DSN TenantDSN `json:"dsn"`
}

// rawTenant es la entrada tal cual viene en tenants.json. Admite el formato legacy
// ({ "dsn": { ... } }) y el v2 ({ "default": "...", "databases": [ ... ] }).
type rawTenant struct {
	DSN       TenantDSN        `json:"dsn"`
	Default   string           `json:"default"`
	Databases []TenantDatabase `json:"databases"`
}

// tenantCatalog es la entrada ya normalizada: BD por defecto + DSNs por nombre físico.
type tenantCatalog struct {
	defaultDatabase string
	databases       map[string]TenantDSN
}

var (
	// ErrUnknownTenant indica que el slug no existe en el catálogo.
	ErrUnknownTenant = errors.New("unknown tenant")
	// ErrUnknownDatabase indica que el nombre físico no existe para ese cliente.
	ErrUnknownDatabase = errors.New("unknown database")

	tenants map[string]tenantCatalog
	pools   sync.Map // "<slug>:<database>" -> *gorm.DB
)

// LoadTenants carga el catálogo de clientes y bases de datos desde un archivo JSON.
// Formatos aceptados por cliente:
//   - legacy: { "dsn": { server, port, database, user, password, encrypt } }
//   - v2:     { "default": "<database>", "databases": [ { "dsn": { ... } }, ... ] }
//
// En ambos casos la identidad de una BD es su campo `dsn.database` (nombre físico).
func LoadTenants(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read tenants config: %w", err)
	}
	var raw map[string]rawTenant
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse tenants config: %w", err)
	}
	if len(raw) == 0 {
		return fmt.Errorf("tenants config is empty")
	}
	catalog := make(map[string]tenantCatalog, len(raw))
	for slug, entry := range raw {
		resolved, err := normalizeTenant(slug, entry)
		if err != nil {
			return err
		}
		catalog[slug] = resolved
	}
	tenants = catalog
	return nil
}

// normalizeTenant valida y convierte una entrada cruda al catálogo interno.
func normalizeTenant(slug string, entry rawTenant) (tenantCatalog, error) {
	if len(entry.Databases) > 0 {
		databases := make(map[string]TenantDSN, len(entry.Databases))
		for i, database := range entry.Databases {
			if err := validateDSN(database.DSN); err != nil {
				return tenantCatalog{}, fmt.Errorf("tenant %q, databases[%d]: %w", slug, i, err)
			}
			name := database.DSN.Database
			if _, exists := databases[name]; exists {
				return tenantCatalog{}, fmt.Errorf("tenant %q: duplicate database %q", slug, name)
			}
			databases[name] = database.DSN
		}
		entry.Default = strings.TrimSpace(entry.Default)
		if entry.Default == "" {
			return tenantCatalog{}, fmt.Errorf("tenant %q: \"default\" is required when \"databases\" is used", slug)
		}
		if _, ok := databases[entry.Default]; !ok {
			return tenantCatalog{}, fmt.Errorf("tenant %q: default database %q is not in \"databases\"", slug, entry.Default)
		}
		return tenantCatalog{defaultDatabase: entry.Default, databases: databases}, nil
	}

	if err := validateDSN(entry.DSN); err != nil {
		return tenantCatalog{}, fmt.Errorf("tenant %q: %w", slug, err)
	}
	name := entry.DSN.Database
	return tenantCatalog{
		defaultDatabase: name,
		databases:       map[string]TenantDSN{name: entry.DSN},
	}, nil
}

func validateDSN(dsn TenantDSN) error {
	if dsn.Server == "" || dsn.Database == "" || dsn.User == "" {
		return fmt.Errorf("incomplete dsn (server, database and user are required)")
	}
	return nil
}

// ResolveDSN devuelve el DSN de la BD indicada; con `database` vacío usa la del cliente.
func ResolveDSN(slug, database string) (TenantDSN, error) {
	tenant, ok := tenants[slug]
	if !ok {
		return TenantDSN{}, fmt.Errorf("%w: %q", ErrUnknownTenant, slug)
	}
	if database == "" {
		database = tenant.defaultDatabase
	}
	dsn, ok := tenant.databases[database]
	if !ok {
		return TenantDSN{}, fmt.Errorf("%w: %q (tenant %q)", ErrUnknownDatabase, database, slug)
	}
	return dsn, nil
}

// GetDB devuelve (y cachea) la conexión GORM del cliente y su base de datos.
// `database` vacío ⇒ base por defecto del cliente.
func GetDB(slug, database string) (*gorm.DB, error) {
	dsn, err := ResolveDSN(slug, database)
	if err != nil {
		return nil, err
	}
	key := slug + ":" + dsn.Database
	if cached, ok := pools.Load(key); ok {
		return cached.(*gorm.DB), nil
	}
	encrypt := "disable"
	if dsn.Encrypt {
		encrypt = "true"
	}
	port := dsn.Port
	if port == 0 {
		port = 1433
	}
	conn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&connection+timeout=30&encrypt=%s",
		dsn.User, dsn.Password, dsn.Server, port, dsn.Database, encrypt)
	db, err := gorm.Open(sqlserver.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tenant %q database %q: %w", slug, dsn.Database, err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(1)
	pools.Store(key, db)
	return db, nil
}
