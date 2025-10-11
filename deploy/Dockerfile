ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/app
COPY go.mod ./
COPY go.sum* ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app ./cmd/uuid-inspector


FROM alpine:latest

RUN apk --no-cache add ca-certificates
COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]
