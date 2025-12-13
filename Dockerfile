# Build stage 

FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/pixelforging .

FROM alpine:3.18

RUN apk add --no-cache ca-certificates 

COPY --from=build /bin/pixelforging /usr/local/bin/pixelforging

EXPOSE 9090


ENTRYPOINT [ "usr/local/bin/pixelforging" ]

CMD [ "start-gRPC-server" ]