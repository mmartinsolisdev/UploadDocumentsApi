# Plan de Pruebas Locales — API UploadDocuments (Auth Firebase)

Guía para probar localmente la API con autenticación Firebase y verificar el login del SPA Nuxt. Ejecución manual por el usuario (esta máquina); ningún agente automatiza Postman/navegador.

## Estado actual verificado

- Backend Steps 1–5 del SECURITY-PLAN.md implementados (`go build`/`go vet` OK).
- API corriendo en `127.0.0.1:3001` (modo development) — log OK: `Firebase Auth initialized`, BD `OrigosVCSPT_Temp` conectada.
- `serviceAccount.json` en la raíz del repo (ignorado por `serviceAccount*.json` en `.gitignore`; verificar con `git status` que no aparece).
- Front: `useAPI.ts` ya envía `Authorization: Bearer <token>` y redirige a login en 401.

## Precondiciones (ya verificadas en esta sesión)

- `.env`: `APP_ENV=` vacío (modo development).
- `.env.development`: `FIREBASE_CREDENTIALS=./serviceAccount.json`, `PORT_APP=3001`, credenciales BD dev.
- Usuario de prueba registrado en el proyecto Firebase `origos-auth`.

## Pasos

### 1. Obtener un ID token de prueba (sin front)

`POST https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=AIzaSyAFokpeyZK03AJOcE2UXwUmrxceIk5Hnt0`

Header: `Content-Type: application/json`

Body:
```json
{
  "email": "<email del usuario de prueba>",
  "password": "<password>",
  "returnSecureToken": true
}
```

Respuesta esperada: `idToken` (válido 1 h), `refreshToken`, `expiresIn: "3600"`. La API key es la del front (`app/plugins/firebase.ts:6`).

Alternativa: capturar el token desde el navegador (DevTools → Network → petición al backend → header `Authorization`) tras loguearse en el SPA.

### 2. Probar la API en Postman/cURL

Base URL: `http://127.0.0.1:3001` — usar `127.0.0.1`, NO `localhost` (evita resolución IPv6 `::1`, misma advertencia del plan de Nginx).

| # | Petición | Resultado esperado |
|---|----------|--------------------|
| 1 | `GET /uploader/CombosList` sin header `Authorization` | `401` `{"error":"missing or invalid authorization header"}` |
| 2 | `GET /uploader/CombosList` con token inválido (`Bearer abc123`) | `401` `{"error":"invalid token"}` |
| 3 | `GET /uploader/CombosList` con token válido | `200` con mapa de combos (posiblemente vacío según datos de la BD dev) |
| 4 | `GET /uploader/MembershipsList?Id=...&Language=...&DocName=...&SaleType=...&ContractCode=...` con token válido | `200` con array (puede ser `[]`) |
| 5 | `POST /uploader/UploadFile?Id=...&Language=...&DocName=...&SaleType=...&ContractCode=...` con token válido | `200` (espera 404/5xx si faltan parámetros de BD; el handler devuelve 404 si no llega el archivo) |

Detalles del POST:
- Body: `form-data` con campo `documents` (archivo `.docx` ≤ 20 MB, límite de `BodyLimit` en `main.go`).
- Params en query: `Id`, `Language`, `DocName`, `SaleType`, `ContractCode` (mismos nombres que en `controllers/uploader/uploader.go`).
- Header: `Authorization: Bearer <idToken>`.

Notas:
- El primer `VerifyIDToken` descarga las claves públicas de Google (`www.googleapis.com`) — requiere internet en esta PC; luego quedan cacheadas.
- CORS es irrelevante para Postman/cURL (no aplica origenes de navegador); el middleware solo añade headers.

### 3. Probar el login y flujo del SPA local (opcional)

1. Front (`upload-contracts-next`): `.env` → `API_BASE_URL=http://127.0.0.1:3001`.
2. Backend `main.go:55`: añadir temporalmente `http://localhost:3000` a `AllowOrigins` (sin esto el navegador bloquea la petición por CORS).
3. `pnpm dev` en el front → login con el usuario de prueba → las llamadas a `/uploader` deben responder 200.
4. Probar 401 real: revocar token (Firebase Console → Auth → deshabilitar usuario) → el front debe hacer logout y redirigir a `/`.

### 4. Restaurar antes de desplegar (pendientes de implementación, en servidor)

- Backend `.env`: `APP_ENV=production` (solo en el servidor) — en esta PC dejarlo vacío para seguir en modo dev.
- Backend `.env.production`: rellenar `DB_SERVER/DB_NAME/DB_USER/DB_PASS` reales y `FIREBASE_CREDENTIALS` apuntando al `serviceAccount.json` del servidor.
- Backend `main.go` CORS: quitar `http://localhost:3000` si se añadió para la prueba.
- Front `.env` (dev): devolver `API_BASE_URL=http://origos.no-ip.com`; `.env.production` del front ya apunta ahí.
- Nginx (VPS): `proxy_pass http://localhost:8000;` → `http://127.0.0.1:8000;`.
- Recompilar el `.exe` (`go build -o UploadDocumentsAPI.exe .`) y reemplazarlo en el servidor.
- Paso 8.6 de seguridad: confirmar que `http://<IP-VPS>:8000/uploader/...` no responde desde otra red.

## Validación

- Comandos ya OK: `go build ./...`, `go vet ./...`, `pnpm build` (front).
- Éxito de la prueba: pasos 2.1–2.2 → 401; 2.3–2.5 → 200 con token válido; SPA local: login exitoso y datos cargados.
- `git status` limpio de `serviceAccount.json` (nunca commitear; el JSON no debe aparecer como untracked).

## Riesgos

- Token expira a la hora: si una petición da 401 con un token que antes funcionó, regenerar el token (Paso 1).
- BD dev `OrigosVCSPT_Temp` puede no tener registros → respuestas 200 con arrays/mapas vacíos son válidas.
- `APP_ENV=production` sin `.env.production` relleno rompe el arranque (`log.Fatal` en Firebase init o fallo de conexión BD).