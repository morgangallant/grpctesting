# Build Step
FROM golang:alpine as build
RUN apk add --no-cache protobuf-dev bash
ADD . /build
WORKDIR /build
RUN bash gen.sh
RUN go build -o example server/main.go

# Run Step
FROM golang:alpine
COPY --from=build /build/example /run/example
ENTRYPOINT ["/run/example"]