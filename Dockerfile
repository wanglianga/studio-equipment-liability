FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /studio-equipment-manager .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /studio-equipment-manager .

EXPOSE 8090

CMD ["./studio-equipment-manager"]
