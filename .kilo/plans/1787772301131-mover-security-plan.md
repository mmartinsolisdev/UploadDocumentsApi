# Plan — Mover SECURITY-PLAN.md a .kilo/plans/

## Objetivo
Reorganizar el repo moviendo el documento de referencia `SECURITY-PLAN.md` desde la raíz del proyecto a la carpeta donde Kilo guarda sus planes, con el resto de planes existente.

## Tareas

1. Mover (con historial git) el archivo `SECURITY-PLAN.md` → `.kilo/plans/SECURITY-PLAN.md`.
   - Usar `git mv "SECURITY-PLAN.md" ".kilo/plans/SECURITY-PLAN.md"` para preservar el historial del fichero.
   - Si `git mv` falla por conflicto, usar `Move-Item` y registrar el cambio con `git add -A`.

2. Verificar el archivo movido:
   - `Test-Path .kilo/plans/SECURITY-PLAN.md` → `True`.
   - `Test-Path SECURITY-PLAN.md` (raíz) → `False`.

3. Revisar referencias al documento:
   - La única referencia conocida es una mención textual "SECURITY-PLAN.md" dentro de `.kilo/plans/1787765211167-prueba-local-api-firebase.md` (línea 13). No es un enlace de ruta, por lo que no necesita cambios.
   - Confirmar con `rg "SECURITY-PLAN" -l` que no haya otros enlaces dependientes de la ruta.

4. No commitear nada (solo dejar el cambio en el working tree), salvo que el usuario lo pida.

## Validación
- El archivo existe en la nueva ubicación y no queda copia en la raíz.
- `git status` muestra el rename correcto (sin borrado/añadido accidental).
- No hay referencias de ruta rotas en el repo.

## Riesgos
- Ninguno material: es un fichero de documentación movido dentro del repo, sin importaciones de código ni configuraciones que lo carguen por ruta.
- Si se abren enlaces/entre comillas, volver a la raíz no afecta a código; es solo reubicación.

## Fuera de alcance
- Editar el contenido de `SECURITY-PLAN.md`.
- Aplicar los pasos pendientes del plan (front Nuxt, `.env.production`, verificación final).