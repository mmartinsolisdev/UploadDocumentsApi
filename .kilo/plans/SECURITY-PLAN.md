# Plan de Seguridad — API UploadDocuments (Autenticación Firebase)

Documento de referencia para adaptar el proyecto de front con el contexto actual del backend.

> **Estado de implementación:** los Pasos 1–5 (backend) están **implementados y verificados** (`go build`/`go vet` OK). Pendientes: Paso 6 (front Nuxt), Paso 7 (valores reales de `.env.production`), Paso 8 (verificación final) y el cambio de `proxy_pass` en Nginx a `127.0.0.1`.

## Alcance de esta etapa

**Se incluye:**
1. Verificación de ID token de Firebase en el backend (middleware Fiber)
2. Bind de la API a `127.0.0.1` (deja de estar expuesta a internet)
3. Restricción de CORS al dominio real
4. Configuración de credenciales del service account en `.env`
5. Guía de cambios para el front Nuxt 4 (SPA)

**Se excluye (segunda etapa):** HTTPS/TLS, Let's Encrypt, renovación de certificados, rate-limiting.

---

## Contexto actual verificado

- **Backend:** Go 1.27, Fiber v2.52, GORM. Entrada: `main.go` → `routes.Register(app)`. Grupo `/uploader` con 3 endpoints:
  - `GET /uploader/MembershipsList`
  - `GET /uploader/CombosList`
  - `POST /uploader/UploadFile`
- **Listen:** `app.Listen("127.0.0.1:" + port_app)` → solo loopback (implementado).
- **CORS:** `AllowOrigins: "http://origos.no-ip.com, https://origos.no-ip.com"` (implementado; incluye http para la etapa actual sin TLS).
- **Front:** Nuxt 4 SPA estático, servido por Nginx (mismo servidor). Ya tiene login Firebase.
- **Nginx:** sirve el SPA en `C:/_Origos/origosWeb/uploadContracts` y hace proxy de `/uploader` → `http://localhost:8000`.
- **Front llama a la API** vía el proxy Nginx `/uploader` (mismo origen).

### Flujo de una petición del usuario (remoto)

```
Browser del usuario (remoto)
   │  pide https://origos.no-ip.com/uploader/CombosList
   ▼
Nginx (IP pública del VPS, puerto 80/443)
   │  reenvía INTERNAMENTE a 127.0.0.1:8000
   ▼
Backend Go (escucha en 127.0.0.1:8000)
```

El front usa rutas relativas (`/uploader/...`), que se resuelven al mismo dominio → Nginx. El navegador solo habla con Nginx; el backend nunca se expone a internet directamente.

---

## Paso 1 — Backend: añadir dependencia Firebase Admin SDK

**Archivo:** `go.mod` (comando)

```bash
go get firebase.google.com/go/v4@latest
```

**Efecto:** descarga el Admin SDK que verifica tokens de forma segura (firma, expiración, audience) sin llamadas a red.

---

## Paso 2 — Backend: crear middleware de autenticación

**Nuevo archivo:** `middleware/firebase.go` (nuevo paquete `middleware`)

Responsabilidades:
1. Leer header `Authorization` con formato `Bearer <token>`.
2. Si falta o el formato es inválido → `401 Unauthorized` con `{"error":"missing or invalid authorization header"}`.
3. Llamar `client.VerifyIDToken(ctx, token)` del Admin SDK.
4. Si el token es inválido/vencido → `401 Unauthorized`.
5. Si es válido → adjuntar UID al contexto (`c.Locals("uid", token.UID)`) y `c.Next()`.

Consideraciones técnicas:
- La inicialización del cliente Firebase (`firebase.NewApp`) se hace en `main.go` y se guarda en una variable global del paquete (mismo patrón que `database.DBConn`).
- `VerifyIDToken` valida expiración automáticamente (el token de Firebase dura 1 hora).
- El middleware se registra antes de las rutas protegidas y devuelve JSON consistente con el estilo actual.
- `VerifyIDToken` **sí requiere salida a internet** contra Google: en el primer uso (y al refrescar el caché) descarga las claves públicas de Google (`https://www.googleapis.com/robot/v1/metadata/x509/...`) y las cachea. El firewall del VPS debe permitir HTTPS de salida, aunque las claves queden cacheadas después.

---

## Paso 3 — Backend: aplicar middleware al grupo /uploader

**Archivo:** `routes/routes.go`

Cambio: aplicar el middleware al grupo `/uploader`:

```go
uploaderGroup := app.Group("/uploader", middleware.FirebaseAuth)
```

Resultado: los 3 endpoints quedan protegidos sin tocar los handlers.

---

## Paso 4 — Backend: bind a 127.0.0.1 y restringir CORS

**Archivo:** `main.go`

