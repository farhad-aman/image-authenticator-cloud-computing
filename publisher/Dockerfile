FROM golang:1.21-alpine AS build

WORKDIR /app
COPY . .

RUN go build -o app ./cmd/publisher

FROM alpine:3.18

WORKDIR /root/
COPY --from=build /app/app .

EXPOSE 8080
CMD ["./app"]
