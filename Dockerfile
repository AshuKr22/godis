FROM golang:1.24 as builder
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Final Image

FROM alpine:latest

WORKDIR /data

COPY --from=builder /app/app .
COPY --from=builder /app/log.txt .

EXPOSE 6379

VOLUME ["/data"]

CMD [ "./app" ]



