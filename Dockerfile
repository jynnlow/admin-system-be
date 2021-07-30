##FROM command is used to specify our base image - golang image is built on top of linux OS (officially found in docker hub)
#FROM golang:1.16-alpine
#
##WORKDIR command is used to define the working directory of a Docker container
#WORKDIR /app
#
##download necessary go modules
##./ = WORKIR
#COPY go.mod ./
#COPY go.sum ./
#RUN go mod download
#
##copies the entire project, recursively into the container for the build
#COPY . ./
#
##RUN is used to specify commands that are run inside the container
##The -o flag forces build to write the resulting executable or object to the named output file or directory
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /admin-system-be
#
#EXPOSE 8080
#
#CMD ["/app/admin-system-be"]
FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

FROM scratch

WORKDIR /app

COPY --from=builder /app/admin-system-be /usr/bin

ENTRYPOINT ["admin-system-be"]