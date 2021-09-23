# Start by selecting the base image for our service
FROM golang:1.16.6-alpine3.14 AS Builder

# Directory in which the app will run
WORKDIR /go/src/app

# Move everything from root to the newly created app directory
COPY . .

# Download all neededed project dependencies
RUN go get -d -v ./...

# Build the project executable binary
#RUN go install -v ./...
RUN go build

# Build a New Image and copy only the binary file and de .env files
FROM alpine:3.14
COPY --from=builder /go/src/app/UploadDocumentsAPI .
COPY --from=builder /go/src/app/.env /go/src/app/.env.development /go/src/app/.env.production .
CMD ["./UploadDocumentsAPI"]

# Run/Starts the app executable binary
# CMD ["UploadDocumentsAPI"]
