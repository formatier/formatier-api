FROM golang:alpine as build_image
RUN apk add --no-cache --virtual .build-deps bash gcc build-base
WORKDIR /build_file
COPY . ./
RUN go work sync
RUN go mod download
RUN CGO_ENABLE=0 GOOS=linux go build -a -o app -installsuffix _f7r -ldflags '-extldflags "-static"' ./services/workspace-service/cmd/workspace-service/main.go

FROM alpine:latest
COPY --from=build_image ./build_file/app ./
EXPOSE 1200:3000
ENTRYPOINT [ "./app" ]
