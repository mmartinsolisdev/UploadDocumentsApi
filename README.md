
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
DB_SERVER=dabataseServer
DB_NAME=databaseName
DB_USER=databaseUser
DB_PASS=databasePass
```

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