1. Cambiar `app.Listen(":" + port_app)` → `app.Listen("127.0.0.1:" + port_app)`
   - La API solo responde a conexiones locales (Nginx y el propio servidor). Deja de ser alcanzable desde internet aunque alguien conozca la URL.
   - **Importante:** este cambio requiere que Nginx haga el proxy (ya lo hace). Si se ejecuta local sin Nginx, seguirá funcionando vía `localhost`.
2. Cambiar CORS: `AllowOrigins: "*"` → `AllowOrigins: "http://origos.no-ip.com, https://origos.no-ip.com"`.
   - Como el SPA se sirve desde el mismo dominio y Nginx proxea `/uploader`, las peticiones son same-origin; CORS queda como capa defensiva.
   - Se incluye `http://` porque el TLS está diferido a la segunda etapa (hoy el front se sirve por HTTP). Cuando agregues HTTPS podrás dejar solo `https://`.

**Cambio correlativo en Nginx (por compatibilidad):**
- En `nginx.conf`, cambiar `proxy_pass http://localhost:8000;` → `proxy_pass http://127.0.0.1:8000;`
- Motivo: en Windows, `localhost` a veces resuelve a `::1` (IPv6). Si el backend escucha solo en `127.0.0.1` (IPv4) y Nginx resuelve `localhost` a `::1`, la conexión fallaría. Forzar `127.0.0.1` evita este caso.

---

## Paso 5 — Backend: credenciales del service account

**Archivos:** `.env`, `.env.development`, `.env.production`

Añadir en cada uno:

```
FIREBASE_CREDENTIALS=C:/ruta/al/serviceAccount.json
```

**`main.go`:** leer `os.Getenv("FIREBASE_CREDENTIALS")` y pasarla a `firebase.NewApp` con `option.WithCredentialsFile(path)`.

**Generación del JSON (manual, solo tú):**
1. Firebase Console → proyecto → Project settings → **Service accounts**.
2. *Generate new private key* → descarga `serviceAccount.json`.
3. Colocarlo en el servidor (ej. junto al exe) y apuntar `FIREBASE_CREDENTIALS` ahí.
4. **Nunca** commitear este archivo.

**`.gitignore`:** añadir patrón para excluir el JSON del repo (ej. `serviceAccount*.json`).

---

## Paso 6 — Front Nuxt 4: enviar token en cada petición

En el wrapper de llamadas HTTP (axios o fetch nativo) que consume `/uploader`:

1. Obtener token fresco: `await user.getIdToken()` (Firebase renueva automáticamente el token de 1 hora de vida).
2. Añadir header en todas las peticiones:

   ```
   Authorization: Bearer <token>
   ```

3. Interceptor global para respuesta `401`:
   - Forzar logout o redirigir a la pantalla de login (el token pudo expirar o ser inválido).
4. Manejar el caso "usuario sin sesión": no hacer llamadas a la API si no hay usuario autenticado.

**Nota para el front:** el backend no genera el token de sesión. Firebase lo genera en el login; el backend solo lo **verifica**. El front debe adjuntar el token ya existente.

---

## Paso 7 — Actualizar .env.production (valores reales)

**Archivo:** `.env.production` (actualmente vacío)

Rellenar en el servidor:

```
PORT_APP=8000
DB_SERVER=<server real>
DB_NAME=<db real>
DB_USER=<user real>
DB_PASS=<pass real>
FIREBASE_CREDENTIALS=C:/ruta/serviceAccount.json
```

Y en `.env`: `APP_ENV=production`.

---

## Paso 8 — Verificación

1. `go build ./...` y `go vet ./...` → sin errores.
2. Recompilar el `.exe` y reemplazarlo en el servidor.
3. Levantar la API → `log` muestra que escucha en `127.0.0.1:8000`.
4. Probar desde el navegador/Postman:
   - Sin header `Authorization` → `401`.
   - Con token inválido → `401`.
   - Con token válido de un usuario logueado → responde 200/datos.
5. Probar el SPA: login Firebase → llamadas a `/uploader` funcionan.
6. Confirmar que la API **no** responde desde internet (abrir `http://<IP-VPS>:8000/uploader/CombosList` desde otra red → timeout/refused).

---

## Riesgos y notas

- **Nginx ya configurado:** solo requiere el cambio de `localhost` → `127.0.0.1` en `proxy_pass`. El bind a 127.0.0.1 es transparente para el flujo normal.
- **Token expira a la hora:** el front debe usar `getIdToken()` fresco; el interceptor de 401 es clave.
- **Compatibilidad:** Fiber v2 no cambia la firma del middleware; el Admin SDK funciona con Go 1.27.
- **Siguiente etapa (fuera de alcance):** HTTPS con Let's Encrypt/win-acme, renovación automática, rate-limiting en Nginx.