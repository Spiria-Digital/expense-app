FROM golang:1.23 AS builder

ENV GOBIN=/go/bin
ENV PATH=$GOBIN:$PATH

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o bun cmd/bun/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o api main.go

FROM alpine:3.12

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite libc6-compat

WORKDIR /app
COPY --from=builder /app /app

RUN /app/bun db init && \
    /app/bun db migrate

EXPOSE 8080

CMD ["/app/api"]
