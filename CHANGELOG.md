# Historial de cambios

Registro de cambios y decisiones técnicas del proyecto para futuras consultas.

## [2026-08-26] - Corrección CORS para entorno local

### Problema
Al consumir la API desde la app web local (`http://localhost:3000`) el navegador bloqueaba las peticiones con error CORS en el preflight:

```
Access to fetch at 'http://127.0.0.1:3001/uploader/CombosList' from origin
'http://localhost:3000' has been blocked by CORS policy: No 'Access-Control-Allow-Origin'
header is present on the requested resource.
```

### Causa
En `main.go`, la configuración de CORS solo permitía los orígenes de producción:

```go
AllowOrigins: "http://origos.no-ip.com, https://origos.no-ip.com",
```

El origen `http://localhost:3000` no estaba en la lista, por lo que el preflight (`OPTIONS`) era rechazado. Los clientes como Yaak no aplican la política CORS del navegador, por eso solo fallaba desde la app web.

### Solución
Se hizo dinámico `AllowOrigins` según `APP_ENV`:

- **development**: `http://localhost:3000, http://127.0.0.1:3000`
- **production**: `http://origos.no-ip.com, https://origos.no-ip.com`

```go
allowedOrigins := "http://origos.no-ip.com, https://origos.no-ip.com"
if env == "development" {
    allowedOrigins = "http://localhost:3000, http://127.0.0.1:3000"
}
```

### Notas
- `AllowMethods` ya incluye `OPTIONS` y `AllowHeaders` incluye `Authorization`, `Content-Type` y `Origin`, por lo que no requirieron cambios.
- Si el front local cambia de puerto, ajustar `allowedOrigins` en `main.go` para el entorno `development`.
