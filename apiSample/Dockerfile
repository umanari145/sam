FROM golang:1.15.7-alpine

ENV ROOT=/go/src/app
ENV CGO_ENABLED 0
WORKDIR ${ROOT}

RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download
EXPOSE 8080

CMD ["go", "run", "main.go"]