# Plan de Pruebas Locales — API UploadDocuments (Auth Firebase)

Guía para probar localmente la API con autenticación Firebase y verificar el login del SPA Nuxt. Ejecución manual por el usuario; ningún agente automatiza Postman/navegador.

## Estado actual (26/08/2026)

**Completado y verificado por el usuario:**

- API corriendo en `127.0.0.1:3001` (modo development) — log OK: `Firebase Auth initialized`, BD `OrigosVCSPT_Temp` conectada.
- `serviceAccount.json` en la raíz del repo (ignorado por `serviceAccount*.json`).
- ID token obtenido vía `signInWithPassword` de Identity Toolkit (REST de Firebase).
- Pruebas de API en Postman OK: sin header → 401, token inválido → 401, token válido → 200 en los endpoints `/uploader/*`.
- Backend Steps 1–5 del SECURITY-PLAN.md: `go build`/`go vet` OK.
- Front (Paso 6) implementado en `upload-contracts-next/app/composables/useAPI.ts`: token fresco `getIdToken()`, header `Authorization: Bearer`, 401 → `signOut` + redirect a `/`, sin sesión → no llama a la API. `pnpm build` OK.

## Pendientes (en orden)

### 1. Dejar de trackear los archivos `.env` (los sigue mostrando GitHub Desktop)

Causa: `.gitignore` ya tiene los patrones (líneas 19–21), pero los tres archivos **ya están trackeados** (fueron commiteados en `e314d94`). `.gitignore` solo aplica a archivos no trackeados — por eso sus cambios siguen apareciendo.

Comandos (ejecutar como agente de implementación o manualmente):

```
git rm --cached .env .env.development .env.production
git commit -m "Stop tracking env files"
```

Efecto: se dejan de trackear sin borrar los archivos del disco; dejan de aparecer en GitHub Desktop.

Advertencias:

- El historial **conserva** las versiones anteriores (`.env.development` tiene credenciales BD reales commiteadas). Untrack no las borra del historial. Si el repo es público, considerar rotar la contraseña de la BD o reescribir el historial (fuera de alcance).
- Tras el untrack, el servidor debe tener sus propios `.env*` copiados a mano en el despliegue (ya existen en el VPS; no se perderán con git).

### 2. Prueba del SPA local (login end-to-end, opcional)

1. Front (`upload-contracts-next`) `.env` → `API_BASE_URL=http://127.0.0.1:3001`.
2. Backend `main.go:55`: añadir temporalmente `http://localhost:3000` a `AllowOrigins` (sin esto el navegador bloquea por CORS).
3. `pnpm dev` → login con el usuario de prueba → llamadas a `/uploader` responden 200.
4. Probar 401 real (opcional): deshabilitar el usuario en Firebase Console → el front hace logout y redirige a `/`.

### 3. Restauración y despliegue (en el servidor/VPS)

- Backend `.env`: `APP_ENV=production` (solo en servidor; en esta PC dejarlo vacío).
- Backend `.env.production` del servidor: `DB_SERVER/DB_NAME/DB_USER/DB_PASS` reales y `FIREBASE_CREDENTIALS=<ruta serviceAccount.json del servidor>`.
- `main.go` CORS: quitar `http://localhost:3000` si se añadió en la prueba local.
- Nginx (VPS): `proxy_pass http://localhost:8000;` → `http://127.0.0.1:8000;`.
- Recompilar el `.exe` (`go build -o UploadDocumentsAPI.exe .`) y reemplazarlo en el servidor.
- Front: rebuild (`pnpm build`) y desplegar el SPA a `C:/_Origos/origosWeb/uploadContracts`; `.env.production` del front ya apunta a `http://origos.no-ip.com`.
- Paso 8.6 de seguridad: confirmar que `http://<IP-VPS>:8000/uploader/...` no responde desde otra red.

## Validación

- `go build ./...`, `go vet ./...`, `pnpm build` → OK (ya verificados).
- Postman: 401/401/200 verificado por el usuario.
- `git status` limpio de `serviceAccount.json` y de los `.env*` tras el untrack.

## Riesgos

- Token Firebase expira a la hora: si una petición da 401 con token antes válido, regenerarlo (signInWithPassword).
- BD dev `OrigosVCSPT_Temp` puede devolver arrays/mapas vacíos: 200 con datos vacíos es válido.
- `APP_ENV=production` sin `.env.production` relleno rompe el arranque (`log.Fatal`).
- .env untrackeados: clonar el repo no incluirá los `.env*`; copiarlos manualmente al servidor.