
# API to Upload Documents to Sql Server database

Api developed in **Go** to upload documents to a **Sql Server** database using **Fiber** framework and **Gorm** library.

* [Fiber](https://gofiber.io/) - An Express-inspired web framework written in Go.
* [Gorm](https://gorm.io/) - ORM Library
* [Air](https://github.com/cosmtrek/air) - Used for hot reload.
## Requeriments
Install Golang environment in your O.S. - https://golang.org/

## Project setup installation

Install project dependencies, from the root path of the project run in terminal:

```bash
  go get -d -v ./...
```

## Environment variables

To run this project, you will need to add the following environment variables to the `.env.development` and `.env.production` files.

```bash
PORT_APP=port
TENANTS_CONFIG=<ruta absoluta>/tenants.json
FIREBASE_CREDENTIALS=<ruta absoluta>/serviceAccount.json
```

`TENANTS_CONFIG` apunta al catálogo de clientes y sus bases de datos (solo DSNs, fuera de git). Formato v2 (varias BDs por cliente):

```json
{
  "sirenis": {
    "default": "OrigosVCSPT_Temp",
    "databases": [
      { "dsn": { "server": "HOST", "port": 1433, "database": "OrigosVCSPT_Temp", "user": "USUARIO", "password": "***", "encrypt": false } },
      { "dsn": { "server": "HOST", "port": 1433, "database": "OrigosVCSPT_Prod", "user": "USUARIO", "password": "***", "encrypt": false } }
    ]
  }
}
```

También se acepta el formato legacy de una sola BD por cliente (se normaliza con `default = dsn.database`):

```json
{
  "gtmark": { "dsn": { "server": "HOST", "port": 1433, "database": "BD_GTMARK", "user": "USUARIO", "password": "***", "encrypt": false } }
}
```

La identidad de una BD es su nombre físico (`dsn.database`). La subida (`POST /uploader/UploadFile`) exige el header `X-Client` con el slug del cliente y acepta `X-Database` con el nombre físico de la BD (vacío ⇒ la del `default`). Sin `X-Client` responde `400`; cliente o BD desconocidos, `404`.

Finally set the .env file to `production` for production or `development` for development.

```bash
APP_ENV=production
```
## Deployment

**Development**

To run project in development mode execute in terminal:

```bash
  go run main.go
```

The application uses the air package to restart the server API every time we update the code.  
 To run project with **air** run in the terminal:

```bash
  air
```
**Production**

To generate the project binary file for production, run the command:

```bash
  go build
```
The binary file will be generated in the root path of the project.

**Production with Docker**

Install Docker in your PC. -
https://www.docker.com/products/docker-desktop

Create a Docker image using the Dockerfile. In the project root path run the command:

```bash
  docker build -t docker-image-name .
```
Run the docker image, in terminal execute:

```bash
  docker run -it --rm --name new-container-name image-name
```
