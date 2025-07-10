FROM golang:alpine AS builder
RUN apk add --no-cache --virtual .build-deps bash gcc git musl-dev

WORKDIR /build-source
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app cmd/orchestrator/main.go
RUN adduser -S -D -H -h /build-source service
USER service

FROM scratch
WORKDIR /app
COPY --from=builder /build-source/app ./
ENTRYPOINT [ "./app" ]