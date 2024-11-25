FROM golang:1.23-alpine as build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/kruisegamectl

FROM golang:1.23-alpine

COPY --from=build /bin/kruisegamectl /
ENV GIN_MODE=release

CMD ["/kruisegamectl"]
