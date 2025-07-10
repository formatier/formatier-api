FROM golang:alpine as build_image
WORKDIR /build_file
COPY . ./
RUN go mod download
RUN GCO_ENABLE=0 go build -a -o gateway-server -installsuffix _f7r /cmd/gateway-server/main.go
COPY ./gateway-server /

FROM alpine:latest
COPY --from=build_image gateway-server .
ENTRYPOINT [ "gateway-server" ]