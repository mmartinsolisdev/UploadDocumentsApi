# Migración de Fiber v2 → v3

## Objetivo

Migrar `UploadDocumentsAPI` de `github.com/gofiber/fiber/v2` (v2.52.15) a `github.com/gofiber/fiber/v3` (última), sin cambiar el comportamiento de la aplicación.

## Estado actual (verificado)

- `go.mod` declara `go 1.27.0` (v3 exige Go 1.25+, cumple).
- `go.mod` línea 7: `github.com/gofiber/fiber/v2 v2.52.15` — **aún en v2**, no hay ninguna línea de `/v3`.
- Archivos que importan Fiber:
  - `main.go` (fiber, `middleware/recover`, `middleware/cors`)
  - `routes/routes.go` (fiber)
  - `controllers/uploader/uploader.go` (fiber)
  - `middleware/firebase.go` (fiber)
- Hay cambios sin commitear en el árbol de trabajo que **no son parte de la migración** y deben preservarse:
  - `main.go`: lógica de `allowedOrigins` según `APP_ENV` (env == "development").
  - `controllers/uploader/uploader.go`: refactor de `combosData` en slices tipadas (`ids`, `docNames`, `languages`, `docTypes`, `saleTypes`).
  - `go.mod`/`go.sum`: solo cambios de fin de línea (CRLF/LF), sin cambios de contenido.

## Cambios de API relevantes (según guía oficial v3)

1. **Import path**: `/v2` → `/v3` en todos los imports.
2. **CORS**: `AllowOrigins`, `AllowMethods`, `AllowHeaders` y `ExposeHeaders` cambian de `string` (separado por comas) a `[]string`. `MaxAge` y `AllowCredentials` no cambian.
3. **Listen**: `app.Listen("host:port")` con string sigue siendo válido. No requiere cambio.
4. **Recover**: `recover.New()` sin cambios.
5. **Ctx**: `c.Query`, `c.Get`, `c.JSON`, `c.SendStatus`, `c.Status`, `c.Locals`, `c.FormFile`, `c.Next`, `fiber.StatusUnauthorized`, `fiber.Map` — sin cambios.
6. **Router**: `app.Group`, `app.Use` — sin cambios.

## Tareas ordenadas

### 1. Actualizar imports `/v2` → `/v3`

- `main.go`:
  - `"github.com/gofiber/fiber/v2"` → `"github.com/gofiber/fiber/v3"`
  - `"github.com/gofiber/fiber/v2/middleware/recover"` → `"github.com/gofiber/fiber/v3/middleware/recover"`
  - `"github.com/gofiber/fiber/v2/middleware/cors"` → `"github.com/gofiber/fiber/v3/middleware/cors"`
- `routes/routes.go`: `"github.com/gofiber/fiber/v2"` → `"github.com/gofiber/fiber/v3"`
- `controllers/uploader/uploader.go`: `"github.com/gofiber/fiber/v2"` → `"github.com/gofiber/fiber/v3"`
- `middleware/firebase.go`: `"github.com/gofiber/fiber/v2"` → `"github.com/gofiber/fiber/v3"`

### 2. Convertir CORS a slices en `main.go`

Reemplazar el bloque de configuración CORS, preservando la lógica de `allowedOrigins` ya existente (sin commitear). Resultado objetivo:

```go
allowedOrigins := []string{"http://origos.no-ip.com", "https://origos.no-ip.com"}
if env == "development" {
    allowedOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
}

app.Use(cors.New(cors.Config{
    Next:             nil,
    AllowOrigins:     allowedOrigins,
    AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Length", "Accept", "Content-Type", "Accept-Encoding", "Accept-Language", "Authorization"},
    AllowCredentials: false,
    ExposeHeaders:    []string{},
    MaxAge:           0,
}))
```

Notas:
- El valor actual de `AllowMethods` contiene un espacio extra en `"PATCH, OPTIONS"`; al convertirlo a slice, separar en tokens limpios.
- `ExposeHeaders: ""` pasa a `[]string{}` (o `nil`).

### 3. Actualizar dependencias

```bash
go get github.com/gofiber/fiber/v3@latest
go mod tidy
```

- `go get` fija la versión de `/v3` en `go.mod` (ahora sí persistirá porque el código lo importa).
- `go mod tidy` elimina `fiber/v2` de `go.mod` y `go.sum` al no quedar imports de `/v2`, y resuelve dependencias transitivas nuevas de v3.

### 4. Verificación de compilación

```bash
go build ./...
go vet ./...
```

### 5. Verificación de dependencias

```bash
go list -m github.com/gofiber/fiber/v3
```

Debe mostrar la versión instalada (ej. `github.com/gofiber/fiber/v3 v3.x.y`), y `go.mod` no debe contener `fiber/v2`.

### 6. Prueba de ejecución

```bash
go run .
```

- Confirmar que el servidor arranca en `127.0.0.1:PORT_APP`.
- Probar un endpoint protegido (`GET /uploader/MembershipsList`) con token válido.
- Verificar CORS: una petición `OPTIONS` de preflight desde un origen permitido (localhost:3000 en dev) debe devolver los headers `Access-Control-Allow-*` correctos, funcionalmente equivalentes a v2.

## Riesgos y consideraciones

- **CORS funcionalmente equivalente**: al pasar a slices, validar que los orígenes/métodos/headers permitidos sean idénticos a los valores actuales (sin perder ni duplicar entradas). El espacio en `"PATCH, OPTIONS"` sugiere un bug previo que conviene corregir, no perpetuar.
- **Versión de Go**: v3 exige 1.25+; `go.mod` usa `1.27.0`, sin problema.
- **Cambios sin commitear**: `main.go` y `uploader.go` tienen ediciones previas no relacionadas. Editar solo las líneas indicadas y no revertir el resto.
- **`go.sum`**: se modificará al resolver v3; no es un error, es esperado.

## Fuera de alcance

- Migración de `ioutil.ReadAll` a `io.ReadAll` (deprecación de stdlib, no relacionada con Fiber).
- Cambios de lógica de negocio (uploader, firebase, modelos).

## Validación final

1. `go build ./...` y `go vet ./...` sin errores.
2. `go mod tidy` sin cambios residuales (ejecutarlo dos veces debe ser idempotente).
3. `go.mod` sin `fiber/v2`, con `fiber/v3` en el bloque `require`.
4. Servidor arranca y endpoints responden igual que con v2 (especialmente preflight CORS).
