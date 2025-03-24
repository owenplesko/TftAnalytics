# syntax=docker/dockerfile:1
FROM golang:1.24 AS build
WORKDIR /src

COPY server/ /src/
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server .

FROM gcr.io/distroless/static-debian11

COPY --from=build /bin/server /bin/server

EXPOSE 9000

CMD ["/bin/server"]