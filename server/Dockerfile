FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o TFTAnalyticsServer .

EXPOSE 9000 9001

CMD ["./TFTAnalyticsServer"]