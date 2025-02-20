FROM golang:1.23-alpine as build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/kruise-game-api-server

FROM golang:1.23-alpine

COPY --from=build /bin/kruise-game-api-server /
ENV GIN_MODE=release

ENTRYPOINT ["/kruise-game-api-server"]
